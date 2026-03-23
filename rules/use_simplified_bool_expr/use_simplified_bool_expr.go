package use_simplified_bool_expr

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseSimplifiedBoolExprRule detects bool expressions that can be simplified.
type UseSimplifiedBoolExprRule struct{}

// Apply applies the rule to given file.
func (r *UseSimplifiedBoolExprRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintSimplifiedBool{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseSimplifiedBoolExprRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintSimplifiedBool{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseSimplifiedBoolExprRule) Name() string {
	return "useSimplifiedBoolExpr"
}

// Group returns the rule group.
func (*UseSimplifiedBoolExprRule) Group() string {
	return "complexity"
}

// CacheTier returns the cache tier for this rule.
func (*UseSimplifiedBoolExprRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintSimplifiedBool struct {
	onFailure func(lint.Failure)
}

func (w *lintSimplifiedBool) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.UnaryExpr:
		w.checkUnaryExpr(n)
	case *ast.BinaryExpr:
		w.checkBinaryExpr(n)
	}
	return w
}

// checkUnaryExpr handles patterns involving negation (!)
func (w *lintSimplifiedBool) checkUnaryExpr(expr *ast.UnaryExpr) {
	if expr.Op != token.NOT {
		return
	}

	inner := unwrapParens(expr.X)

	// Pattern 1: Double negation !(!(x)) -> x
	if innerUnary, ok := inner.(*ast.UnaryExpr); ok && innerUnary.Op == token.NOT {
		innerExpr := unwrapParens(innerUnary.X)
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryComplexity,
			Confidence: 1,
			Node:       expr,
			Failure:    fmt.Sprintf("can simplify !(!(%s)) to %s", astutils.GoFmt(innerExpr), astutils.GoFmt(innerExpr)),
		})
		return
	}

	// Pattern 3: Inverted comparisons !(x >= y) -> x < y
	if binExpr, ok := inner.(*ast.BinaryExpr); ok {
		if neg := negateComparisonOp(binExpr.Op); neg != token.ILLEGAL {
			// Skip float-related expressions
			if containsFloatLiteral(binExpr.X) || containsFloatLiteral(binExpr.Y) {
				return
			}
			lhs := astutils.GoFmt(binExpr.X)
			rhs := astutils.GoFmt(binExpr.Y)
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryComplexity,
				Confidence: 1,
				Node:       expr,
				Failure:    fmt.Sprintf("can simplify !(%s %s %s) to %s %s %s", lhs, binExpr.Op, rhs, lhs, neg, rhs),
			})
		}
	}
}

// checkBinaryExpr handles patterns involving binary expressions
func (w *lintSimplifiedBool) checkBinaryExpr(expr *ast.BinaryExpr) {
	// Pattern 2: Negated equality !(a) == !(b) -> (a) == (b)
	w.checkNegatedEquality(expr)

	// Pattern 4: Combined checks: x > c || x == c -> x >= c
	w.checkCombinedComparisons(expr)

	// Pattern 5: Increment/decrement removal: x > y-1 -> x >= y
	w.checkIncrementDecrement(expr)

	// Pattern 6: Range folding: x >= c && x <= c -> x == c
	w.checkRangeFolding(expr)
}

// checkNegatedEquality detects !(a) == !(b) -> (a) == (b)
func (w *lintSimplifiedBool) checkNegatedEquality(expr *ast.BinaryExpr) {
	if expr.Op != token.EQL && expr.Op != token.NEQ {
		return
	}

	lhsUnary, lhsOk := unwrapParens(expr.X).(*ast.UnaryExpr)
	rhsUnary, rhsOk := unwrapParens(expr.Y).(*ast.UnaryExpr)

	if !lhsOk || !rhsOk || lhsUnary.Op != token.NOT || rhsUnary.Op != token.NOT {
		return
	}

	lhs := astutils.GoFmt(lhsUnary.X)
	rhs := astutils.GoFmt(rhsUnary.X)
	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryComplexity,
		Confidence: 1,
		Node:       expr,
		Failure:    fmt.Sprintf("can simplify !(%s) %s !(%s) to (%s) %s (%s)", lhs, expr.Op, rhs, lhs, expr.Op, rhs),
	})
}

