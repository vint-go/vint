package use_time_since

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseTimeSinceRule detects time.Now().Sub(x) and suggests time.Since(x).
type UseTimeSinceRule struct{}

// Apply applies the rule to given file.
func (r *UseTimeSinceRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintTimeSince{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseTimeSinceRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintTimeSince{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*UseTimeSinceRule) Name() string {
	return "useTimeSince"
}

// Group returns the rule group.
func (*UseTimeSinceRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseTimeSinceRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintTimeSince struct {
	onFailure func(lint.Failure)
}

func (w *lintTimeSince) Visit(node ast.Node) ast.Visitor {
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

	// The receiver must be a call expression: time.Now()
	recvCall, ok := sel.X.(*ast.CallExpr)
	if !ok {
		return w
	}

	// time.Now() has no arguments
	if len(recvCall.Args) != 0 {
		return w
	}

	// The function of the receiver call must be time.Now
	recvSel, ok := recvCall.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	ident, ok := recvSel.X.(*ast.Ident)
	if !ok {
		return w
	}

	if ident.Name != "time" || recvSel.Sel.Name != "Now" {
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       call,
		Failure:    "replace time.Now().Sub(x) with time.Since(x)",
	})

	return w
}
