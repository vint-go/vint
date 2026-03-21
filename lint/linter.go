package lint

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"

	goversion "github.com/hashicorp/go-version"
	"github.com/strowk/vint/internal/rulecache"
	"golang.org/x/mod/modfile"
	"golang.org/x/sync/errgroup"
)

// ReadFile defines an abstraction for reading files.
type ReadFile func(path string) (result []byte, err error)

// Linter is used for linting set of files.
type Linter struct {
	reader         ReadFile
	fileReadTokens chan struct{}
	cache          *rulecache.RuleCache
	importer       packageImporter
}

// New creates a new Linter.
func New(reader ReadFile, maxOpenFiles int) Linter {
	var fileReadTokens chan struct{}
	if maxOpenFiles > 0 {
		fileReadTokens = make(chan struct{}, maxOpenFiles)
	}

	var imp packageImporter
	if os.Getenv("VINT_SAFE_IMPORT") == "1" {
		imp = newSafeImporter()
	} else {
		imp = newSharedImporter()
	}
	l := Linter{
		reader:         reader,
		fileReadTokens: fileReadTokens,
		importer:       imp,
	}

	// Auto-configure disk cache from environment.
	if os.Getenv("VINT_NO_CACHE") != "1" {
		if dir := os.Getenv("VINT_CACHE_DIR"); dir != "" {
			if rc, err := rulecache.NewWithDir(dir); err == nil {
				l.cache = rc
			}
		}
	}

	return l
}

func (l *Linter) readFile(path string) (result []byte, err error) {
	if l.fileReadTokens != nil {
		// "take" a token by writing to the channel.
		// It will block if no more space in the channel's buffer
		l.fileReadTokens <- struct{}{}
		defer func() {
			// "free" a token by reading from the channel
			<-l.fileReadTokens
		}()
	}

	return l.reader(path)
}

// fileData holds a file's name and content after reading from disk.
type fileData struct {
	filename string
	content  []byte
}

// computeConfigHash produces a single hash covering all rule configurations
// and global settings that affect lint output. Computed once per Lint() call.
func computeConfigHash(ruleSet []Rule, config Config) [32]byte {
	h := sha256.New()
	names := make([]string, len(ruleSet))
	for i, r := range ruleSet {
		names[i] = r.Name()
	}
	sort.Strings(names)
	for _, name := range names {
		rc := config.Rules[name]
		fmt.Fprintf(h, "rule %s\n", name)
		if len(rc.Arguments) > 0 {
			var buf bytes.Buffer
			_ = gob.NewEncoder(&buf).Encode(rc.Arguments)
			h.Write(buf.Bytes())
		}
		for _, ex := range rc.Exclude {
			fmt.Fprintf(h, "exclude %s\n", ex)
		}
	}
	fmt.Fprintf(h, "confidence %g\n", config.Confidence)
	// Sort directive keys for deterministic hashing.
	dirKeys := make([]string, 0, len(config.Directives))
	for k := range config.Directives {
		dirKeys = append(dirKeys, k)
	}
	sort.Strings(dirKeys)
	for _, k := range dirKeys {
		fmt.Fprintf(h, "directive %s %s\n", k, config.Directives[k].Severity)
	}
	var id [32]byte
	copy(id[:], h.Sum(nil))
	return id
}

var (
	generatedPrefix  = []byte("// Code generated ")
	generatedSuffix  = []byte(" DO NOT EDIT.")
	defaultGoVersion = goversion.Must(goversion.NewVersion("1.0"))
)

// Lint lints a set of files with the specified rule.
func (l *Linter) Lint(packages [][]string, ruleSet []Rule, config Config) (<-chan Failure, error) {
	failures := make(chan Failure)

	// Register rules with cache if available.
	if l.cache != nil {
		for _, r := range ruleSet {
			ruleConfig := config.Rules[r.Name()]
			if err := l.cache.RegisterRule(r.Name(), ruleConfig.Arguments); err != nil {
				return nil, fmt.Errorf("cache: register rule %q: %w", r.Name(), err)
			}
		}
	}

	// Compute config hash once for package-level caching.
	var configHash [32]byte
	if l.cache != nil {
		configHash = computeConfigHash(ruleSet, config)
	}

	perModVersions := map[string]*goversion.Version{}
	perPkgVersions := make([]*goversion.Version, len(packages))
	for n, files := range packages {
		if len(files) == 0 {
			continue
		}
		if config.GoVersion != nil {
			perPkgVersions[n] = config.GoVersion
			continue
		}

		dir, err := filepath.Abs(filepath.Dir(files[0]))
		if err != nil {
			return nil, err
		}

		alreadyKnownMod := false
		for d, v := range perModVersions {
			if strings.HasPrefix(dir, d) {
				perPkgVersions[n] = v
				alreadyKnownMod = true
				break
			}
		}
		if alreadyKnownMod {
			continue
		}

		d, v, err := detectGoMod(dir)
		if err != nil {
			// No luck finding the go.mod file thus set the default Go version
			v = defaultGoVersion
			d = dir
		}
		perModVersions[d] = v
		perPkgVersions[n] = v
	}

	var wg errgroup.Group
	wg.SetLimit(2 * runtime.GOMAXPROCS(0))

	go func() {
		for n := range packages {
			wg.Go(func() error {
				pkg := packages[n]
				gover := perPkgVersions[n]
				if err := l.lintPackage(pkg, gover, ruleSet, config, configHash, failures); err != nil {
					return fmt.Errorf("error during linting: %w", err)
				}
				return nil
			})
		}
		err := wg.Wait()
		if err != nil {
			failures <- NewInternalFailure(err.Error())
		}
		// Flush cache to disk after all packages are linted.
		if l.cache != nil {
			if flushErr := l.cache.FlushAll(); flushErr != nil {
				failures <- NewInternalFailure(flushErr.Error())
			}
		}
		close(failures)
	}()

	return failures, nil
}