// checkCombinedComparisons detects patterns like x > c || x == c -> x >= c
// and x < c || x == c -> x <= c
func (w *lintSimplifiedBool) checkCombinedComparisons(expr *ast.BinaryExpr) {
	if expr.Op != token.LOR {
		return
	}

	lhs, lhsOk := unwrapParens(expr.X).(*ast.BinaryExpr)
	rhs, rhsOk := unwrapParens(expr.Y).(*ast.BinaryExpr)
	if !lhsOk || !rhsOk {
		return
	}

	// Skip float-related expressions
	if containsFloatLiteral(lhs.X) || containsFloatLiteral(lhs.Y) ||
		containsFloatLiteral(rhs.X) || containsFloatLiteral(rhs.Y) {
		return
	}

	// Detect: x > c || x == c -> x >= c
	// Detect: x < c || x == c -> x <= c
	var cmpExpr, eqlExpr *ast.BinaryExpr
	if isStrictComparison(lhs.Op) && rhs.Op == token.EQL {
		cmpExpr, eqlExpr = lhs, rhs
	} else if isStrictComparison(rhs.Op) && lhs.Op == token.EQL {
		cmpExpr, eqlExpr = rhs, lhs
	} else {
		return
	}

	cmpLHS := astutils.GoFmt(cmpExpr.X)
	cmpRHS := astutils.GoFmt(cmpExpr.Y)
	eqlLHS := astutils.GoFmt(eqlExpr.X)
	eqlRHS := astutils.GoFmt(eqlExpr.Y)

	if cmpLHS != eqlLHS || cmpRHS != eqlRHS {
		return
	}

	combined := combineComparisonOp(cmpExpr.Op)
	if combined == token.ILLEGAL {
		return
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryComplexity,
		Confidence: 1,
		Node:       expr,
		Failure:    fmt.Sprintf("can simplify %s %s %s || %s == %s to %s %s %s", cmpLHS, cmpExpr.Op, cmpRHS, eqlLHS, eqlRHS, cmpLHS, combined, cmpRHS),
	})
}

// checkIncrementDecrement detects patterns like x > y-1 -> x >= y and x < y+1 -> x <= y
func (w *lintSimplifiedBool) checkIncrementDecrement(expr *ast.BinaryExpr) {
	// Skip float-related expressions
	if containsFloatLiteral(expr.X) || containsFloatLiteral(expr.Y) {
		return
	}

	// Pattern: x > y-1 -> x >= y (RHS has subtraction of 1)
	if expr.Op == token.GTR {
		if binRHS, ok := unwrapParens(expr.Y).(*ast.BinaryExpr); ok && binRHS.Op == token.SUB {
			if isIntLiteral(binRHS.Y, "1") {
				lhs := astutils.GoFmt(expr.X)
				rhs := astutils.GoFmt(binRHS.X)
				w.onFailure(lint.Failure{
					Category:   lint.FailureCategoryComplexity,
					Confidence: 1,
					Node:       expr,
					Failure:    fmt.Sprintf("can simplify %s > %s-1 to %s >= %s", lhs, rhs, lhs, rhs),
				})
			}
		}
	}

	// Pattern: x < y+1 -> x <= y (RHS has addition of 1)
	if expr.Op == token.LSS {
		if binRHS, ok := unwrapParens(expr.Y).(*ast.BinaryExpr); ok && binRHS.Op == token.ADD {
			if isIntLiteral(binRHS.Y, "1") {
				lhs := astutils.GoFmt(expr.X)
				rhs := astutils.GoFmt(binRHS.X)
				w.onFailure(lint.Failure{
					Category:   lint.FailureCategoryComplexity,
					Confidence: 1,
					Node:       expr,
					Failure:    fmt.Sprintf("can simplify %s < %s+1 to %s <= %s", lhs, rhs, lhs, rhs),
				})
			}
		}
	}

	// Pattern: x >= y+1 -> x > y (RHS has addition of 1)
	if expr.Op == token.GEQ {
		if binRHS, ok := unwrapParens(expr.Y).(*ast.BinaryExpr); ok && binRHS.Op == token.ADD {
			if isIntLiteral(binRHS.Y, "1") {
				lhs := astutils.GoFmt(expr.X)
				rhs := astutils.GoFmt(binRHS.X)
				w.onFailure(lint.Failure{
					Category:   lint.FailureCategoryComplexity,
					Confidence: 1,
					Node:       expr,
					Failure:    fmt.Sprintf("can simplify %s >= %s+1 to %s > %s", lhs, rhs, lhs, rhs),
				})
			}
		}
	}

	// Pattern: x <= y-1 -> x < y (RHS has subtraction of 1)
	if expr.Op == token.LEQ {
		if binRHS, ok := unwrapParens(expr.Y).(*ast.BinaryExpr); ok && binRHS.Op == token.SUB {
			if isIntLiteral(binRHS.Y, "1") {
				lhs := astutils.GoFmt(expr.X)
				rhs := astutils.GoFmt(binRHS.X)
				w.onFailure(lint.Failure{
					Category:   lint.FailureCategoryComplexity,
					Confidence: 1,
					Node:       expr,
					Failure:    fmt.Sprintf("can simplify %s <= %s-1 to %s < %s", lhs, rhs, lhs, rhs),
				})
			}
		}
	}
}

