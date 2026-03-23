package no_unnecessary_blank_identifier

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoUnnecessaryBlankIdentifierRule detects unnecessary use of the blank identifier
// in range statements. When only the value from a range is needed, the blank
// identifier for the index is unnecessary and can be omitted.
type NoUnnecessaryBlankIdentifierRule struct{}

// Apply applies the rule to given file.
func (r *NoUnnecessaryBlankIdentifierRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUnnecessaryBlank{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoUnnecessaryBlankIdentifierRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUnnecessaryBlank{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoUnnecessaryBlankIdentifierRule) Name() string {
	return "noUnnecessaryBlankIdentifier"
}

// Group returns the rule group.
func (*NoUnnecessaryBlankIdentifierRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnnecessaryBlankIdentifierRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintUnnecessaryBlank struct {
	onFailure func(lint.Failure)
}

func (w *lintUnnecessaryBlank) Visit(node ast.Node) ast.Visitor {
	rangeStmt, ok := node.(*ast.RangeStmt)
	if !ok {
		return w
	}

	// Check if the key is a blank identifier
	keyIdent, ok := rangeStmt.Key.(*ast.Ident)
	if !ok || keyIdent.Name != "_" {
		return w
	}

	// If the value is nil (no value variable), then `_ = range x` is unnecessary
	// and should be `range x`
	if rangeStmt.Value == nil {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       rangeStmt,
			Category:   lint.FailureCategoryStyle,
			Failure:    "unnecessary blank identifier in range statement, use 'for range' instead",
		})
		return w
	}

	// If the value is also a blank identifier, then `_, _ = range x` is unnecessary
	// and should be `for range x`
	if valIdent, ok := rangeStmt.Value.(*ast.Ident); ok && valIdent.Name == "_" {
		// Both key and value are blank: `for _, _ = range x`
		// This is also unnecessary, should be `for range x`
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       rangeStmt,
			Category:   lint.FailureCategoryStyle,
			Failure:    "unnecessary blank identifiers in range statement, use 'for range' instead",
		})
		return w
	}

	// If key is `_` and value is a real variable with `:=`, then check if the token is DEFINE.
	// `for _, v := range x` is valid Go and the blank identifier is needed here.
	// But `for _ = range x` with just key being blank and no value is already caught above.

	// If key is `_` and value is a real variable with `=` (not `:=`), the blank is still
	// needed syntactically to have `_, v = range x`, so this is valid.

	// Only case left: key=_ and value is a real variable. This requires the blank
	// identifier syntactically since Go needs `_, v := range` to receive the value.
	// So this is valid and we don't flag it.

	return w
}

