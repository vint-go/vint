package no_impossible_condition

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoImpossibleConditionRule detects suspicious condition expressions that are
// likely logic errors, such as inverted loop conditions, contradictory range
// checks, and simultaneous equality comparisons to different values.
type NoImpossibleConditionRule struct{}

// Apply applies the rule to given file.
func (r *NoImpossibleConditionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if file.Pkg.TypeCheck() != nil {
		return nil
	}

	w := &lintImpossibleCondition{
		file: file,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoImpossibleConditionRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintImpossibleCondition{file: file, onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoImpossibleConditionRule) Name() string {
	return "noImpossibleCondition"
}

// Group returns the rule group.
func (*NoImpossibleConditionRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoImpossibleConditionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule uses type info for constant evaluation.
func (*NoImpossibleConditionRule) RequiresTypecheck() bool {
	return true
}

type lintImpossibleCondition struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintImpossibleCondition) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.ForStmt:
		w.checkForStmt(n)
	case *ast.BinaryExpr:
		if n.Op == token.LAND {
			w.checkImpossibleAnd(n)
		}
	}
	return w
}

// checkForStmt detects inverted loop conditions like `for i := 0; i > n; i++`.
func (w *lintImpossibleCondition) checkForStmt(forStmt *ast.ForStmt) {
	if forStmt.Cond == nil || forStmt.Init == nil || forStmt.Post == nil {
		return
	}

	// Check for init of the form `i := 0` or `i = 0`.
	initVar, initVal := w.extractAssignment(forStmt.Init)
	if initVar == "" {
		return
	}

	// Check post is increment (i++ or i += ...)
	if !w.isIncrement(forStmt.Post, initVar) {
		return
	}

	// Check condition is a comparison involving initVar.
	cond, ok := forStmt.Cond.(*ast.BinaryExpr)
	if !ok {
		return
	}

	// For `i := 0; i > n; i++`, the loop variable is on one side.
	condVar, condOp, otherExpr := w.extractComparison(cond, initVar)
	if condVar == "" {
		return
	}

	// If we know the init value, check if the condition is impossible.
	// For ascending loops (i++), the condition should be i < n or i <= n or i != n.
	// If the condition is i > n or i >= n, this is likely inverted.
	switch condOp {
	case token.GTR: // i > n
		// Only flag if we know init val is small or n is likely positive.
		if w.isLikelyInvertedAscending(initVal, otherExpr, condOp) {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 0.8,
				Node:       forStmt.Cond,
				Failure:    "suspicious loop condition: loop variable increments but condition uses >",
			})
		}
	case token.GEQ: // i >= n
		if w.isLikelyInvertedAscending(initVal, otherExpr, condOp) {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 0.8,
				Node:       forStmt.Cond,
				Failure:    "suspicious loop condition: loop variable increments but condition uses >=",
			})
		}
	}
}

// extractAssignment extracts variable name and initial value from an init statement.
func (w *lintImpossibleCondition) extractAssignment(stmt ast.Stmt) (string, *int64) {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		if len(s.Lhs) != 1 || len(s.Rhs) != 1 {
			return "", nil
		}
		ident, ok := s.Lhs[0].(*ast.Ident)
		if !ok {
			return "", nil
		}
		val := w.constIntValue(s.Rhs[0])
		return ident.Name, val
	}
	return "", nil
}

// isIncrement checks if the post statement increments the given variable.
func (w *lintImpossibleCondition) isIncrement(stmt ast.Stmt, varName string) bool {
	switch s := stmt.(type) {
	case *ast.IncDecStmt:
		if s.Tok == token.INC {
			if ident, ok := s.X.(*ast.Ident); ok && ident.Name == varName {
				return true
			}
		}
	case *ast.AssignStmt:
		if s.Tok == token.ADD_ASSIGN && len(s.Lhs) == 1 {
			if ident, ok := s.Lhs[0].(*ast.Ident); ok && ident.Name == varName {
				return true
			}
		}
	}
	return false
}

// extractComparison extracts comparison info given a variable name.
// Returns the matched variable name, the operator as seen from the variable's side, and the other expression.
func (w *lintImpossibleCondition) extractComparison(expr *ast.BinaryExpr, varName string) (string, token.Token, ast.Expr) {
	// Check if left side is the variable.
	if ident, ok := expr.X.(*ast.Ident); ok && ident.Name == varName {
		return varName, expr.Op, expr.Y
	}
	// Check if right side is the variable (flip the operator).
	if ident, ok := expr.Y.(*ast.Ident); ok && ident.Name == varName {
		return varName, flipOp(expr.Op), expr.X
	}
	return "", 0, nil
}

