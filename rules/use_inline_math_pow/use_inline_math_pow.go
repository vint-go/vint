package use_inline_math_pow

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseInlineMathPowRule detects calls to math.Pow with small integer exponents
// that can be replaced with inline multiplication for better performance.
type UseInlineMathPowRule struct{}

// Apply applies the rule to given file.
func (r *UseInlineMathPowRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintInlineMathPow{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseInlineMathPowRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintInlineMathPow{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*UseInlineMathPowRule) Name() string {
	return "useInlineMathPow"
}

// Group returns the rule group.
func (*UseInlineMathPowRule) Group() string {
	return "performance"
}

// CacheTier returns the cache tier for this rule.
func (*UseInlineMathPowRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintInlineMathPow struct {
	onFailure func(lint.Failure)
}

func (w *lintInlineMathPow) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(call.Fun, "math", "Pow") {
		return w
	}

	if len(call.Args) != 2 {
		return w
	}

	exponent := call.Args[1]

	// Check if the exponent is an integer literal
	lit, ok := exponent.(*ast.BasicLit)
	if !ok {
		return w
	}

	if lit.Kind != token.INT && lit.Kind != token.FLOAT {
		return w
	}

	// Parse the exponent value
	var expVal int64
	var isSmallInt bool

	if lit.Kind == token.INT {
		val, err := strconv.ParseInt(lit.Value, 0, 64)
		if err != nil {
			return w
		}
		expVal = val
		isSmallInt = true
	} else if lit.Kind == token.FLOAT {
		val, err := strconv.ParseFloat(lit.Value, 64)
		if err != nil {
			return w
		}
		// Only flag if the float is actually a whole number
		if val != float64(int64(val)) {
			return w
		}
		expVal = int64(val)
		isSmallInt = true
	}

	if !isSmallInt {
		return w
	}

	// Only flag small exponents (0 through 4)
	if expVal < 0 || expVal > 4 {
		return w
	}

	base := astutils.GoFmt(call.Args[0])
	var replacement string
	switch expVal {
	case 0:
		replacement = "1"
	case 1:
		replacement = base
	case 2:
		replacement = fmt.Sprintf("%s * %s", base, base)
	case 3:
		replacement = fmt.Sprintf("%s * %s * %s", base, base, base)
	case 4:
		replacement = fmt.Sprintf("%s * %s * %s * %s", base, base, base, base)
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryOptimization,
		Confidence: 1,
		Node:       call,
		Failure:    fmt.Sprintf("math.Pow can be replaced by inline multiplication: %s", replacement),
	})

	return w
}