func (l *Linter) lintPackage(filenames []string, gover *goversion.Version, ruleSet []Rule, config Config, configHash [32]byte, failures chan Failure) error {
	if len(filenames) == 0 {
		return nil
	}

	// Read all files, hash content, and filter generated files.
	allFiles := make([]fileData, 0, len(filenames))
	for _, filename := range filenames {
		content, err := l.readFile(filename)
		if err != nil {
			return err
		}
		if !config.IgnoreGeneratedHeader && isGenerated(content) {
			continue
		}
		if l.cache != nil {
			l.cache.CacheFileHash(filename, content)
		}
		allFiles = append(allFiles, fileData{filename: filename, content: content})
	}
	if len(allFiles) == 0 {
		return nil
	}

	// Check for uncacheable rules — if any exist, skip package-level cache.
	hasUncacheable := false
	for _, r := range ruleSet {
		if ur, ok := r.(UncacheableRule); ok && ur.Uncacheable() {
			hasUncacheable = true
			break
		}
	}
	usePkgCache := l.cache != nil && !hasUncacheable

	if usePkgCache {
		// Compute package-level action ID.
		sortedNames := make([]string, len(allFiles))
		for i, f := range allFiles {
			sortedNames[i] = f.filename
		}
		sort.Strings(sortedNames)
		pkgKey := filepath.Dir(allFiles[0].filename)
		pkgActionID, ok := l.cache.PackageActionID(sortedNames, configHash)

		if ok {
			if cached, hit := l.cache.GetPackage(pkgActionID, pkgKey); hit {
				for _, cf := range cached {
					failures <- fromCachedFailure(cf)
				}
				return nil
			}
		}

		// Package cache miss — run full analysis, collecting failures.
		localCh := make(chan Failure, 64)
		var collected []Failure
		done := make(chan struct{})
		go func() {
			for f := range localCh {
				collected = append(collected, f)
			}
			close(done)
		}()

		err := l.lintPackageCore(allFiles, gover, ruleSet, config, localCh)
		close(localCh)
		<-done

		// Store in package cache.
		if ok {
			cachedFs := make([]rulecache.CachedFailure, len(collected))
			for i, f := range collected {
				cachedFs[i] = toCachedFailure(f)
			}
			l.cache.PutPackage(pkgActionID, pkgKey, cachedFs)
		}

		// Forward to real channel.
		for _, f := range collected {
			failures <- f
		}
		return err
	}

	return l.lintPackageCore(allFiles, gover, ruleSet, config, failures)
}

