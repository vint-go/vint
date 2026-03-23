package no_deprecated_usage

import (
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
	"sync"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoDeprecatedUsageRule flags usage of deprecated functions, variables, constants,
// or fields from the standard library, third-party packages, and the current package.
type NoDeprecatedUsageRule struct{}

// pkgDeprecationInfo holds all deprecation data for a single import path,
// computed once from a single parsePackageDoc call.
type pkgDeprecationInfo struct {
	pkgDeprecated     bool
	deprecatedSymbols map[string]bool
}

// pkgInfoEntry ensures each import path is resolved exactly once,
// even under concurrent access from multiple goroutines.
type pkgInfoEntry struct {
	once sync.Once
	info *pkgDeprecationInfo
}

// pkgInfoCache maps import path → *pkgInfoEntry.
var pkgInfoCache sync.Map

// Apply applies the rule to given file.
func (r *NoDeprecatedUsageRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	info := file.Pkg.TypesInfo()
	if info == nil {
		return nil
	}

	// Phase 1: Collect deprecated declarations in the current package.
	localDeprecated := collectLocalDeprecated(file.Pkg)

	// Phase 2: Walk through all identifier usages in this file and check for deprecated references.
	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	w := &lintDeprecatedUsage{
		pkg:             file.Pkg,
		info:            info,
		localDeprecated: localDeprecated,
		onFailure:       onFailure,
	}

	ast.Inspect(file.AST, w.inspect)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoDeprecatedUsageRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	info := file.Pkg.TypesInfo()
	if info == nil {
		return nil
	}

	localDeprecated := collectLocalDeprecated(file.Pkg)

	w := &lintDeprecatedUsage{
		pkg:             file.Pkg,
		info:            info,
		localDeprecated: localDeprecated,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	w.inspect(node)

	return failures
}

type lintDeprecatedUsage struct {
	pkg             *lint.Package
	info            *types.Info
	localDeprecated map[types.Object]bool
	onFailure       func(lint.Failure)
}

func (w *lintDeprecatedUsage) inspect(node ast.Node) bool {
	ident, ok := node.(*ast.Ident)
	if !ok {
		return true
	}

	obj, ok := w.info.Uses[ident]
	if !ok {
		return true
	}

	// Skip the universe scope (built-ins).
	if obj.Pkg() == nil {
		return true
	}

	// Skip package name references — we report the actual symbols, not the package import itself.
	if _, ok := obj.(*types.PkgName); ok {
		return true
	}

	// Check if this is a deprecated local (current package) object.
	if w.localDeprecated[obj] {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryBadPractice,
			Confidence: 1,
			Node:       ident,
			Failure:    deprecatedMessage(obj),
		})
		return true
	}

	// For objects from external packages, check if the specific symbol is deprecated,
	// or if the entire package is deprecated.
	objPkg := obj.Pkg()
	if objPkg != nil {
		importPath := objPkg.Path()
		depInfo := getDeprecationInfo(importPath, w.pkg)
		if depInfo.deprecatedSymbols[obj.Name()] || depInfo.pkgDeprecated {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryBadPractice,
				Confidence: 1,
				Node:       ident,
				Failure:    deprecatedMessage(obj),
			})
		}
	}

	return true
}

// deprecatedMessage produces a failure message for a deprecated object.
func deprecatedMessage(obj types.Object) string {
	pkg := obj.Pkg()
	if pkg == nil {
		return obj.Name() + " is deprecated"
	}
	return pkg.Path() + "." + obj.Name() + " is deprecated"
}

// collectLocalDeprecated scans all files in the current package's AST
// and returns a set of types.Object that are marked as deprecated.
func collectLocalDeprecated(pkg *lint.Package) map[types.Object]bool {
	info := pkg.TypesInfo()
	if info == nil {
		return nil
	}

	deprecated := map[types.Object]bool{}

	for _, f := range pkg.Files() {
		for _, decl := range f.AST.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if isDocDeprecated(d.Doc) {
					if obj := info.Defs[d.Name]; obj != nil {
						deprecated[obj] = true
					}
				}
			case *ast.GenDecl:
				// For grouped declarations (var, const, type), check each spec.
				groupDoc := d.Doc
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.ValueSpec:
						docToCheck := s.Doc
						if docToCheck == nil {
							docToCheck = groupDoc
						}
						if isDocDeprecated(docToCheck) {
							for _, name := range s.Names {
								if obj := info.Defs[name]; obj != nil {
									deprecated[obj] = true
								}
							}
						}
					case *ast.TypeSpec:
						docToCheck := s.Doc
						if docToCheck == nil {
							docToCheck = groupDoc
						}
						if isDocDeprecated(docToCheck) {
							if obj := info.Defs[s.Name]; obj != nil {
								deprecated[obj] = true
							}
						}
					}
				}
			}
		}
	}

	return deprecated
}

