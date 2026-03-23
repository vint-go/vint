package use_time_method

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseTimeMethodRule detects manual conversions to milli- or microseconds
// and suggests using the time package's built-in methods instead.
type UseTimeMethodRule struct{}

// Apply applies the rule to given file.
func (r *UseTimeMethodRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintTimeMethod{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseTimeMethodRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintTimeMethod{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*UseTimeMethodRule) Name() string {
	return "useTimeMethod"
}

// Group returns the rule group.
func (*UseTimeMethodRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseTimeMethodRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintTimeMethod struct {
	onFailure func(lint.Failure)
}

func (w *lintTimeMethod) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	if binExpr.Op != token.QUO {
		return w
	}

	// Check if the left side is a call to .UnixNano()
	call, ok := binExpr.X.(*ast.CallExpr)
	if !ok {
		return w
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	if sel.Sel.Name != "UnixNano" {
		return w
	}

	// Check if there are no arguments to UnixNano()
	if len(call.Args) != 0 {
		return w
	}

	// Check the divisor
	divisor := intValue(binExpr.Y)
	if divisor == 0 {
		return w
	}

	var suggestion string
	switch divisor {
	case 1000000:
		suggestion = "UnixMilli"
	case 1000:
		suggestion = "UnixMicro"
	default:
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       binExpr,
		Failure:    fmt.Sprintf("use %s() instead of manual time conversion", suggestion),
	})

	return w
}

// intValue extracts the integer value from a basic literal expression.
// Returns 0 if the expression is not an integer literal.
func intValue(expr ast.Expr) int64 {
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		return 0
	}

	if lit.Kind != token.INT {
		return 0
	}

	var val int64
	for _, c := range lit.Value {
		if c == '_' {
			continue
		}
		if c < '0' || c > '9' {
			return 0
		}
		val = val*10 + int64(c-'0')
	}

	return val
}
