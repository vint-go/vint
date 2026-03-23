package no_invalid_strconv_arg

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoInvalidStrconvArgRule detects invalid arguments passed to strconv.ParseInt,
// strconv.ParseUint, strconv.ParseFloat, and strconv.FormatInt.
type NoInvalidStrconvArgRule struct{}

// Apply applies the rule to given file.
func (r *NoInvalidStrconvArgRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintInvalidStrconvArg{onFailure: onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoInvalidStrconvArgRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintInvalidStrconvArg{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoInvalidStrconvArgRule) Name() string {
	return "noInvalidStrconvArg"
}

// Group returns the rule group.
func (*NoInvalidStrconvArgRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoInvalidStrconvArgRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintInvalidStrconvArg struct {
	onFailure func(lint.Failure)
}

func (w *lintInvalidStrconvArg) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if astutils.IsPkgDotName(call.Fun, "strconv", "ParseInt") {
		w.checkParseIntUint(call, "strconv.ParseInt")
		return w
	}
	if astutils.IsPkgDotName(call.Fun, "strconv", "ParseUint") {
		w.checkParseIntUint(call, "strconv.ParseUint")
		return w
	}
	if astutils.IsPkgDotName(call.Fun, "strconv", "ParseFloat") {
		w.checkParseFloat(call)
		return w
	}
	if astutils.IsPkgDotName(call.Fun, "strconv", "FormatInt") {
		w.checkFormatInt(call)
		return w
	}

	return w
}

// checkParseIntUint checks the base and bitSize arguments of strconv.ParseInt and strconv.ParseUint.
// Signature: ParseInt(s string, base int, bitSize int) (int64, error)
// base must be 0 or between 2 and 36.
// bitSize must be 0 or between 1 and 64.
func (w *lintInvalidStrconvArg) checkParseIntUint(call *ast.CallExpr, funcName string) {
	if len(call.Args) != 3 {
		return
	}

	// Check base (second argument, index 1)
	if val, ok := intLitValue(call.Args[1]); ok {
		if val != 0 && (val < 2 || val > 36) {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 1,
				Node:       call,
				Failure:    "invalid base argument to " + funcName + ": base must be 0 or between 2 and 36",
			})
		}
	}

	// Check bitSize (third argument, index 2)
	if val, ok := intLitValue(call.Args[2]); ok {
		if val < 0 || val > 64 {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 1,
				Node:       call,
				Failure:    "invalid bitSize argument to " + funcName + ": bitSize must be between 0 and 64",
			})
		}
	}
}

// checkParseFloat checks the bitSize argument of strconv.ParseFloat.
// Signature: ParseFloat(s string, bitSize int) (float64, error)
// bitSize must be 32 or 64.
func (w *lintInvalidStrconvArg) checkParseFloat(call *ast.CallExpr) {
	if len(call.Args) != 2 {
		return
	}

	// Check bitSize (second argument, index 1)
	if val, ok := intLitValue(call.Args[1]); ok {
		if val != 32 && val != 64 {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 1,
				Node:       call,
				Failure:    "invalid bitSize argument to strconv.ParseFloat: bitSize must be 32 or 64",
			})
		}
	}
}

// checkFormatInt checks the base argument of strconv.FormatInt.
// Signature: FormatInt(i int64, base int) string
// base must be between 2 and 36.
func (w *lintInvalidStrconvArg) checkFormatInt(call *ast.CallExpr) {
	if len(call.Args) != 2 {
		return
	}

	// Check base (second argument, index 1)
	if val, ok := intLitValue(call.Args[1]); ok {
		if val < 2 || val > 36 {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 1,
				Node:       call,
				Failure:    "invalid base argument to strconv.FormatInt: base must be between 2 and 36",
			})
		}
	}
}

// intLitValue extracts the integer value from an *ast.BasicLit node.
// It handles positive integer literals and negative unary expressions (e.g. -1).
// Returns the value and true if successful, or 0 and false otherwise.
func intLitValue(expr ast.Expr) (int64, bool) {
	// Handle negative unary expressions like -1
	if unary, ok := expr.(*ast.UnaryExpr); ok && unary.Op == token.SUB {
		if lit, ok := unary.X.(*ast.BasicLit); ok && lit.Kind == token.INT {
			v, err := strconv.ParseInt(lit.Value, 0, 64)
			if err != nil {
				return 0, false
			}
			return -v, true
		}
		return 0, false
	}

	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.INT {
		return 0, false
	}
	v, err := strconv.ParseInt(lit.Value, 0, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}
