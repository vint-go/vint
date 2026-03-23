package no_weak_slice_guard

import (
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoWeakSliceGuardRule detects conditions where a nil check on a slice is
// insufficient before indexing. A nil slice and an empty slice both have
// length 0, so `x != nil` does not protect against out-of-bounds panics.
type NoWeakSliceGuardRule struct{}

// Apply applies the rule to given file.
func (r *NoWeakSliceGuardRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintWeakSliceGuard{onFailure: onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoWeakSliceGuardRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintWeakSliceGuard{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoWeakSliceGuardRule) Name() string {
	return "noWeakSliceGuard"
}

// Group returns the rule group.
func (*NoWeakSliceGuardRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoWeakSliceGuardRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintWeakSliceGuard struct {
	onFailure func(lint.Failure)
}

func (w *lintWeakSliceGuard) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	// Pattern 1: x != nil && ...x[i]...
	// Pattern 2: x == nil || ...x[i]...
	switch binExpr.Op {
	case token.LAND: // &&
		w.checkNilGuard(binExpr, binExpr.X, binExpr.Y, token.NEQ)
	case token.LOR: // ||
		w.checkNilGuard(binExpr, binExpr.X, binExpr.Y, token.EQL)
	}

	return w
}

// checkNilGuard checks if the left side is a nil comparison with the given
// operator and the right side contains an index expression on the same variable.
func (w *lintWeakSliceGuard) checkNilGuard(node ast.Node, left, right ast.Expr, expectedOp token.Token) {
	cmp, ok := left.(*ast.BinaryExpr)
	if !ok {
		return
	}

	if cmp.Op != expectedOp {
		return
	}

	varName := extractNilComparedIdent(cmp)
	if varName == "" {
		return
	}

	if containsIndexExpr(right, varName) {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       node,
			Failure:    "nil check for '" + varName + "' is not enough, use len(" + varName + ") != 0",
		})
	}
}

// extractNilComparedIdent returns the name of the identifier compared to nil,
// or empty string if the expression is not a nil comparison.
func extractNilComparedIdent(binExpr *ast.BinaryExpr) string {
	if isNilIdent(binExpr.Y) {
		if ident, ok := binExpr.X.(*ast.Ident); ok {
			return ident.Name
		}
	}
	if isNilIdent(binExpr.X) {
		if ident, ok := binExpr.Y.(*ast.Ident); ok {
			return ident.Name
		}
	}
	return ""
}

// isNilIdent checks if the expression is the nil identifier.
func isNilIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "nil"
}

// containsIndexExpr checks whether the expression tree contains an index
// expression where the indexed variable matches varName.
func containsIndexExpr(expr ast.Expr, varName string) bool {
	found := false
	ast.Inspect(expr, func(n ast.Node) bool {
		if found {
			return false
		}
		indexExpr, ok := n.(*ast.IndexExpr)
		if !ok {
			return true
		}
		if ident, ok := indexExpr.X.(*ast.Ident); ok && ident.Name == varName {
			found = true
			return false
		}
		return true
	})
	return found
}
