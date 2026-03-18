package use_equal_fold

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseEqualFoldRule detects case-insensitive string comparisons that use
// strings.ToLower or strings.ToUpper before comparing, instead of the
// more efficient strings.EqualFold.
type UseEqualFoldRule struct{}

// Apply applies the rule to given file.
func (r *UseEqualFoldRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintEqualFold{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseEqualFoldRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintEqualFold{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseEqualFoldRule) Name() string {
	return "useEqualFold"
}

// Group returns the rule group.
func (*UseEqualFoldRule) Group() string {
	return "performance"
}

// CacheTier returns the cache tier for this rule.
func (*UseEqualFoldRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintEqualFold struct {
	onFailure func(lint.Failure)
}

func (w *lintEqualFold) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	if binExpr.Op != token.EQL && binExpr.Op != token.NEQ {
		return w
	}

	// Check if at least one side is strings.ToLower or strings.ToUpper
	leftIsCaseConv := isCaseConversionCall(binExpr.X)
	rightIsCaseConv := isCaseConversionCall(binExpr.Y)

	if !leftIsCaseConv && !rightIsCaseConv {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       node,
		Category:   lint.FailureCategoryOptimization,
		Failure:    "use strings.EqualFold instead of strings.ToLower or strings.ToUpper for case-insensitive comparison",
	})

	return w
}

// isCaseConversionCall checks if the expression is a call to
// strings.ToLower or strings.ToUpper.
func isCaseConversionCall(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	return astutils.IsPkgDotName(call.Fun, "strings", "ToLower") ||
		astutils.IsPkgDotName(call.Fun, "strings", "ToUpper")
}