// isDocDeprecated checks whether a doc comment group contains a "Deprecated:" marker.
func isDocDeprecated(doc *ast.CommentGroup) bool {
	if doc == nil {
		return false
	}
	text := doc.Text()
	// According to Go convention, a deprecation notice is a paragraph
	// that starts with "Deprecated: ".
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "Deprecated:") {
			return true
		}
	}
	return false
}

// getDeprecationInfo returns the deprecation info for an import path,
// computing it at most once via sync.Once. Uses the type checker's
// already-resolved package data to find source directories without
// spawning subprocesses.
func getDeprecationInfo(importPath string, pkg *lint.Package) *pkgDeprecationInfo {
	v, _ := pkgInfoCache.LoadOrStore(importPath, &pkgInfoEntry{})
	entry := v.(*pkgInfoEntry)
	entry.once.Do(func() {
		entry.info = loadDeprecationInfo(importPath, pkg)
	})
	return entry.info
}

// loadDeprecationInfo does the actual work: finds the package source dir
// (from the type checker if possible), parses, and extracts doc comments.
func loadDeprecationInfo(importPath string, pkg *lint.Package) *pkgDeprecationInfo {
	docPkg := parsePackageDoc(importPath, pkg)
	if docPkg == nil {
		return &pkgDeprecationInfo{}
	}

	symbols := map[string]bool{}

	for _, f := range docPkg.Funcs {
		if isTextDeprecated(f.Doc) {
			symbols[f.Name] = true
		}
	}

	for _, v := range docPkg.Vars {
		if isTextDeprecated(v.Doc) {
			for _, name := range v.Names {
				symbols[name] = true
			}
		}
	}

	for _, c := range docPkg.Consts {
		if isTextDeprecated(c.Doc) {
			for _, name := range c.Names {
				symbols[name] = true
			}
		}
	}

	for _, t := range docPkg.Types {
		if isTextDeprecated(t.Doc) {
			symbols[t.Name] = true
		}
		for _, f := range t.Funcs {
			if isTextDeprecated(f.Doc) {
				symbols[f.Name] = true
			}
		}
		for _, m := range t.Methods {
			if isTextDeprecated(m.Doc) {
				symbols[m.Name] = true
			}
		}
		for _, v := range t.Vars {
			if isTextDeprecated(v.Doc) {
				for _, name := range v.Names {
					symbols[name] = true
				}
			}
		}
		for _, c := range t.Consts {
			if isTextDeprecated(c.Doc) {
				for _, name := range c.Names {
					symbols[name] = true
				}
			}
		}
	}

	return &pkgDeprecationInfo{
		pkgDeprecated:     isTextDeprecated(docPkg.Doc),
		deprecatedSymbols: symbols,
	}
}

// parsePackageDoc finds and parses the source of a package, returning a go/doc.Package.
// It first tries to resolve the source directory from the type checker's importer
// (which already resolved these packages during TypeCheck), avoiding expensive
// build.Import calls that spawn subprocesses in module mode.
// Falls back to build.Import for packages not found in the type checker's cache.
func parsePackageDoc(importPath string, pkg *lint.Package) *doc.Package {
	dir, ok := pkg.ImportedPkgSourceDir(importPath)
	if !ok {
		return nil
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return nil
	}

	// Return the first non-test package.
	for _, p := range pkgs {
		if !strings.HasSuffix(p.Name, "_test") {
			return doc.New(p, importPath, doc.AllDecls)
		}
	}

	return nil
}

// isTextDeprecated checks if a doc string contains "Deprecated:" as a paragraph start.
func isTextDeprecated(text string) bool {
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "Deprecated:") {
			return true
		}
	}
	return false
}

// Name returns the rule name.
func (*NoDeprecatedUsageRule) Name() string {
	return "noDeprecatedUsage"
}

// Group returns the rule group.
func (*NoDeprecatedUsageRule) Group() string {
	return "correctness"
}

// RequiresTypecheck returns true since this rule needs type information.
func (*NoDeprecatedUsageRule) RequiresTypecheck() bool {
	return true
}

// CacheTier returns the cache tier for this rule.
func (*NoDeprecatedUsageRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}
