package no_unused_function

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoUnusedFunctionRule detects functions that are declared but never called
// or referenced anywhere in the package.
type NoUnusedFunctionRule struct{}

// Apply applies the rule to given file.
func (r *NoUnusedFunctionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	pkg := file.Pkg
	isMain := pkg.IsMain()

	// Phase 1: Collect all function declarations across all files in the package.
	entryPoints := map[string]bool{}
	allFuncs := map[string]bool{}

	for _, f := range pkg.Files() {
		collectFunctions(f, isMain, allFuncs, entryPoints)
	}

	// Phase 2: Build a call graph: for each function, which other functions does it reference?
	callGraph := map[string]map[string]bool{}
	for _, f := range pkg.Files() {
		buildCallGraph(f, allFuncs, callGraph)
	}

	// Phase 3: Compute reachability from entry points using BFS.
	reachable := computeReachable(entryPoints, callGraph)

	// Phase 4: Report functions declared in this file that are not reachable.
	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Recv != nil {
			continue
		}

		name := funcDecl.Name.Name

		// Skip if this is an entry point.
		if entryPoints[name] {
			continue
		}

		// Skip if reachable from an entry point.
		if reachable[name] {
			continue
		}

		// Skip if not in our known function set (shouldn't happen, but be safe).
		if !allFuncs[name] {
			continue
		}

		failures = append(failures, lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       funcDecl,
			Failure:    fmt.Sprintf("func %s is unused", name),
		})
	}

	return failures
}

// methodKey returns the call graph key for a method: "TypeName.MethodName".
func methodKey(funcDecl *ast.FuncDecl) string {
	if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
		return funcDecl.Name.Name
	}
	recvType := funcDecl.Recv.List[0].Type
	// Unwrap pointer receiver.
	if star, ok := recvType.(*ast.StarExpr); ok {
		recvType = star.X
	}
	if ident, ok := recvType.(*ast.Ident); ok {
		return ident.Name + "." + funcDecl.Name.Name
	}
	return funcDecl.Name.Name
}

// collectFunctions scans a file for function declarations and classifies them
// as entry points or regular functions. Methods are also tracked: exported
// methods are entry points and all method names are recorded so their bodies
// can be walked for call-graph edges.
func collectFunctions(file *lint.File, isMain bool, allFuncs map[string]bool, entryPoints map[string]bool) {
	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if funcDecl.Recv != nil {
			// Track methods: exported methods are entry points.
			key := methodKey(funcDecl)
			allFuncs[key] = true
			if ast.IsExported(funcDecl.Name.Name) {
				entryPoints[key] = true
			}
			continue
		}

		name := funcDecl.Name.Name
		allFuncs[name] = true

		// init functions are always considered used.
		if name == "init" {
			entryPoints[name] = true
			continue
		}

		// main function in main package is always used.
		if name == "main" && isMain {
			entryPoints[name] = true
			continue
		}

		// Exported functions are considered used (the package exports them).
		if ast.IsExported(name) {
			entryPoints[name] = true
			continue
		}

		// Check for //export (cgo) and //go:linkname directives.
		if hasCgoExportOrLinkname(funcDecl) {
			entryPoints[name] = true
			continue
		}
	}
}

// hasCgoExportOrLinkname checks if a function has //export or //go:linkname
// directives in its doc comments.
func hasCgoExportOrLinkname(funcDecl *ast.FuncDecl) bool {
	if funcDecl.Doc != nil {
		for _, comment := range funcDecl.Doc.List {
			text := comment.Text
			if strings.HasPrefix(text, "//export ") || strings.HasPrefix(text, "//go:linkname ") {
				return true
			}
		}
	}
	return false
}

// buildCallGraph walks each function/method body and records which other
// package-level functions it references (by name). It also walks method bodies
// so that method-to-function call edges are captured.
func buildCallGraph(file *lint.File, allFuncs map[string]bool, callGraph map[string]map[string]bool) {
	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}

		callerName := methodKey(funcDecl)
		if callGraph[callerName] == nil {
			callGraph[callerName] = map[string]bool{}
		}

		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			ident, ok := n.(*ast.Ident)
			if !ok {
				return true
			}
			// Check if this identifier refers to a known function in the package.
			name := ident.Name
			if name == callerName {
				return true // skip self-references
			}
			// Use the set of known function names instead of ident.Obj, which
			// is nil for cross-file references (parser resolves per-file only).
			if allFuncs[name] {
				callGraph[callerName][name] = true
			}
			return true
		})
	}

	// Also track references from package-level variable initializers.
	for _, decl := range file.AST.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		ast.Inspect(genDecl, func(n ast.Node) bool {
			ident, ok := n.(*ast.Ident)
			if !ok {
				return true
			}
			if allFuncs[ident.Name] {
				// Package-level init references count as entry point references.
				// We model this by adding these to a synthetic "__pkg_init__" entry.
				const pkgInit = "__pkg_init__"
				if callGraph[pkgInit] == nil {
					callGraph[pkgInit] = map[string]bool{}
				}
				callGraph[pkgInit][ident.Name] = true
			}
			return true
		})
	}
}

// computeReachable performs a BFS from all entry points through the call graph
// and returns the set of all reachable function names.
func computeReachable(entryPoints map[string]bool, callGraph map[string]map[string]bool) map[string]bool {
	reachable := map[string]bool{}
	queue := []string{}

	// Seed the BFS with all entry points.
	for name := range entryPoints {
		queue = append(queue, name)
		reachable[name] = true
	}

	// Also add the synthetic package init node.
	const pkgInit = "__pkg_init__"
	if _, exists := callGraph[pkgInit]; exists {
		queue = append(queue, pkgInit)
		reachable[pkgInit] = true
	}

	// BFS traversal.
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for callee := range callGraph[current] {
			if !reachable[callee] {
				reachable[callee] = true
				queue = append(queue, callee)
			}
		}
	}

	return reachable
}

// Name returns the rule name.
func (*NoUnusedFunctionRule) Name() string {
	return "noUnusedFunction"
}

// Group returns the rule group.
func (*NoUnusedFunctionRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnusedFunctionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}