// lintPackageCore runs per-rule caching (Phase 1) and AST parsing + linting (Phase 2).
func (l *Linter) lintPackageCore(files []fileData, gover *goversion.Version, ruleSet []Rule, config Config, failures chan Failure) error {
	// Classify rules by tier once.
	var fileOnlyRules, otherRules []Rule
	for _, r := range ruleSet {
		if tr, ok := r.(TieredRule); ok && tr.CacheTier() == rulecache.TierFileOnly {
			fileOnlyRules = append(fileOnlyRules, r)
		} else {
			otherRules = append(otherRules, r)
		}
	}

	// Phase 1: Check cache for TierFileOnly rules BEFORE parsing AST.
	type fileInfo struct {
		filename          string
		content           []byte
		fileOnlyCacheHits map[string]bool
	}
	var filesToLint []fileInfo

	for _, fd := range files {
		info := fileInfo{
			filename:          fd.filename,
			content:           fd.content,
			fileOnlyCacheHits: make(map[string]bool, len(fileOnlyRules)),
		}

		if l.cache != nil {
			for _, r := range fileOnlyRules {
				ruleConfig := config.Rules[r.Name()]
				if ruleConfig.MustExclude(fd.filename) {
					info.fileOnlyCacheHits[r.Name()] = true
					continue
				}
				if ur, ok := r.(UncacheableRule); ok && ur.Uncacheable() {
					continue
				}
				if cached, hit := l.cache.Get(r.Name(), fd.filename, rulecache.TierFileOnly, nil, nil); hit {
					for _, cf := range cached {
						failure := fromCachedFailure(cf)
						if failure.Confidence >= config.Confidence {
							failures <- failure
						}
					}
					info.fileOnlyCacheHits[r.Name()] = true
				}
			}
		}

		allHit := len(info.fileOnlyCacheHits) == len(fileOnlyRules) && len(otherRules) == 0
		if allHit {
			continue
		}

		filesToLint = append(filesToLint, info)
	}

	if len(filesToLint) == 0 {
		return nil
	}

	// Phase 2: Parse AST only for files that had cache misses.
	pkg := &Package{
		fset:      token.NewFileSet(),
		importer:  l.importer,
		files:     map[string]*File{},
		goVersion: gover,
	}
	fileOnlyHits := map[string]map[string]bool{} // filename → rule hits from phase 1
	for _, info := range filesToLint {
		file, err := NewFile(info.filename, info.content, pkg)
		if err != nil {
			addInvalidFileFailure(info.filename, err.Error(), failures)
			continue
		}
		pkg.files[info.filename] = file
		fileOnlyHits[info.filename] = info.fileOnlyCacheHits
	}

	if len(pkg.files) == 0 {
		return nil
	}

	return pkg.lint(ruleSet, config, failures, l.cache, fileOnlyHits)
}

func detectGoMod(dir string) (rootDir string, ver *goversion.Version, err error) {
	modFileName, err := retrieveModFile(dir)
	if err != nil {
		return "", nil, fmt.Errorf("%q doesn't seem to be part of a Go module", dir)
	}

	mod, err := os.ReadFile(modFileName) //nolint:gosec // ignore G304: potential file inclusion via variable
	if err != nil {
		return "", nil, fmt.Errorf("failed to read %q, got %w", modFileName, err)
	}

	modAst, err := modfile.ParseLax(modFileName, mod, nil)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse %q, got %w", modFileName, err)
	}

	if modAst.Go == nil {
		return "", nil, fmt.Errorf("%q does not specify a Go version", modFileName)
	}

	ver, err = goversion.NewVersion(modAst.Go.Version)
	return filepath.Dir(modFileName), ver, err
}

func retrieveModFile(dir string) (string, error) {
	const lookingForFile = "go.mod"
	for {
		// filepath.Dir returns 'C:\' on Windows, and '/' on Unix
		isRootDir := (dir == filepath.VolumeName(dir)+string(filepath.Separator))
		if dir == "." || isRootDir {
			return "", fmt.Errorf("did not found %q file", lookingForFile)
		}

		lookingForFilePath := filepath.Join(dir, lookingForFile)
		info, err := os.Stat(lookingForFilePath)
		if err != nil || info.IsDir() {
			// lets check the parent dir
			dir = filepath.Dir(dir)
			continue
		}

		return lookingForFilePath, nil
	}
}

// isGenerated reports whether the source file is generated code
// according to the rules from https://go.dev/s/generatedcode.
// This is inherited from the original go lint.
func isGenerated(src []byte) bool {
	sc := bufio.NewScanner(bytes.NewReader(src))
	for sc.Scan() {
		b := sc.Bytes()
		if bytes.HasPrefix(b, generatedPrefix) && bytes.HasSuffix(b, generatedSuffix) && len(b) >= len(generatedPrefix)+len(generatedSuffix) {
			return true
		}
	}
	return false
}

// addInvalidFileFailure adds a failure for an invalid formatted file.
func addInvalidFileFailure(filename, errStr string, failures chan Failure) {
	position := getPositionInvalidFile(filename, errStr)
	failures <- Failure{
		Confidence: 1,
		Failure:    fmt.Sprintf("invalid file %s: %v", filename, errStr),
		Category:   failureCategoryValidity,
		Position:   position,
	}
}

// errPosRegexp matches with a NewFile error message:
//
//	corrupted.go:10:4: expected '}', found 'EOF
//
// The first group matches the line, and the second group matches the column.
var errPosRegexp = regexp.MustCompile(`.*:(\d*):(\d*):.*$`)

// getPositionInvalidFile gets the position of the error in an invalid file.
func getPositionInvalidFile(filename, s string) FailurePosition {
	pos := errPosRegexp.FindStringSubmatch(s)
	if len(pos) < 3 {
		return FailurePosition{}
	}
	line, err := strconv.Atoi(pos[1])
	if err != nil {
		return FailurePosition{}
	}
	column, err := strconv.Atoi(pos[2])
	if err != nil {
		return FailurePosition{}
	}

	return FailurePosition{
		Start: token.Position{
			Filename: filename,
			Line:     line,
			Column:   column,
		},
	}
}