// isLikelyInvertedAscending determines whether the loop condition is likely inverted for an ascending loop.
func (w *lintImpossibleCondition) isLikelyInvertedAscending(initVal *int64, _ ast.Expr, _ token.Token) bool {
	// If init value is 0 or a small non-negative number, i > n is almost certainly wrong.
	if initVal != nil && *initVal >= 0 {
		return true
	}
	// If we don't know the init value, still flag it as suspicious
	// since ascending loops with > conditions are almost always wrong.
	return initVal != nil
}

// checkImpossibleAnd checks for impossible conditions in && expressions.
func (w *lintImpossibleCondition) checkImpossibleAnd(expr *ast.BinaryExpr) {
	left, okL := expr.X.(*ast.BinaryExpr)
	right, okR := expr.Y.(*ast.BinaryExpr)
	if !okL || !okR {
		return
	}

	// Check for contradictory comparisons: x == a && x == b (where a != b).
	if w.checkDoubleEquality(expr, left, right) {
		return
	}

	// Check for contradictory range: x < a && x > b (where a <= b).
	w.checkContradictoryRange(expr, left, right)
}

// checkDoubleEquality detects patterns like x == 1 && x == 2.
func (w *lintImpossibleCondition) checkDoubleEquality(parent *ast.BinaryExpr, left, right *ast.BinaryExpr) bool {
	if left.Op != token.EQL || right.Op != token.EQL {
		return false
	}

	// Extract the common variable and the two values.
	leftVar, leftVal := w.extractEqualityParts(left)
	rightVar, rightVal := w.extractEqualityParts(right)

	if leftVar == "" || rightVar == "" || leftVar != rightVar {
		return false
	}

	// If both are constants and different, this is impossible.
	if leftVal != nil && rightVal != nil && *leftVal != *rightVal {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       parent,
			Failure:    fmt.Sprintf("impossible condition: %s cannot equal both values simultaneously", leftVar),
		})
		return true
	}

	// If the compared-to expressions are different identifiers, still flag.
	leftOther := w.otherSide(left, leftVar)
	rightOther := w.otherSide(right, rightVar)
	if leftOther != nil && rightOther != nil {
		leftIdent, lok := leftOther.(*ast.Ident)
		rightIdent, rok := rightOther.(*ast.Ident)
		if lok && rok && leftIdent.Name != rightIdent.Name {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 0.8,
				Node:       parent,
				Failure:    fmt.Sprintf("suspicious condition: %s compared for equality with two different values", leftVar),
			})
			return true
		}
	}

	return false
}

// extractEqualityParts extracts the variable name and optional constant value from an equality check.
func (w *lintImpossibleCondition) extractEqualityParts(expr *ast.BinaryExpr) (string, *int64) {
	if ident, ok := expr.X.(*ast.Ident); ok {
		val := w.constIntValue(expr.Y)
		return ident.Name, val
	}
	if ident, ok := expr.Y.(*ast.Ident); ok {
		val := w.constIntValue(expr.X)
		return ident.Name, val
	}
	return "", nil
}

// otherSide returns the expression on the other side of an equality from the named variable.
func (w *lintImpossibleCondition) otherSide(expr *ast.BinaryExpr, varName string) ast.Expr {
	if ident, ok := expr.X.(*ast.Ident); ok && ident.Name == varName {
		return expr.Y
	}
	if ident, ok := expr.Y.(*ast.Ident); ok && ident.Name == varName {
		return expr.X
	}
	return nil
}

// checkContradictoryRange detects patterns like x < 5 && x > 10.
func (w *lintImpossibleCondition) checkContradictoryRange(parent *ast.BinaryExpr, left, right *ast.BinaryExpr) {
	// Normalize: find a common variable and check if the ranges are contradictory.
	leftVar, leftOp, leftBound := w.extractRangePart(left)
	rightVar, rightOp, rightBound := w.extractRangePart(right)

	if leftVar == "" || rightVar == "" || leftVar != rightVar {
		return
	}

	if leftBound == nil || rightBound == nil {
		return
	}

	// Check if the range is impossible.
	// For x < a && x > b: impossible if a <= b
	// For x <= a && x > b: impossible if a <= b
	// For x < a && x >= b: impossible if a <= b
	// For x <= a && x >= b: impossible if a < b
	if w.isContradictory(leftOp, *leftBound, rightOp, *rightBound) {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       parent,
			Failure:    fmt.Sprintf("impossible condition: %s cannot satisfy both comparisons simultaneously", leftVar),
		})
	}
}