// checkRangeFolding detects x >= c && x <= c -> x == c
func (w *lintSimplifiedBool) checkRangeFolding(expr *ast.BinaryExpr) {
	if expr.Op != token.LAND {
		return
	}

	lhs, lhsOk := unwrapParens(expr.X).(*ast.BinaryExpr)
	rhs, rhsOk := unwrapParens(expr.Y).(*ast.BinaryExpr)
	if !lhsOk || !rhsOk {
		return
	}

	// Skip float-related expressions
	if containsFloatLiteral(lhs.X) || containsFloatLiteral(lhs.Y) ||
		containsFloatLiteral(rhs.X) || containsFloatLiteral(rhs.Y) {
		return
	}

	// Check: x >= c && x <= c  or  x <= c && x >= c
	var geqExpr, leqExpr *ast.BinaryExpr
	if lhs.Op == token.GEQ && rhs.Op == token.LEQ {
		geqExpr, leqExpr = lhs, rhs
	} else if lhs.Op == token.LEQ && rhs.Op == token.GEQ {
		geqExpr, leqExpr = rhs, lhs
	} else {
		return
	}

	geqLHS := astutils.GoFmt(geqExpr.X)
	geqRHS := astutils.GoFmt(geqExpr.Y)
	leqLHS := astutils.GoFmt(leqExpr.X)
	leqRHS := astutils.GoFmt(leqExpr.Y)

	if geqLHS != leqLHS || geqRHS != leqRHS {
		return
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryComplexity,
		Confidence: 1,
		Node:       expr,
		Failure:    fmt.Sprintf("can simplify %s >= %s && %s <= %s to %s == %s", geqLHS, geqRHS, leqLHS, leqRHS, geqLHS, geqRHS),
	})
}

// negateComparisonOp returns the negated comparison operator or token.ILLEGAL if not applicable.
func negateComparisonOp(op token.Token) token.Token {
	switch op {
	case token.EQL:
		return token.NEQ
	case token.NEQ:
		return token.EQL
	case token.LSS:
		return token.GEQ
	case token.GEQ:
		return token.LSS
	case token.GTR:
		return token.LEQ
	case token.LEQ:
		return token.GTR
	default:
		return token.ILLEGAL
	}
}

// combineComparisonOp maps a strict comparison to its non-strict version (for combining with ==).
func combineComparisonOp(op token.Token) token.Token {
	switch op {
	case token.GTR:
		return token.GEQ
	case token.LSS:
		return token.LEQ
	default:
		return token.ILLEGAL
	}
}

// isStrictComparison returns true for > and <.
func isStrictComparison(op token.Token) bool {
	return op == token.GTR || op == token.LSS
}

// unwrapParens removes parentheses around an expression.
func unwrapParens(expr ast.Expr) ast.Expr {
	for {
		p, ok := expr.(*ast.ParenExpr)
		if !ok {
			return expr
		}
		expr = p.X
	}
}

// isIntLiteral checks if an expression is an integer literal with the given value.
func isIntLiteral(expr ast.Expr, value string) bool {
	lit, ok := unwrapParens(expr).(*ast.BasicLit)
	return ok && lit.Kind == token.INT && lit.Value == value
}

// containsFloatLiteral checks if an expression is or contains a float literal.
func containsFloatLiteral(expr ast.Expr) bool {
	expr = unwrapParens(expr)
	if lit, ok := expr.(*ast.BasicLit); ok {
		return lit.Kind == token.FLOAT
	}
	return false
}
