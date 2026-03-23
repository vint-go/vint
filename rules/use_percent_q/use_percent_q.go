package use_percent_q

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UsePercentQRule detects fmt.Sprintf calls where "%s" is wrapped in escaped
// quotes and can be replaced with the %q verb.
type UsePercentQRule struct{}

// Apply applies the rule to given file.
func (r *UsePercentQRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	w := &lintPercentQ{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UsePercentQRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	w := &lintPercentQ{onFailure: onFailure}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*UsePercentQRule) Name() string {
	return "usePercentQ"
}

// Group returns the rule group.
func (*UsePercentQRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UsePercentQRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintPercentQ struct {
	onFailure func(lint.Failure)
}

// formatFunctions lists fmt functions that accept a format string as the first argument.
var formatFunctions = []string{
	"Sprintf",
	"Fprintf",
	"Errorf",
	"Printf",
}

func (w *lintPercentQ) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check if this is a call to a fmt format function
	if !w.isFmtFormatCall(call) {
		return w
	}

	// The first argument (or second for Fprintf) is the format string.
	formatArg := w.getFormatArg(call)
	if formatArg == nil {
		return w
	}

	// Must be a string literal
	lit, ok := formatArg.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return w
	}

	// Check if the format string contains \"%s\" pattern
	// In the raw BasicLit.Value (which includes surrounding quotes),
	// the pattern looks like \"%s\" for interpreted strings
	if containsQuotedPercentS(lit.Value) {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryStyle,
			Confidence: 1,
			Node:       call,
			Failure:    `use %q instead of \"%s\" for quoted string formatting`,
		})
	}

	return w
}

// isFmtFormatCall checks if the call is to one of the fmt format functions.
func (w *lintPercentQ) isFmtFormatCall(call *ast.CallExpr) bool {
	for _, fn := range formatFunctions {
		if astutils.IsPkgDotName(call.Fun, "fmt", fn) {
			return true
		}
	}
	return false
}

// getFormatArg returns the format string argument from the call.
// For Fprintf, the format string is the second argument; for others it's the first.
func (w *lintPercentQ) getFormatArg(call *ast.CallExpr) ast.Expr {
	if len(call.Args) == 0 {
		return nil
	}

	// Fprintf takes (writer, format, args...) so format is at index 1
	if astutils.IsPkgDotName(call.Fun, "fmt", "Fprintf") {
		if len(call.Args) < 2 {
			return nil
		}
		return call.Args[1]
	}

	// All others: format is the first argument
	return call.Args[0]
}

// containsQuotedPercentS checks if the raw string literal value contains
// the pattern \"%s\" (an escaped-quote-wrapped %s verb).
func containsQuotedPercentS(rawValue string) bool {
	// rawValue includes the surrounding quotes from the source code.
	// For an interpreted string like "value is \"%s\"", rawValue is:
	//   "value is \"%s\""
	// We look for the sequence \"%s\"
	return strings.Contains(rawValue, `\"%s\"`)
}
