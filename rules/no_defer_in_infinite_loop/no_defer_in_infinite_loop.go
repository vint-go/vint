package no_defer_in_infinite_loop

import (
	"go/ast"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoDeferInInfiniteLoopRule detects defer statements inside infinite loops.
// Deferred function calls are executed when the surrounding function returns.
// In an infinite loop, the function never returns, so deferred calls will never
// execute and resources will leak.
type NoDeferInInfiniteLoopRule struct{}

// Apply applies the rule to given file.
func (r *NoDeferInInfiniteLoopRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintDeferInInfiniteLoop{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoDeferInInfiniteLoopRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintDeferInInfiniteLoop{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoDeferInInfiniteLoopRule) Name() string {
	return "noDeferInInfiniteLoop"
}

// Group returns the rule group.
func (*NoDeferInInfiniteLoopRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoDeferInInfiniteLoopRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintDeferInInfiniteLoop struct {
	onFailure func(lint.Failure)
}

func (w *lintDeferInInfiniteLoop) Visit(node ast.Node) ast.Visitor {
	forStmt, ok := node.(*ast.ForStmt)
	if !ok {
		return w
	}

	if !isInfiniteLoop(forStmt) {
		return w
	}

	w.findDeferInBlock(forStmt.Body)
	return w
}

// isInfiniteLoop checks whether a for statement is an infinite loop.
// A for loop is infinite when it has no condition (for { ... }) or when
// its condition is a boolean literal true (for true { ... }).
func isInfiniteLoop(forStmt *ast.ForStmt) bool {
	if forStmt.Cond == nil {
		return true
	}

	ident, ok := forStmt.Cond.(*ast.Ident)
	if ok && ident.Name == "true" {
		return true
	}

	return false
}

// findDeferInBlock walks a block statement looking for defer statements,
// but does not descend into nested function literals (closures) since those
// have their own scope and defer is fine there. It also does not descend
// into nested infinite for-loops, since those will be handled by the
// walker separately (preventing duplicate reports).
func (w *lintDeferInInfiniteLoop) findDeferInBlock(block *ast.BlockStmt) {
	if block == nil {
		return
	}
	ast.Inspect(block, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		switch v := n.(type) {
		case *ast.FuncLit:
			// Don't descend into closures; defer inside a closure
			// defined within an infinite loop is fine (it runs when the closure returns).
			return false
		case *ast.ForStmt:
			// Skip nested infinite loops — they will be handled by the main walker.
			if isInfiniteLoop(v) {
				return false
			}
		}

		deferStmt, ok := n.(*ast.DeferStmt)
		if !ok {
			return true
		}

		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       deferStmt,
			Category:   lint.FailureCategoryLogic,
			Failure:    "defers in infinite loops will never execute",
		})

		return true
	})
}