// extractRangePart extracts variable, normalized operator (from variable's perspective), and constant bound.
func (w *lintImpossibleCondition) extractRangePart(expr *ast.BinaryExpr) (string, token.Token, *int64) {
	// Check left side is ident.
	if ident, ok := expr.X.(*ast.Ident); ok {
		val := w.constIntValue(expr.Y)
		return ident.Name, expr.Op, val
	}
	// Check right side is ident, flip operator.
	if ident, ok := expr.Y.(*ast.Ident); ok {
		val := w.constIntValue(expr.X)
		return ident.Name, flipOp(expr.Op), val
	}
	return "", 0, nil
}

// isContradictory checks if two range conditions on the same variable are contradictory.
// leftOp/leftBound: first condition (e.g., x < 5 means op=LSS, bound=5)
// rightOp/rightBound: second condition (e.g., x > 10 means op=GTR, bound=10)
func (w *lintImpossibleCondition) isContradictory(leftOp token.Token, leftBound int64, rightOp token.Token, rightBound int64) bool {
	// Determine the upper and lower bounds.
	// We need to check: is there any value that satisfies both conditions?

	// Normalize to: var has upper bound (from < or <=) and lower bound (from > or >=).
	var hasUpper, hasLower bool
	var upper, lower int64
	var upperStrict, lowerStrict bool // strict means < or >, not <= or >=

	switch leftOp {
	case token.LSS: // x < bound
		hasUpper = true
		upper = leftBound
		upperStrict = true
	case token.LEQ: // x <= bound
		hasUpper = true
		upper = leftBound
		upperStrict = false
	case token.GTR: // x > bound
		hasLower = true
		lower = leftBound
		lowerStrict = true
	case token.GEQ: // x >= bound
		hasLower = true
		lower = leftBound
		lowerStrict = false
	}

	switch rightOp {
	case token.LSS:
		hasUpper = true
		upper = rightBound
		upperStrict = true
	case token.LEQ:
		hasUpper = true
		upper = rightBound
		upperStrict = false
	case token.GTR:
		hasLower = true
		lower = rightBound
		lowerStrict = true
	case token.GEQ:
		hasLower = true
		lower = rightBound
		lowerStrict = false
	}

	if !hasUpper || !hasLower {
		return false
	}

	// Check if range is empty.
	if upperStrict && lowerStrict {
		// x < upper && x > lower: impossible if upper <= lower
		return upper <= lower
	} else if upperStrict {
		// x < upper && x >= lower: impossible if upper <= lower
		return upper <= lower
	} else if lowerStrict {
		// x <= upper && x > lower: impossible if upper <= lower
		return upper <= lower
	} else {
		// x <= upper && x >= lower: impossible if upper < lower
		return upper < lower
	}
}

// constIntValue tries to extract a constant integer value from an expression using type info.
func (w *lintImpossibleCondition) constIntValue(expr ast.Expr) *int64 {
	info := w.file.Pkg.TypesInfo()
	if info == nil {
		return w.constIntValueFromLiteral(expr)
	}

	tv, ok := info.Types[expr]
	if !ok || tv.Value == nil {
		return w.constIntValueFromLiteral(expr)
	}

	if tv.Value.Kind() != constant.Int {
		return nil
	}

	val, exact := constant.Int64Val(tv.Value)
	if !exact {
		return nil
	}

	return &val
}

// constIntValueFromLiteral is a fallback that extracts constant value from a BasicLit.
func (w *lintImpossibleCondition) constIntValueFromLiteral(expr ast.Expr) *int64 {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.INT {
		return nil
	}

	val := constant.MakeFromLiteral(lit.Value, token.INT, 0)
	if val.Kind() != constant.Int {
		return nil
	}

	v, exact := constant.Int64Val(val)
	if !exact {
		return nil
	}

	return &v
}

// flipOp returns the mirrored comparison operator.
func flipOp(op token.Token) token.Token {
	switch op {
	case token.LSS:
		return token.GTR
	case token.GTR:
		return token.LSS
	case token.LEQ:
		return token.GEQ
	case token.GEQ:
		return token.LEQ
	case token.EQL:
		return token.EQL
	case token.NEQ:
		return token.NEQ
	default:
		return op
	}
}
