package lint

import (
	"errors"
	"fmt"
	"go/ast"
	"go/importer"
	"go/token"
	"go/types"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	goversion "github.com/hashicorp/go-version"
	"golang.org/x/sync/errgroup"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/internal/typeparams"
)

// sharedImporter wraps a types.ImporterFrom with a lock-free cache so that
// already-resolved packages are returned instantly without filesystem
// syscalls. The inner importer's FindPkg is expensive (EvalSymlinks etc.)
// and is called even on map hits, so we cache results in a sync.Map keyed
// by import path to bypass it entirely for known packages.
type sharedImporter struct {
	cache sync.Map   // import path → *importResult (lock-free reads)
	mu    sync.Mutex // serializes actual import resolution (cache misses)
	inner types.ImporterFrom
	fset  *token.FileSet // fset used by gcimporter; positions of imported objects live here
}

type importResult struct {
	pkg *types.Package
	err error
}

func newSharedImporter() *sharedImporter {
	fset := token.NewFileSet()
	imp := importer.ForCompiler(fset, runtime.Compiler, nil)
	return &sharedImporter{
		fset:  fset,
		inner: imp.(types.ImporterFrom),
	}
}

func (s *sharedImporter) Import(path string) (*types.Package, error) {
	return s.ImportFrom(path, "", 0)
}

func (s *sharedImporter) ImportFrom(path, srcDir string, mode types.ImportMode) (*types.Package, error) {
	// Fast path: lock-free cache check — no FindPkg, no syscalls.
	if v, ok := s.cache.Load(path); ok {
		r := v.(*importResult)
		return r.pkg, r.err
	}

	// Slow path: resolve under lock (first time only per import path).
	s.mu.Lock()
	defer s.mu.Unlock()

	// Double-check after acquiring lock.
	if v, ok := s.cache.Load(path); ok {
		r := v.(*importResult)
		return r.pkg, r.err
	}

	pkg, err := s.inner.ImportFrom(path, srcDir, mode)
	s.cache.Store(path, &importResult{pkg: pkg, err: err})
	return pkg, err
}

// Package represents a package in the project.
type Package struct {
	fset     *token.FileSet
	importer *sharedImporter

	mu        sync.RWMutex
	files     map[string]*File
	goVersion *goversion.Version
	typesPkg  *types.Package
	typesInfo *types.Info
	// sortable is the set of types in the package that implement sort.Interface.
	sortable map[string]bool
	// main is whether this is a "main" package.
	main int
}

var (
	trueValue  = 1
	falseValue = 2

	// Go115 is a constant representing the Go version 1.15.
	Go115 = goversion.Must(goversion.NewVersion("1.15"))
	// Go121 is a constant representing the Go version 1.21.
	Go121 = goversion.Must(goversion.NewVersion("1.21"))
	// Go122 is a constant representing the Go version 1.22.
	Go122 = goversion.Must(goversion.NewVersion("1.22"))
	// Go124 is a constant representing the Go version 1.24.
	Go124 = goversion.Must(goversion.NewVersion("1.24"))
	// Go125 is a constant representing the Go version 1.25.
	Go125 = goversion.Must(goversion.NewVersion("1.25"))
)

// Files return package's files.
func (p *Package) Files() map[string]*File {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.files
}

// IsMain returns if that's the main package.
func (p *Package) IsMain() bool {
	// Fast path: read lock only to avoid contention when
	// multiple routines check IsMain on the same package.
	p.mu.RLock()
	cached := p.main
	p.mu.RUnlock()
	if cached != 0 {
		return cached == trueValue
	}

	// Slow path: acquire write lock to compute and cache result.
	p.mu.Lock()
	defer p.mu.Unlock()

	switch p.main {
	case trueValue:
		return true
	case falseValue:
		return false
	}
	for _, f := range p.files {
		if f.isMain() {
			p.main = trueValue
			return true
		}
	}
	p.main = falseValue
	return false
}

// TypesPkg yields information on this package.
func (p *Package) TypesPkg() *types.Package {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.typesPkg
}

// TypesInfo yields type information of this package identifiers.
func (p *Package) TypesInfo() *types.Info {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.typesInfo
}

// ImportedPkgSourceDir resolves an import path to its source directory
// using position data already loaded by the type checker's importer.
// Returns ("", false) if the package wasn't imported or has no position info.
func (p *Package) ImportedPkgSourceDir(importPath string) (string, bool) {
	if p.importer == nil {
		return "", false
	}

	v, ok := p.importer.cache.Load(importPath)
	if !ok {
		return "", false
	}

	r := v.(*importResult)
	if r.err != nil || r.pkg == nil {
		return "", false
	}

	// Find any exported object with a valid position — its filename
	// reveals the source directory. The gcimporter uses "$GOROOT" as a
	// placeholder in filenames, so we expand it to the real path.
	goroot := runtime.GOROOT()
	scope := r.pkg.Scope()
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		if obj == nil || !obj.Pos().IsValid() {
			continue
		}

		pos := p.importer.fset.Position(obj.Pos())
		if pos.Filename == "" {
			continue
		}

		filename := pos.Filename
		if strings.HasPrefix(filename, "$GOROOT") {
			filename = goroot + filename[len("$GOROOT"):]
		}

		return filepath.Dir(filename), true
	}

	return "", false
}

