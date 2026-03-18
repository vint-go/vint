package use_strings_contains

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseStringsContainsRule detects comparisons of strings.Index results
// against -1 or 0 that can be replaced with strings.Contains.
type UseStringsContainsRule struct{}

// Apply applies the rule to given file.
func (r *UseStringsContainsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintStringsContains{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseStringsContainsRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintStringsContains{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseStringsContainsRule) Name() string {
	return "useStringsContains"
}

// Group returns the rule group.
func (*UseStringsContainsRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseStringsContainsRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintStringsContains struct {
	onFailure func(lint.Failure)
}

func (w *lintStringsContains) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	// Check both orientations: Index(...) op literal, and literal op Index(...)
	if w.checkIndexComparison(binExpr, binExpr.X, binExpr.Y, binExpr.Op) {
		return w
	}
	// Reversed: literal op Index(...)
	w.checkIndexComparison(binExpr, binExpr.Y, binExpr.X, reverseOp(binExpr.Op))

	return w
}

// checkIndexComparison checks if callSide is strings.Index(...) and litSide is a numeric literal,
// combined with the given operator, forming a pattern replaceable by strings.Contains.
// Returns true if a failure was reported.
func (w *lintStringsContains) checkIndexComparison(node ast.Node, callSide, litSide ast.Expr, op token.Token) bool {
	call, ok := callSide.(*ast.CallExpr)
	if !ok {
		return false
	}

	if !astutils.IsPkgDotName(call.Fun, "strings", "Index") {
		return false
	}

	litVal, ok := intLitValue(litSide)
	if !ok {
		return false
	}

	// Detect the patterns:
	//   strings.Index(s, sub) != -1  => strings.Contains
	//   strings.Index(s, sub) >= 0   => strings.Contains
	//   strings.Index(s, sub) > -1   => strings.Contains
	//   strings.Index(s, sub) == -1  => !strings.Contains
	//   strings.Index(s, sub) < 0    => !strings.Contains
	match := false
	switch {
	case op == token.NEQ && litVal == -1:
		match = true
	case op == token.GEQ && litVal == 0:
		match = true
	case op == token.GTR && litVal == -1:
		match = true
	case op == token.EQL && litVal == -1:
		match = true
	case op == token.LSS && litVal == 0:
		match = true
	}

	if !match {
		return false
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       node,
		Category:   lint.FailureCategoryStyle,
		Failure:    "replace call to strings.Index with strings.Contains",
	})

	return true
}

// intLitValue extracts the integer value from an ast expression if it is
// a basic integer literal or a unary minus applied to one.
func intLitValue(expr ast.Expr) (int, bool) {
	// Handle negative: -1 is represented as UnaryExpr{Op: -, X: BasicLit{1}}
	if unary, ok := expr.(*ast.UnaryExpr); ok && unary.Op == token.SUB {
		lit, ok := unary.X.(*ast.BasicLit)
		if !ok || lit.Kind != token.INT {
			return 0, false
		}
		if lit.Value == "1" {
			return -1, true
		}
		return 0, false
	}

	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.INT {
		return 0, false
	}

	switch lit.Value {
	case "0":
		return 0, true
	case "1":
		return 1, true
	default:
		return 0, false
	}
}

// reverseOp returns the comparison operator with operands swapped.
func reverseOp(op token.Token) token.Token {
	switch op {
	case token.LSS:
		return token.GTR
	case token.GTR:
		return token.LSS
	case token.LEQ:
		return token.GEQ
	case token.GEQ:
		return token.LEQ
	default:
		return op // EQL, NEQ are symmetric
	}
}
