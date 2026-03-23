package no_defer_in_loop

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoDeferInLoopRule detects defer statements inside loops.
type NoDeferInLoopRule struct{}

// Apply applies the rule to given file.
func (r *NoDeferInLoopRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintDeferInLoop{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoDeferInLoopRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintDeferInLoop{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoDeferInLoopRule) Name() string {
	return "noDeferInLoop"
}

// Group returns the rule group.
func (*NoDeferInLoopRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoDeferInLoopRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintDeferInLoop struct {
	onFailure func(lint.Failure)
}

func (w *lintDeferInLoop) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.ForStmt:
		w.findDeferInBlock(n.Body)
		return nil
	case *ast.RangeStmt:
		w.findDeferInBlock(n.Body)
		return nil
	}
	return w
}

// findDeferInBlock walks a block statement looking for defer statements,
// but does not descend into nested function literals (closures) since those
// have their own scope and defer is fine there.
func (w *lintDeferInLoop) findDeferInBlock(block *ast.BlockStmt) {
	if block == nil {
		return
	}
	ast.Inspect(block, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.FuncLit:
			// Don't descend into closures; defer inside a closure
			// defined within a loop is fine (it runs when the closure returns).
			return false
		}

		deferStmt, ok := n.(*ast.DeferStmt)
		if !ok {
			return true
		}

		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       deferStmt,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "defer statement inside a loop; deferred call executes only when the function returns, not on each iteration",
		})

		return true
	})
}
