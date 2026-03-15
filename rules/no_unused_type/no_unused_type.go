package no_unused_type

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoUnusedTypeRule detects named types that are declared but never used
// anywhere in the codebase.
type NoUnusedTypeRule struct {
	generatedIsUsed bool
}

const defaultGeneratedIsUsed = true

// Configure validates and applies the rule configuration.
func (r *NoUnusedTypeRule) Configure(arguments lint.Arguments) error {
	r.generatedIsUsed = defaultGeneratedIsUsed

	if len(arguments) == 0 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noUnusedType" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		switch strings.ReplaceAll(k, "_", "-") {
		case "generated-is-used":
			b, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid value for "generated-is-used" in "noUnusedType"; need bool but got %T`, v)
			}
			r.generatedIsUsed = b
		}
	}

	return nil
}

// pkgTypeInfo holds information about a package-level type declaration.
type pkgTypeInfo struct {
	name     string
	node     ast.Node // the *ast.Ident of the type name
	fileName string
}

// Apply applies the rule to the given file.
func (r *NoUnusedTypeRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	pkg := file.Pkg

	// Phase 1: Collect all package-level type declarations across all files.
	allTypes := map[string]*pkgTypeInfo{}
	entryTypes := map[string]bool{} // types considered always used

	for fname, f := range pkg.Files() {
		// If generated-is-used is enabled and file is generated, mark its types as entries.
		if r.generatedIsUsed && isGeneratedFile(f) {
			for _, decl := range f.AST.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}
				for _, spec := range genDecl.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					entryTypes[ts.Name.Name] = true
				}
			}
			continue
		}

		for _, decl := range f.AST.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}
			for _, spec := range genDecl.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				name := ts.Name.Name

				// Blank identifier is always considered used.
				if name == "_" {
					continue
				}

				// Exported types are always considered used.
				if ast.IsExported(name) {
					entryTypes[name] = true
					continue
				}

				allTypes[name] = &pkgTypeInfo{
					name:     name,
					node:     ts.Name,
					fileName: fname,
				}
			}
		}
	}

	if len(allTypes) == 0 {
		return nil
	}

	// Phase 2: Scan for type references from "entry point" contexts:
	// - Function declarations (signatures + bodies)
	// - Exported var/const type annotations and values
	// Types referenced from unexported var/const or other type declarations
	// only count if the referencing entity is itself used.
	directlyUsed := map[string]bool{}

	// Also build a dependency graph: type A depends on type B if A's definition
	// references B.
	typeDeps := map[string]map[string]bool{}

	for _, f := range pkg.Files() {
		for _, decl := range f.AST.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				// Function signatures and bodies are entry points for type usage.
				collectTypeRefsFromFuncDecl(d, allTypes, directlyUsed)
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.ValueSpec:
						// Only count references from exported var/const as direct usage.
						hasExportedName := false
						for _, nameIdent := range s.Names {
							if ast.IsExported(nameIdent.Name) {
								hasExportedName = true
								break
							}
						}
						if hasExportedName {
							if s.Type != nil {
								collectExprTypeRefs(s.Type, allTypes, directlyUsed, nil)
							}
							for _, val := range s.Values {
								collectExprTypeRefs(val, allTypes, directlyUsed, nil)
							}
						}
					case *ast.TypeSpec:
						// Build type dependency graph.
						typeName := s.Name.Name
						deps := map[string]bool{}
						skipNames := map[string]bool{typeName: true}
						collectExprTypeRefs(s.Type, allTypes, deps, skipNames)
						if s.TypeParams != nil {
							for _, field := range s.TypeParams.List {
								if field.Type != nil {
									collectExprTypeRefs(field.Type, allTypes, deps, skipNames)
								}
							}
						}
						if len(deps) > 0 {
							typeDeps[typeName] = deps
						}

						// If this type is exported or an entry, its deps are directly used.
						if entryTypes[typeName] {
							for dep := range deps {
								directlyUsed[dep] = true
							}
						}
					}
				}
			}
		}
	}

	// Phase 2.5: Propagate usage through type dependency graph.
	// If type A is used (directly or transitively) and A's definition
	// references type B, then B is also used.
	reachable := computeReachable(directlyUsed, entryTypes, typeDeps)

	// Phase 3: Report unused types declared in this file.
	thisFileName := fileNameFromFile(file)

	for _, decl := range file.AST.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			name := ts.Name.Name
			if name == "_" {
				continue
			}
			if ast.IsExported(name) {
				continue
			}
			if entryTypes[name] {
				continue
			}

			info, exists := allTypes[name]
			if !exists {
				continue
			}
			// Only report if declared in this file.
			if info.fileName != thisFileName {
				continue
			}

			if reachable[name] {
				continue
			}

			failures = append(failures, lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 1,
				Node:       ts.Name,
				Failure:    fmt.Sprintf("type %s is unused", name),
			})
		}
	}

	return failures
}

// collectTypeRefsFromFuncDecl collects type references from a function declaration
// (signatures and body).
func collectTypeRefsFromFuncDecl(funcDecl *ast.FuncDecl, allTypes map[string]*pkgTypeInfo, refs map[string]bool) {
	// Scan receiver.
	if funcDecl.Recv != nil {
		for _, field := range funcDecl.Recv.List {
			if field.Type != nil {
				collectExprTypeRefs(field.Type, allTypes, refs, nil)
			}
		}
	}
	// Scan parameters.
	if funcDecl.Type.Params != nil {
		for _, field := range funcDecl.Type.Params.List {
			if field.Type != nil {
				collectExprTypeRefs(field.Type, allTypes, refs, nil)
			}
		}
	}
	// Scan results.
	if funcDecl.Type.Results != nil {
		for _, field := range funcDecl.Type.Results.List {
			if field.Type != nil {
				collectExprTypeRefs(field.Type, allTypes, refs, nil)
			}
		}
	}
	// Scan type params (generics).
	if funcDecl.Type.TypeParams != nil {
		for _, field := range funcDecl.Type.TypeParams.List {
			if field.Type != nil {
				collectExprTypeRefs(field.Type, allTypes, refs, nil)
			}
		}
	}
	// Scan function body.
	if funcDecl.Body != nil {
		scanBodyForTypeRefs(funcDecl.Body, allTypes, refs)
	}
}

// scanBodyForTypeRefs walks a function body looking for type references.
func scanBodyForTypeRefs(node ast.Node, allTypes map[string]*pkgTypeInfo, refs map[string]bool) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			if _, isType := allTypes[x.Name]; isType {
				if x.Obj != nil && x.Obj.Kind == ast.Typ {
					refs[x.Name] = true
				}
			}
		case *ast.CompositeLit:
			// Check composite literal type: e.g. config{...}
			collectExprTypeRefs(x.Type, allTypes, refs, nil)
		case *ast.TypeAssertExpr:
			// Check type assertion: e.g. x.(config)
			if x.Type != nil {
				collectExprTypeRefs(x.Type, allTypes, refs, nil)
			}
		}
		return true
	})
}

// collectExprTypeRefs scans an expression for type identifier references.
// skipNames is a set of type names to skip (e.g. the type being declared).
func collectExprTypeRefs(expr ast.Expr, allTypes map[string]*pkgTypeInfo, refs map[string]bool, skipNames map[string]bool) {
	if expr == nil {
		return
	}
	ast.Inspect(expr, func(n ast.Node) bool {
		ident, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		name := ident.Name
		if skipNames != nil && skipNames[name] {
			return true
		}
		if _, isType := allTypes[name]; isType {
			refs[name] = true
		}
		return true
	})
}

// computeReachable determines which types are reachable from entry points
// through the type dependency graph.
func computeReachable(directlyUsed map[string]bool, entryTypes map[string]bool, typeDeps map[string]map[string]bool) map[string]bool {
	reachable := map[string]bool{}
	queue := []string{}

	// Seed with all directly used types.
	for name := range directlyUsed {
		if !reachable[name] {
			reachable[name] = true
			queue = append(queue, name)
		}
	}

	// BFS: if a type is used, all types it depends on are also used.
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for dep := range typeDeps[current] {
			if !reachable[dep] {
				reachable[dep] = true
				queue = append(queue, dep)
			}
		}
	}

	return reachable
}

// fileNameFromFile extracts the file name that matches pkg.Files() keys.
func fileNameFromFile(file *lint.File) string {
	for name, f := range file.Pkg.Files() {
		if f == file {
			return name
		}
	}
	return ""
}

// isGeneratedFile checks if a file has a "Code generated" comment indicating
// it was automatically generated.
func isGeneratedFile(file *lint.File) bool {
	for _, cg := range file.AST.Comments {
		for _, c := range cg.List {
			if strings.Contains(c.Text, "Code generated") && strings.Contains(c.Text, "DO NOT EDIT") {
				return true
			}
		}
	}
	return false
}

// Name returns the rule name.
func (*NoUnusedTypeRule) Name() string {
	return "noUnusedType"
}

// Group returns the rule group.
func (*NoUnusedTypeRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnusedTypeRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}
