package no_empty_critical_section

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoEmptyCriticalSectionRule detects locking a mutex and immediately unlocking
// it without doing any work in between. This is pointless and likely a mistake
// where the Unlock was meant to be deferred.
type NoEmptyCriticalSectionRule struct{}

// Apply applies the rule to given file.
func (r *NoEmptyCriticalSectionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintEmptyCriticalSection{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoEmptyCriticalSectionRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintEmptyCriticalSection{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoEmptyCriticalSectionRule) Name() string {
	return "noEmptyCriticalSection"
}

// Group returns the rule group.
func (*NoEmptyCriticalSectionRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoEmptyCriticalSectionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintEmptyCriticalSection struct {
	onFailure func(lint.Failure)
}

func (w *lintEmptyCriticalSection) Visit(node ast.Node) ast.Visitor {
	block, ok := node.(*ast.BlockStmt)
	if !ok {
		return w
	}

	// lockUnlockPairs defines the lock/unlock method pairs to check.
	lockUnlockPairs := [][2]string{
		{"Lock", "Unlock"},
		{"RLock", "RUnlock"},
	}

	for i := 0; i < len(block.List)-1; i++ {
		for _, pair := range lockUnlockPairs {
			lockReceiver, isLock := extractMutexCall(block.List[i], pair[0])
			if !isLock {
				continue
			}

			unlockReceiver, isUnlock := extractMutexCall(block.List[i+1], pair[1])
			if !isUnlock {
				continue
			}

			if astutils.GoFmt(lockReceiver) == astutils.GoFmt(unlockReceiver) {
				w.onFailure(lint.Failure{
					Category:   lint.FailureCategoryLogic,
					Confidence: 1,
					Node:       block.List[i],
					Failure:    "empty critical section, did you mean to defer the unlock?",
				})
			}
		}
	}

	return w
}

// extractMutexCall checks if the statement is an expression statement calling
// receiver.methodName() (e.g., mu.Lock() or mu.Unlock()). If it matches, it
// returns the receiver expression and true.
func extractMutexCall(stmt ast.Stmt, methodName string) (ast.Expr, bool) {
	exprStmt, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return nil, false
	}

	callExpr, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return nil, false
	}

	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}

	if selExpr.Sel.Name != methodName {
		return nil, false
	}

	return selExpr.X, true
}
