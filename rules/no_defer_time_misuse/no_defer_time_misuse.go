package no_defer_time_misuse

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoDeferTimeMisuseRule detects deferred calls to time.Since which evaluates its
// argument immediately rather than at defer time.
type NoDeferTimeMisuseRule struct{}

// Apply applies the rule to given file.
func (r *NoDeferTimeMisuseRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintDeferTimeMisuse{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoDeferTimeMisuseRule) Name() string {
	return "noDeferTimeMisuse"
}

// Group returns the rule group.
func (*NoDeferTimeMisuseRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoDeferTimeMisuseRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintDeferTimeMisuse struct {
	onFailure func(lint.Failure)
}

func (w *lintDeferTimeMisuse) Visit(node ast.Node) ast.Visitor {
	deferStmt, ok := node.(*ast.DeferStmt)
	if !ok {
		return w
	}

	// Check the deferred call expression for time.Since calls that are NOT
	// inside a func literal (closure). If time.Since is called directly as
	// an argument to the deferred function, it is evaluated immediately.
	w.checkCallForTimeSince(deferStmt.Call)

	return w
}

// checkCallForTimeSince checks the arguments of a call expression for direct
// calls to time.Since. It does NOT recurse into func literals (closures),
// because wrapping in a closure is the correct fix.
func (w *lintDeferTimeMisuse) checkCallForTimeSince(call *ast.CallExpr) {
	// Check if the deferred call itself is time.Since
	if astutils.IsPkgDotName(call.Fun, "time", "Since") {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       call,
			Category:   lint.FailureCategoryLogic,
			Failure:    "time.Since is evaluated immediately in defer, not when the deferred function runs",
		})
		return
	}

	// Check the arguments of the deferred call for time.Since
	for _, arg := range call.Args {
		w.findTimeSinceInExpr(arg)
	}
}

// findTimeSinceInExpr recursively searches for time.Since calls within an expression,
// but does NOT descend into func literals (closures).
func (w *lintDeferTimeMisuse) findTimeSinceInExpr(expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.CallExpr:
		if astutils.IsPkgDotName(e.Fun, "time", "Since") {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       e,
				Category:   lint.FailureCategoryLogic,
				Failure:    "time.Since is evaluated immediately in defer, not when the deferred function runs",
			})
			return
		}
		// Check arguments of nested calls
		for _, arg := range e.Args {
			w.findTimeSinceInExpr(arg)
		}
	case *ast.FuncLit:
		// Do not descend into closures — wrapping in a closure is the correct pattern
		return
	case *ast.ParenExpr:
		w.findTimeSinceInExpr(e.X)
	case *ast.UnaryExpr:
		w.findTimeSinceInExpr(e.X)
	case *ast.BinaryExpr:
		w.findTimeSinceInExpr(e.X)
		w.findTimeSinceInExpr(e.Y)
	}
}
