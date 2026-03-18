package no_external_error_reassign

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoExternalErrorReassignRule detects suspicious reassignment of exported error
// variables from other packages. Package-level error sentinel values (like
// io.EOF or http.ErrServerClosed) should be compared against, not reassigned.
type NoExternalErrorReassignRule struct{}

// Apply applies the rule to given file.
func (r *NoExternalErrorReassignRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoExternalErrorReassign{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoExternalErrorReassignRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoExternalErrorReassign{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoExternalErrorReassignRule) Name() string {
	return "noExternalErrorReassign"
}

// Group returns the rule group.
func (*NoExternalErrorReassignRule) Group() string {
	return "correctness"
}

// RequiresTypecheck returns true because this rule needs type information.
func (*NoExternalErrorReassignRule) RequiresTypecheck() bool {
	return true
}

// CacheTier returns the cache tier for this rule.
func (*NoExternalErrorReassignRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

type lintNoExternalErrorReassign struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintNoExternalErrorReassign) Visit(node ast.Node) ast.Visitor {
	assign, ok := node.(*ast.AssignStmt)
	if !ok {
		return w
	}

	for _, lhs := range assign.Lhs {
		sel, ok := lhs.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		// Check that the selector's X is an identifier (package name)
		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			continue
		}

		// Use type info to check that the identifier refers to a package
		info := w.pkg.TypesInfo()
		if info == nil {
			continue
		}

		obj := info.ObjectOf(ident)
		if obj == nil {
			continue
		}

		_, isPkg := obj.(*types.PkgName)
		if !isPkg {
			continue
		}

		// Check that the selected field/var implements the error interface
		selObj := info.ObjectOf(sel.Sel)
		if selObj == nil {
			continue
		}

		if !isErrorType(selObj.Type()) {
			continue
		}

		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryErrors,
			Confidence: 1,
			Node:       assign,
			Failure:    fmt.Sprintf("reassigning external error variable %s.%s", ident.Name, sel.Sel.Name),
		})
	}

	return w
}

// isErrorType checks whether the given type implements the error interface.
func isErrorType(t types.Type) bool {
	if t == nil {
		return false
	}

	// Check if the type is exactly the error interface
	errorType := types.Universe.Lookup("error").Type()
	if errorType == nil {
		return false
	}

	return types.Implements(t, errorType.Underlying().(*types.Interface))
}
