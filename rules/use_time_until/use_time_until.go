package use_time_until

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseTimeUntilRule detects x.Sub(time.Now()) and suggests time.Until(x).
type UseTimeUntilRule struct{}

// Apply applies the rule to given file.
func (r *UseTimeUntilRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintTimeUntil{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseTimeUntilRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintTimeUntil{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*UseTimeUntilRule) Name() string {
	return "useTimeUntil"
}

// Group returns the rule group.
func (*UseTimeUntilRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseTimeUntilRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintTimeUntil struct {
	onFailure func(lint.Failure)
}

func (w *lintTimeUntil) Visit(node ast.Node) ast.Visitor {
	// Look for a call expression: something.Sub(arg)
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// The function must be a selector expression: receiver.Sub
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	if sel.Sel.Name != "Sub" {
		return w
	}

	// Sub takes exactly one argument
	if len(call.Args) != 1 {
		return w
	}

	// The argument must be a call expression: time.Now()
	argCall, ok := call.Args[0].(*ast.CallExpr)
	if !ok {
		return w
	}

	// time.Now() has no arguments
	if len(argCall.Args) != 0 {
		return w
	}

	// The function of the argument call must be time.Now
	argSel, ok := argCall.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	ident, ok := argSel.X.(*ast.Ident)
	if !ok {
		return w
	}

	if ident.Name != "time" || argSel.Sel.Name != "Now" {
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       call,
		Failure:    "replace x.Sub(time.Now()) with time.Until(x)",
	})

	return w
}
