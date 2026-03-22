package use_errors_as

import (
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseErrorsAsRule detects type assertions and type switches on error values
// that should use errors.As() instead.
type UseErrorsAsRule struct{}

// Apply applies the rule to given file.
func (*UseErrorsAsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUseErrorsAs{file: file, onFailure: onFailure}
	if w.file.Pkg.TypeCheck() != nil {
		return nil
	}

	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (*UseErrorsAsRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintUseErrorsAs{file: file, onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseErrorsAsRule) Name() string {
	return "useErrorsAs"
}

// Group returns the rule group.
func (*UseErrorsAsRule) Group() string {
	return "correctness"
}

// RequiresTypecheck returns true because this rule uses type information.
func (*UseErrorsAsRule) RequiresTypecheck() bool { return true }

// CacheTier returns the cache tier for this rule.
func (*UseErrorsAsRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

type lintUseErrorsAs struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintUseErrorsAs) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.TypeAssertExpr:
		// Skip type assertions without an explicit type (e.g. x.(type) in switch)
		if n.Type == nil {
			return w
		}
		if w.isErrorExpr(n.X) {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryErrors,
				Confidence: 1,
				Node:       n,
				Failure:    "type assertion on error will fail on wrapped errors, use errors.As",
			})
		}
	case *ast.TypeSwitchStmt:
		var switchExpr ast.Expr
		if assign, ok := n.Assign.(*ast.AssignStmt); ok && len(assign.Rhs) == 1 {
			if tae, ok := assign.Rhs[0].(*ast.TypeAssertExpr); ok {
				switchExpr = tae.X
			}
		} else if exprStmt, ok := n.Assign.(*ast.ExprStmt); ok {
			if tae, ok := exprStmt.X.(*ast.TypeAssertExpr); ok {
				switchExpr = tae.X
			}
		}

		if switchExpr != nil && w.isErrorExpr(switchExpr) {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryErrors,
				Confidence: 1,
				Node:       n,
				Failure:    "type switch on error will fail on wrapped errors, use errors.As",
			})
		}
	}
	return w
}

func (w *lintUseErrorsAs) isErrorExpr(expr ast.Expr) bool {
	typ := w.file.Pkg.TypeOf(expr)
	if typ == nil {
		return false
	}
	errorIface := types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
	return types.Implements(typ, errorIface)
}
