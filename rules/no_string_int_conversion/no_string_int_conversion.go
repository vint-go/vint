package no_string_int_conversion

import (
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoStringIntConversionRule flags type conversions from integers to strings.
// In Go, string(n) where n is an integer produces the string containing the
// Unicode code point, not the decimal representation of the number.
type NoStringIntConversionRule struct{}

// Apply applies the rule to given file.
func (r *NoStringIntConversionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	file.Pkg.TypeCheck()

	w := &lintStringIntConversion{file: file, onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoStringIntConversionRule) Name() string {
	return "noStringIntConversion"
}

// Group returns the rule group.
func (*NoStringIntConversionRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoStringIntConversionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoStringIntConversionRule) RequiresTypecheck() bool {
	return true
}

type lintStringIntConversion struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintStringIntConversion) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !w.isStringConversion(ce.Fun) {
		return w
	}

	if !w.isIntegerArg(ce.Args) {
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryLogic,
		Confidence: 1,
		Node:       ce,
		Failure:    "string(int) produces a rune character, not a decimal string; use strconv.Itoa or fmt.Sprint instead",
	})

	return w
}

// isStringConversion checks if the function expression is a string type conversion.
func (w *lintStringIntConversion) isStringConversion(e ast.Expr) bool {
	t := w.file.Pkg.TypeOf(e)
	if t == nil {
		return false
	}

	tb, ok := t.Underlying().(*types.Basic)
	return ok && tb.Kind() == types.String
}

// isIntegerArg checks if the call has exactly one argument and that argument
// is an integer type (excluding byte and rune which are intentional).
func (w *lintStringIntConversion) isIntegerArg(args []ast.Expr) bool {
	if len(args) != 1 {
		return false
	}

	t := w.file.Pkg.TypeOf(args[0])
	if t == nil {
		return false
	}

	ut, ok := t.Underlying().(*types.Basic)
	if !ok || ut.Info()&types.IsInteger == 0 {
		return false
	}

	// Exclude byte and rune conversions which are intentional
	switch ut.Kind() {
	case types.Byte, types.Rune, types.UntypedRune:
		return false
	}

	return true
}
