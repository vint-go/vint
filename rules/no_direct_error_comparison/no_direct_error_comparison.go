package no_direct_error_comparison

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/lint"
)

// NoDirectErrorComparisonRule flags direct comparisons of error values using == or != and recommends errors.Is().
type NoDirectErrorComparisonRule struct{}

// Apply applies the rule to given file.
func (*NoDirectErrorComparisonRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintNoDirectErrorComparison{file: file, onFailure: onFailure}
	if w.file.Pkg.TypeCheck() != nil {
		return nil
	}

	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (*NoDirectErrorComparisonRule) Name() string {
	return "noDirectErrorComparison"
}

// Group returns the rule group.
func (*NoDirectErrorComparisonRule) Group() string {
	return "correctness"
}

func (*NoDirectErrorComparisonRule) RequiresTypecheck() bool {
	return true
}

type lintNoDirectErrorComparison struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintNoDirectErrorComparison) Visit(node ast.Node) ast.Visitor {
	expr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	switch expr.Op {
	case token.EQL, token.NEQ:
	default:
		return w
	}

	// Check if either side is nil — always allowed
	if isNilIdent(expr.X) || isNilIdent(expr.Y) {
		return w
	}

	// Check if either side is io.EOF — always allowed
	if astutils.IsPkgDotName(expr.X, "io", "EOF") || astutils.IsPkgDotName(expr.Y, "io", "EOF") {
		return w
	}

	typeOfX := w.file.Pkg.TypeOf(expr.X)
	typeOfY := w.file.Pkg.TypeOf(expr.Y)

	if typeOfX == nil || typeOfY == nil {
		return w
	}

	// Check if at least one side implements the error interface
	errorIface := types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
	xIsError := types.Implements(typeOfX, errorIface)
	yIsError := types.Implements(typeOfY, errorIface)

	if !xIsError && !yIsError {
		return w
	}

	xStr := astutils.GoFmt(expr.X)
	yStr := astutils.GoFmt(expr.Y)

	var failureMsg string
	var replacementLine string
	if expr.Op == token.EQL {
		failureMsg = fmt.Sprintf("avoid direct error comparison, use errors.Is(%s, %s) instead", xStr, yStr)
		replacementLine = fmt.Sprintf("errors.Is(%s, %s)", xStr, yStr)
	} else {
		failureMsg = fmt.Sprintf("avoid direct error comparison, use !errors.Is(%s, %s) instead", xStr, yStr)
		replacementLine = fmt.Sprintf("!errors.Is(%s, %s)", xStr, yStr)
	}

	w.onFailure(lint.Failure{
		Category:        lint.FailureCategoryErrors,
		Confidence:      1,
		Node:            node,
		Failure:         failureMsg,
		ReplacementLine: replacementLine,
	})

	return w
}

func isNilIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "nil"
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name lowercased and without hyphens.
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}
