package lint

import (
	"bufio"
	"bytes"
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
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
}

// New creates a new Linter.
func New(reader ReadFile, maxOpenFiles int) Linter {
	var fileReadTokens chan struct{}
	if maxOpenFiles > 0 {
		fileReadTokens = make(chan struct{}, maxOpenFiles)
	}

	l := Linter{
		reader:         reader,
		fileReadTokens: fileReadTokens,
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
	for n := range packages {
		wg.Go(func() error {
			pkg := packages[n]
			gover := perPkgVersions[n]
			if err := l.lintPackage(pkg, gover, ruleSet, config, failures); err != nil {
				return fmt.Errorf("error during linting: %w", err)
			}
			return nil
		})
	}

	go func() {
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

func (l *Linter) lintPackage(filenames []string, gover *goversion.Version, ruleSet []Rule, config Config, failures chan Failure) error {
	if len(filenames) == 0 {
		return nil
	}

	// Classify rules by tier once.
	var fileOnlyRules, otherRules []Rule
	for _, r := range ruleSet {
		if tr, ok := r.(TieredRule); ok && tr.CacheTier() == rulecache.TierFileOnly {
			fileOnlyRules = append(fileOnlyRules, r)
		} else {
			otherRules = append(otherRules, r)
		}
	}

	// Phase 1: Read files, hash content, and check cache for TierFileOnly rules
	// BEFORE parsing AST. This lets us skip parsing entirely when all rules hit cache.
	type fileInfo struct {
		filename string
		content  []byte
		// fileOnlyCacheHits tracks which TierFileOnly rules had cache hits.
		// If a rule name is present, its cached results have already been emitted.
		fileOnlyCacheHits map[string]bool
	}
	var filesToLint []fileInfo

	for _, filename := range filenames {
		content, err := l.readFile(filename)
		if err != nil {
			return err
		}
		if !config.IgnoreGeneratedHeader && isGenerated(content) {
			continue
		}

		// Cache file content hash if cache is available.
		if l.cache != nil {
			l.cache.CacheFileHash(filename, content)
		}

		info := fileInfo{
			filename:          filename,
			content:           content,
			fileOnlyCacheHits: make(map[string]bool, len(fileOnlyRules)),
		}

		// Pre-parse cache check for TierFileOnly rules.
		// These only need (rule config hash + file content hash) — no AST required.
		if l.cache != nil {
			for _, r := range fileOnlyRules {
				ruleConfig := config.Rules[r.Name()]
				if ruleConfig.MustExclude(filename) {
					info.fileOnlyCacheHits[r.Name()] = true
					continue
				}
				if ur, ok := r.(UncacheableRule); ok && ur.Uncacheable() {
					continue
				}
				if cached, hit := l.cache.Get(r.Name(), filename, rulecache.TierFileOnly, nil); hit {
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

		// If all rules hit cache for this file, skip it entirely — no AST parse needed.
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
