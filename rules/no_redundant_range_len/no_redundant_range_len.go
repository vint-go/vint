package no_redundant_range_len

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoRedundantRangeLenRule detects `for i := range len(slice)` that can be simplified.
type NoRedundantRangeLenRule struct{}

// Apply applies the rule to given file.
func (r *NoRedundantRangeLenRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRangeLen{file: file, onFailure: onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoRedundantRangeLenRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintRangeLen{file: file, onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoRedundantRangeLenRule) Name() string {
	return "noRedundantRangeLen"
}

// Group returns the rule group.
func (*NoRedundantRangeLenRule) Group() string {
	return "style"
}

// RequiresTypecheck returns true because this rule needs type info.
func (*NoRedundantRangeLenRule) RequiresTypecheck() bool {
	return true
}

// CacheTier returns the cache tier for this rule.
func (*NoRedundantRangeLenRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

type lintRangeLen struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintRangeLen) Visit(node ast.Node) ast.Visitor {
	rs, ok := node.(*ast.RangeStmt)
	if !ok {
		return w
	}

	// Check if the range expression is a call to len()
	call, ok := rs.X.(*ast.CallExpr)
	if !ok {
		return w
	}

	fn, ok := call.Fun.(*ast.Ident)
	if !ok || fn.Name != "len" || len(call.Args) != 1 {
		return w
	}

	// Get the argument to len()
	arg := call.Args[0]

	// Use type checking to confirm it's a slice or array
	if err := w.file.Pkg.TypeCheck(); err != nil {
		return w
	}

	typ := w.file.Pkg.TypeOf(arg)
	if typ == nil {
		return w
	}

	// Check underlying type is slice or array
	underlying := typ.Underlying()
	switch underlying.(type) {
	case *types.Slice, *types.Array:
		// OK, proceed
	default:
		return w
	}

	// Determine the argument name for the suggestion
	argStr := w.file.Render(arg)

	// Check if the loop variable is unused (blank identifier or absent)
	if isBlankOrNil(rs.Key) {
		w.onFailure(lint.Failure{
			Failure:    fmt.Sprintf("unnecessary len() in range; use `for range %s` instead", argStr),
			Confidence: 1,
			Node:       call,
			Category:   lint.FailureCategoryStyle,
		})
	} else {
		w.onFailure(lint.Failure{
			Failure:    fmt.Sprintf("unnecessary len() in range; use `for %s %s range %s` instead", w.file.Render(rs.Key), rs.Tok, argStr),
			Confidence: 1,
			Node:       call,
			Category:   lint.FailureCategoryStyle,
		})
	}

	return w
}

func isBlankOrNil(expr ast.Expr) bool {
	if expr == nil {
		return true
	}
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "_"
}
