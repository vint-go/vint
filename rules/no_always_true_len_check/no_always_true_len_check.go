package no_always_true_len_check

import (
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoAlwaysTrueLenCheckRule detects usage of len where the result is obvious,
// such as len(s) >= 0 (always true) or len(s) < 0 (always false).
type NoAlwaysTrueLenCheckRule struct{}

// Apply applies the rule to given file.
func (r *NoAlwaysTrueLenCheckRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintAlwaysTrueLen{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoAlwaysTrueLenCheckRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintAlwaysTrueLen{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoAlwaysTrueLenCheckRule) Name() string {
	return "noAlwaysTrueLenCheck"
}

// Group returns the rule group.
func (*NoAlwaysTrueLenCheckRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoAlwaysTrueLenCheckRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintAlwaysTrueLen struct {
	onFailure func(lint.Failure)
}

func (w *lintAlwaysTrueLen) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	// Check for patterns: len(x) >= 0, len(x) < 0, 0 <= len(x), 0 > len(x)
	switch {
	case isLenCall(binExpr.X) && isZeroLiteral(binExpr.Y):
		w.checkLenOpZero(binExpr, binExpr.Op)
	case isZeroLiteral(binExpr.X) && isLenCall(binExpr.Y):
		// Flip the operator: 0 <= len(x) is the same as len(x) >= 0
		w.checkLenOpZero(binExpr, flipOp(binExpr.Op))
	}

	return w
}

func (w *lintAlwaysTrueLen) checkLenOpZero(node ast.Node, op token.Token) {
	switch op {
	case token.GEQ: // len(x) >= 0 is always true
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       node,
			Failure:    "len(x) >= 0 is always true",
		})
	case token.LSS: // len(x) < 0 is always false
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       node,
			Failure:    "len(x) < 0 is always false",
		})
	}
}

// isLenCall checks if the expression is a call to the builtin len function.
func isLenCall(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}
	ident, ok := call.Fun.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == "len" && len(call.Args) == 1
}

// isZeroLiteral checks if the expression is the integer literal 0.
func isZeroLiteral(expr ast.Expr) bool {
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		return false
	}
	return lit.Kind == token.INT && lit.Value == "0"
}

// flipOp reverses the direction of a comparison operator.
func flipOp(op token.Token) token.Token {
	switch op {
	case token.LEQ:
		return token.GEQ
	case token.GEQ:
		return token.LEQ
	case token.LSS:
		return token.GTR
	case token.GTR:
		return token.LSS
	default:
		return op
	}
}