// Sortable yields a map of sortable types in this package.
func (p *Package) Sortable() map[string]bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.sortable
}

// TypeCheck performs type checking for given package.
func (p *Package) TypeCheck() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	alreadyTypeChecked := p.typesInfo != nil || p.typesPkg != nil
	if alreadyTypeChecked {
		return nil
	}

	var imp types.Importer
	if p.importer != nil {
		imp = p.importer
	} else {
		imp = importer.Default()
	}
	config := &types.Config{
		// By setting a no-op error reporter, the type checker does as much work as possible.
		Error:    func(error) {},
		Importer: imp,
	}
	info := &types.Info{
		Types:      map[ast.Expr]types.TypeAndValue{},
		Defs:       map[*ast.Ident]types.Object{},
		Uses:       map[*ast.Ident]types.Object{},
		Scopes:     map[ast.Node]*types.Scope{},
		Selections: map[*ast.SelectorExpr]*types.Selection{},
	}
	var anyFile *File
	var astFiles []*ast.File
	for _, f := range p.files {
		anyFile = f
		astFiles = append(astFiles, f.AST)
	}

	if anyFile == nil {
		// this is unlikely to happen, but technically guarantees anyFile to not be nil
		return errors.New("no ast.File found")
	}

	typesPkg, err := check(config, anyFile.AST.Name.Name, p.fset, astFiles, info)

	// Remember the typechecking info, even if config.Check failed,
	// since we will get partial information.
	p.typesPkg = typesPkg
	p.typesInfo = info

	return err
}

// check function encapsulates the call to [go/types.Config.Check] method and
// recovers if the called method panics (see issue #59).
func check(config *types.Config, n string, fset *token.FileSet, astFiles []*ast.File, info *types.Info) (p *types.Package, err error) {
	defer func() {
		if r := recover(); r != nil {
			err, _ = r.(error)
			p = nil
			return
		}
	}()

	return config.Check(n, fset, astFiles, info)
}

// TypeOf returns the type of expression.
func (p *Package) TypeOf(expr ast.Expr) types.Type {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.typesInfo == nil {
		return nil
	}

	return p.typesInfo.TypeOf(expr)
}

type sortableMethodsFlags int

// flags for sortable interface methods.
const (
	bfLen sortableMethodsFlags = 1 << iota
	bfLess
	bfSwap
)

func (p *Package) scanSortable() {
	p.mu.Lock()
	defer p.mu.Unlock()

	sortableFlags := map[string]sortableMethodsFlags{}
	for _, f := range p.files {
		for _, decl := range f.AST.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			isAMethodDeclaration := ok && fn.Recv != nil && len(fn.Recv.List) != 0
			if !isAMethodDeclaration {
				continue
			}

			recvType := typeparams.ReceiverType(fn)
			sortableFlags[recvType] |= getSortableMethodFlagForFunction(fn)
		}
	}

	p.sortable = make(map[string]bool, len(sortableFlags))
	for typ, ms := range sortableFlags {
		if ms == bfLen|bfLess|bfSwap {
			p.sortable[typ] = true
		}
	}
}

func (p *Package) lint(rules []Rule, config Config, failures chan Failure, rc *rulecache.RuleCache, preCachedHits map[string]map[string]bool) error {
	p.scanSortable()
	var eg errgroup.Group
	eg.SetLimit(runtime.GOMAXPROCS(0))
	for name, file := range p.Files() {
		hits := preCachedHits[name] // may be nil
		eg.Go(func() error {
			return file.lint(rules, config, failures, rc, hits)
		})
	}

	return eg.Wait()
}

// IsAtLeastGoVersion returns true if the Go version for this package is v or higher, false otherwise.
func (p *Package) IsAtLeastGoVersion(v *goversion.Version) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.goVersion.GreaterThanOrEqual(v)
}

// GoVersionString returns the Go version for this package as a "go1.X" string
// suitable for use with golang.org/x/tools/internal/versions.
func (p *Package) GoVersionString() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.goVersion == nil {
		return ""
	}
	segments := p.goVersion.Segments()
	if len(segments) < 2 {
		return "go" + p.goVersion.String()
	}
	return fmt.Sprintf("go%d.%d", segments[0], segments[1])
}

func getSortableMethodFlagForFunction(fn *ast.FuncDecl) sortableMethodsFlags {
	switch {
	case astutils.FuncSignatureIs(fn, "Len", []string{}, []string{"int"}):
		return bfLen
	case astutils.FuncSignatureIs(fn, "Less", []string{"int", "int"}, []string{"bool"}):
		return bfLess
	case astutils.FuncSignatureIs(fn, "Swap", []string{"int", "int"}, []string{}):
		return bfSwap
	default:
		return 0
	}
}
