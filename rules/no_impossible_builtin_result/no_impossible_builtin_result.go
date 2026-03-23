package no_impossible_builtin_result

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoImpossibleBuiltinResultRule detects comparisons of builtin function
// results against values they can never produce. For example, len() and cap()
// always return non-negative values, so len(x) < 0 is always false.
type NoImpossibleBuiltinResultRule struct{}

// Apply applies the rule to given file.
func (r *NoImpossibleBuiltinResultRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintImpossibleBuiltin{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoImpossibleBuiltinResultRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintImpossibleBuiltin{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoImpossibleBuiltinResultRule) Name() string {
	return "noImpossibleBuiltinResult"
}

// Group returns the rule group.
func (*NoImpossibleBuiltinResultRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoImpossibleBuiltinResultRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintImpossibleBuiltin struct {
	onFailure func(lint.Failure)
}

// nonNegativeBuiltins is the set of builtin functions that always return >= 0.
var nonNegativeBuiltins = map[string]bool{
	"len": true,
	"cap": true,
}

func (w *lintImpossibleBuiltin) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	// Check both orientations: builtin(x) op literal and literal op builtin(x)
	if builtinName, ok := getNonNegativeBuiltinCall(binExpr.X); ok {
		if val, ok := negativeIntLiteral(binExpr.Y); ok {
			w.checkBuiltinVsNegative(binExpr, builtinName, binExpr.Op, val)
		}
	} else if builtinName, ok := getNonNegativeBuiltinCall(binExpr.Y); ok {
		if val, ok := negativeIntLiteral(binExpr.X); ok {
			// Flip the operator since the builtin is on the right side
			w.checkBuiltinVsNegative(binExpr, builtinName, flipOp(binExpr.Op), val)
		}
	}

	return w
}

// checkBuiltinVsNegative checks if comparing a non-negative builtin result
// against a negative value produces an impossible or always-true condition.
// The op is always from the builtin's perspective: builtin op value.
func (w *lintImpossibleBuiltin) checkBuiltinVsNegative(node ast.Node, builtinName string, op token.Token, val int64) {
	// builtin() always returns >= 0, so comparisons against negative values
	// have predetermined results:
	//   builtin(x) < negative  -> always false
	//   builtin(x) <= negative -> always false
	//   builtin(x) == negative -> always false
	//   builtin(x) > negative  -> always true
	//   builtin(x) >= negative -> always true
	//   builtin(x) != negative -> always true
	switch op {
	case token.LSS, token.LEQ, token.EQL:
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       node,
			Failure:    fmt.Sprintf("%s() never returns a negative value, comparison always evaluates to false", builtinName),
		})
	case token.GTR, token.GEQ, token.NEQ:
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       node,
			Failure:    fmt.Sprintf("%s() never returns a negative value, comparison always evaluates to true", builtinName),
		})
	}
}

// getNonNegativeBuiltinCall checks if the expression is a call to a builtin
// that always returns a non-negative value (len, cap). Returns the builtin
// name if so.
func getNonNegativeBuiltinCall(expr ast.Expr) (string, bool) {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return "", false
	}
	ident, ok := call.Fun.(*ast.Ident)
	if !ok {
		return "", false
	}
	if nonNegativeBuiltins[ident.Name] && len(call.Args) == 1 {
		return ident.Name, true
	}
	return "", false
}

// negativeIntLiteral checks if the expression is a unary negation of an integer
// literal (e.g., -1, -5). Returns the negative value if so.
func negativeIntLiteral(expr ast.Expr) (int64, bool) {
	unary, ok := expr.(*ast.UnaryExpr)
	if !ok {
		return 0, false
	}
	if unary.Op != token.SUB {
		return 0, false
	}
	lit, ok := unary.X.(*ast.BasicLit)
	if !ok {
		return 0, false
	}
	if lit.Kind != token.INT {
		return 0, false
	}
	v, err := strconv.ParseInt(lit.Value, 0, 64)
	if err != nil {
		return 0, false
	}
	return -v, true
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
