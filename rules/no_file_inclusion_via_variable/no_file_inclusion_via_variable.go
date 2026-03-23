package no_file_inclusion_via_variable

import (
	"fmt"
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoFileInclusionViaVariableRule detects potential file inclusion via variable,
// where file paths provided as taint input are used to open or read files.
type NoFileInclusionViaVariableRule struct{}

// Apply applies the rule to given file.
func (r *NoFileInclusionViaVariableRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoFileInclusionViaVariable{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoFileInclusionViaVariableRule) Name() string {
	return "noFileInclusionViaVariable"
}

// Group returns the rule group.
func (*NoFileInclusionViaVariableRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoFileInclusionViaVariableRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoFileInclusionViaVariable struct {
	onFailure func(lint.Failure)
}

// fileFuncPathArgIndex maps os and ioutil package functions to the
// zero-based index of their file path argument.
var fileFuncPathArgIndex = map[string]map[string]int{
	"os": {
		"Open":     0, // os.Open(name)
		"OpenFile": 0, // os.OpenFile(name, flag, perm)
		"ReadFile": 0, // os.ReadFile(name)
		"Create":   0, // os.Create(name)
	},
	"ioutil": {
		"ReadFile": 0, // ioutil.ReadFile(filename)
	},
}

// sanitizationFuncs lists functions from path/filepath that sanitize paths.
var sanitizationFuncs = map[string]bool{
	"Clean":        true,
	"Rel":          true,
	"EvalSymlinks": true,
}

func (w *lintNoFileInclusionViaVariable) Visit(node ast.Node) ast.Visitor {
	// We process at the function declaration level to track local variable assignments.
	fd, ok := node.(*ast.FuncDecl)
	if !ok {
		return w
	}

	if fd.Body == nil {
		return nil
	}

	// Collect variables assigned from sanitization function results within this function.
	sanitizedVars := collectSanitizedVars(fd.Body)

	// Walk the function body looking for file operation calls.
	fw := &funcWalker{
		onFailure:     w.onFailure,
		sanitizedVars: sanitizedVars,
	}
	ast.Walk(fw, fd.Body)

	return nil // don't recurse into the function again
}

// funcWalker walks a function body to find file operation calls.
type funcWalker struct {
	onFailure     func(lint.Failure)
	sanitizedVars map[string]bool
}

func (fw *funcWalker) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return fw
	}

	for pkg, funcs := range fileFuncPathArgIndex {
		for funcName, argIdx := range funcs {
			if astutils.IsPkgDotName(ce.Fun, pkg, funcName) {
				if len(ce.Args) > argIdx {
					fw.checkPathArg(ce, ce.Args[argIdx], pkg, funcName)
				}
				return fw
			}
		}
	}

	return fw
}

// checkPathArg reports a failure if the file path argument is not a
// string literal constant and has not been sanitized by a filepath function.
func (fw *funcWalker) checkPathArg(call *ast.CallExpr, pathArg ast.Expr, pkg, funcName string) {
	if astutils.IsStringLiteral(pathArg) {
		return // hardcoded string literal is safe
	}

	if isSanitizedExpr(pathArg) {
		return // direct call to sanitization function
	}

	// Check if the argument is a variable that was assigned from a sanitization function.
	if ident, ok := pathArg.(*ast.Ident); ok {
		if fw.sanitizedVars[ident.Name] {
			return // variable was assigned from sanitization function
		}
	}

	fw.onFailure(lint.Failure{
		Confidence: 1,
		Node:       call,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    fmt.Sprintf("potential file inclusion via variable: path passed to %s.%s is not a hardcoded constant", pkg, funcName),
	})
}

// collectSanitizedVars scans a block statement for variable assignments
// where the right-hand side is a call to a path/filepath sanitization function.
// It returns the set of variable names that hold sanitized values.
func collectSanitizedVars(body *ast.BlockStmt) map[string]bool {
	sanitized := map[string]bool{}

	ast.Inspect(body, func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.AssignStmt:
			// Check short variable declarations and assignments: x := filepath.Clean(...)
			for i, rhs := range stmt.Rhs {
				if isSanitizedExpr(rhs) && i < len(stmt.Lhs) {
					if ident, ok := stmt.Lhs[i].(*ast.Ident); ok {
						sanitized[ident.Name] = true
					}
				}
				// Also handle multi-return: x, err := filepath.Rel(...)
				if call, ok := rhs.(*ast.CallExpr); ok {
					if isSanitizedCall(call) {
						// First return value is the sanitized path
						if i < len(stmt.Lhs) {
							if ident, ok := stmt.Lhs[i].(*ast.Ident); ok {
								sanitized[ident.Name] = true
							}
						}
					}
				}
			}
		case *ast.ValueSpec:
			// var x = filepath.Clean(...)
			for i, val := range stmt.Values {
				if isSanitizedExpr(val) && i < len(stmt.Names) {
					sanitized[stmt.Names[i].Name] = true
				}
			}
		}
		return true
	})

	return sanitized
}

// isSanitizedExpr checks whether the given expression is a call to a
// path/filepath sanitization function (Clean, Rel, EvalSymlinks).
func isSanitizedExpr(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}
	return isSanitizedCall(call)
}

// isSanitizedCall checks whether the given call expression is a call to a
// path/filepath sanitization function.
func isSanitizedCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "filepath" && sanitizationFuncs[sel.Sel.Name]
}
