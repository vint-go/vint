package vintlint0

import (
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	goversion "github.com/hashicorp/go-version"
	"github.com/sourcegraph/conc/pool"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// fileEntry holds a file's name and content after reading from disk.
type fileEntry struct {
	name    string
	content []byte
}

// Linter is the vintlint0 orchestrator. It uses two worker pools to separate
// AST-only rules (no type info needed) from TypeCheck rules, eliminating
// the TypeCheck lock contention that plagues the original architecture.
type Linter struct {
	importer       lint.PackageImporter
	index          *RuleIndex
	config         lint.Config
	rules          []lint.Rule // all rules (needed for disabled-interval computation)
	cache          *rulecache.RuleCache
	configHash     [32]byte
	hasUncacheable bool // precomputed: true if any rule is UncacheableRule
}

// New creates a new vintlint0 Linter.
func New(rules []lint.Rule, config lint.Config) *Linter {
	var imp lint.PackageImporter
	imp = lint.NewSafeSharedImporter()
	l := &Linter{
		importer: imp,
		index:    NewRuleIndex(rules),
		config:   config,
		rules:    rules,
	}

	// Precompute uncacheable flag.
	for _, r := range rules {
		if ur, ok := r.(lint.UncacheableRule); ok && ur.Uncacheable() {
			l.hasUncacheable = true
			break
		}
	}

	// Auto-configure disk cache from environment (same as old linter).
	if os.Getenv("VINT_NO_CACHE") != "1" {
		if dir := os.Getenv("VINT_CACHE_DIR"); dir != "" {
			if rc, err := rulecache.NewWithDir(dir); err == nil {
				l.cache = rc
			}
		}
	}

	// Register rules with cache and compute config hash.
	if l.cache != nil {
		for _, r := range rules {
			ruleConfig := config.Rules[lint.FullRuleName(r)]
			_ = l.cache.RegisterRule(lint.FullRuleName(r), ruleConfig.Arguments)
		}
		l.configHash = lint.ComputeConfigHash(rules, config)
	}

	return l
}

// Lint discovers packages from include/exclude patterns and lints them using
// two worker pools. Returns a channel of failures that is closed when all
// work is complete.
func (l *Linter) Lint(includes, excludes []string) (<-chan lint.Failure, error) {
	packages, err := DiscoverPackages(includes, excludes)
	if err != nil {
		return nil, fmt.Errorf("vintlint0: package discovery: %w", err)
	}

	// Detect Go versions per package from go.mod files.
	perPkgVersions, err := l.detectGoVersions(packages)
	if err != nil {
		return nil, err
	}

	failures := make(chan lint.Failure, 256)

	go l.run(packages, perPkgVersions, failures)

	return failures, nil
}

// detectGoVersions resolves the Go version for each package by finding and
// parsing go.mod files. Caches results per module root to avoid redundant
// parsing.
func (l *Linter) detectGoVersions(packages [][]string) ([]*goversion.Version, error) {
	perModVersions := map[string]*goversion.Version{}
	perPkgVersions := make([]*goversion.Version, len(packages))

	for n, files := range packages {
		if len(files) == 0 {
			continue
		}
		if l.config.GoVersion != nil {
			perPkgVersions[n] = l.config.GoVersion
			continue
		}

		dir, err := filepath.Abs(filepath.Dir(files[0]))
		if err != nil {
			return nil, err
		}

		// Check if we already know the version for this module.
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

		d, v, err := lint.DetectGoMod(dir)
		if err != nil {
			v = lint.DefaultGoVersion
			d = dir
		}
		perModVersions[d] = v
		perPkgVersions[n] = v
	}

	return perPkgVersions, nil
}

// run is the main pipeline goroutine. It dispatches work to two pools and
// closes the failures channel when all work is done.
func (l *Linter) run(packages [][]string, goVersions []*goversion.Version, failures chan<- lint.Failure) {
	defer close(failures)

	// halfProcs := runtime.GOMAXPROCS(0) // trying to make better CPU use..
	halfProcs := runtime.GOMAXPROCS(0) / 2
	if halfProcs < 1 {
		halfProcs = 1
	}

	// AST-only pool: runs rules that don't need type information.
	// One task per file — maximum parallelism, no lock contention.
	astPool := pool.New().WithMaxGoroutines(halfProcs).WithErrors()

	// TypeCheck pool: runs TypeCheck + typecheck-requiring rules.
	// One task per PACKAGE — TypeCheck once, then all rules sequentially.
	// This eliminates the RWMutex contention on Package.
	tcPool := pool.New().WithMaxGoroutines(halfProcs).WithErrors()

	// collectorWg tracks background goroutines that collect failures for
	// package-level cache population. Must complete before FlushAll.
	var collectorWg sync.WaitGroup

	for n, pkgFiles := range packages {
		if len(pkgFiles) == 0 {
			continue
		}
		gover := goVersions[n]

		err := l.processPackage(pkgFiles, gover, astPool, tcPool, failures, &collectorWg)
		if err != nil {
			failures <- lint.NewInternalFailure(
				fmt.Sprintf("error processing package %s: %v", filepath.Dir(pkgFiles[0]), err))
		}
	}

	// Wait for both pools to drain all dispatched work.
	if err := astPool.Wait(); err != nil {
		failures <- lint.NewInternalFailure(fmt.Sprintf("AST pool error: %v", err))
	}
	if err := tcPool.Wait(); err != nil {
		failures <- lint.NewInternalFailure(fmt.Sprintf("TypeCheck pool error: %v", err))
	}

	// Wait for all package cache collectors to finish storing results.
	collectorWg.Wait()

	// Finalize aggregating rules — all files have been collected.
	for _, ar := range l.index.Aggregating {
		fullName := lint.FullRuleName(ar)
		for _, f := range ar.Finalize() {
			if f.RuleName == "" {
				f.RuleName = fullName
			}
			if f.Confidence == 0 {
				f.Confidence = 1
			}
			if f.Confidence >= l.config.Confidence {
				failures <- f
			}
		}
	}

	// Flush cache to disk after all packages are linted.
	if l.cache != nil {
		if err := l.cache.FlushAll(); err != nil {
			failures <- lint.NewInternalFailure(fmt.Sprintf("cache flush error: %v", err))
		}
	}
}

// processPackage reads, parses, and dispatches a single package's files to
// the two worker pools. When cache is enabled, uses a per-package collector
// to populate the package-level cache without sacrificing parallelism.
func (l *Linter) processPackage(
	filenames []string,
	gover *goversion.Version,
	astPool, tcPool *pool.ErrorPool,
	failures chan<- lint.Failure,
	collectorWg *sync.WaitGroup,
) error {
	// Step 1: Read files and filter generated ones.
	var files []fileEntry
	for _, name := range filenames {
		content, err := os.ReadFile(name) //nolint:gosec
		if err != nil {
			return fmt.Errorf("reading %s: %w", name, err)
		}
		if !l.config.IgnoreGeneratedHeader && lint.IsGenerated(content) {
			continue
		}
		// Register file content hash with cache.
		if l.cache != nil {
			l.cache.CacheFileHash(name, content)
		}
		files = append(files, fileEntry{name: name, content: content})
	}
	if len(files) == 0 {
		return nil
	}

	// Step 2: Check package-level cache.
	var pkgActionID [32]byte
	var pkgKey string
	populatePackageCache := false

	if l.cache != nil && !l.hasUncacheable {
		sortedNames := make([]string, len(files))
		for i, f := range files {
			sortedNames[i] = f.name
		}
		sort.Strings(sortedNames)
		pkgKey = filepath.Dir(files[0].name)
		var ok bool
		pkgActionID, ok = l.cache.PackageActionID(sortedNames, l.configHash)
		if ok {
			if cached, hit := l.cache.GetPackage(pkgActionID, pkgKey); hit {
				for _, cf := range cached {
					failures <- lint.FromCachedFailure(cf)
				}
				return nil
			}
			populatePackageCache = true
		}
	}

	// Step 3: Parse AST and create Package + File objects.
	fset := token.NewFileSet()
	pkg := lint.NewPackage(fset, gover, l.importer)

	var lintFiles []*lint.File
	for _, f := range files {
		file, err := pkg.AddFile(f.name, f.content)
		if err != nil {
			lint.AddInvalidFileFailure(f.name, err.Error(), failures)
			continue
		}
		lintFiles = append(lintFiles, file)
	}
	if len(lintFiles) == 0 {
		return nil
	}

	pkg.ExportedScanSortable()
	pkg.Freeze() // all pre-computed fields are now lock-free

	// Step 4: Dispatch to pools.
	//
	// When populatePackageCache is true, we interpose a local channel
	// between the pool tasks and the main failures channel. A collector
	// goroutine forwards failures immediately AND accumulates them so
	// it can store the package cache entry once all tasks complete.
	// This preserves full parallelism — pool tasks are identical in
	// both paths; only the target channel differs.
	targetCh := failures
	var pkgWg sync.WaitGroup
	var localCh chan lint.Failure

	if populatePackageCache {
		localCh = make(chan lint.Failure, 64)
		targetCh = localCh

		// Collector goroutine: forwards failures to main channel and
		// stores package cache entry when localCh is closed.
		collectorWg.Add(1)
		go func() {
			defer collectorWg.Done()
			var collected []lint.Failure
			for f := range localCh {
				collected = append(collected, f)
				failures <- f // forward immediately
			}
			// All pool tasks for this package are done — store cache.
			cachedFs := make([]rulecache.CachedFailure, len(collected))
			for i, f := range collected {
				cachedFs[i] = lint.ToCachedFailure(f)
			}
			l.cache.PutPackage(pkgActionID, pkgKey, cachedFs)
		}()
	}

	// Count tasks BEFORE dispatching so all Add() calls complete
	// before the closer goroutine's Wait() can observe the counter.
	if populatePackageCache {
		taskCount := 0
		if len(l.index.ASTOnly) > 0 {
			taskCount += len(lintFiles)
		}
		if len(l.index.TypeCheck) > 0 {
			taskCount++
		}
		pkgWg.Add(taskCount)
	}

	// Dispatch AST-only rules — one task per file.
	if len(l.index.ASTOnly) > 0 {
		for _, file := range lintFiles {
			astPool.Go(func() error {
				if populatePackageCache {
					defer pkgWg.Done()
				}
				return runRulesOnFile(file, l.index.ASTOnly, l.config, targetCh)
			})
		}
	}

	// Dispatch aggregating rule collection — one task per file.
	// Collect calls go through the AST pool for parallelism but
	// bypass the per-package cache (results come from Finalize later).
	if len(l.index.Aggregating) > 0 {
		for _, file := range lintFiles {
			astPool.Go(func() error {
				for _, ar := range l.index.Aggregating {
					ruleConfig := l.config.Rules[lint.FullRuleName(ar)]
					if ruleConfig.MustExclude(file.Name) {
						continue
					}
					ar.Collect(file, ruleConfig.Arguments)
				}
				return nil
			})
		}
	}

	// Dispatch TypeCheck rules — one task per package.
	if len(l.index.TypeCheck) > 0 {
		tcPool.Go(func() error {
			if populatePackageCache {
				defer pkgWg.Done()
			}
			_ = pkg.TypeCheck()
			for _, file := range lintFiles {
				if err := runRulesOnFile(file, l.index.TypeCheck, l.config, targetCh); err != nil {
					return err
				}
			}
			return nil
		})
	}

	// Closer goroutine: launched AFTER all Add() calls are done,
	// so Wait() cannot see a zero counter prematurely.
	if populatePackageCache {
		go func() {
			pkgWg.Wait()
			close(localCh)
		}()
	}

	return nil
}
