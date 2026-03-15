package no_bad_lock_pattern

import (
	"fmt"
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoBadLockPatternRule detects suspicious mutex lock/unlock operations such as
// double locking without an unlock, double unlocking, or locking immediately
// followed by unlocking without any work in the critical section.
type NoBadLockPatternRule struct{}

// Apply applies the rule to given file.
func (r *NoBadLockPatternRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintBadLockPattern{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoBadLockPatternRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintBadLockPattern{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoBadLockPatternRule) Name() string {
	return "noBadLockPattern"
}

// Group returns the rule group.
func (*NoBadLockPatternRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoBadLockPatternRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// lockMethods is the set of lock-related methods we track.
var lockMethods = map[string]bool{
	"Lock":    true,
	"Unlock":  true,
	"RLock":   true,
	"RUnlock": true,
}

// lockPair maps a lock method to its corresponding unlock method.
var lockPair = map[string]string{
	"Lock":  "Unlock",
	"RLock": "RUnlock",
}

// unlockPair maps an unlock method to its corresponding lock method.
var unlockPair = map[string]string{
	"Unlock":  "Lock",
	"RUnlock": "RLock",
}

type lintBadLockPattern struct {
	onFailure func(lint.Failure)
}

func (w *lintBadLockPattern) Visit(node ast.Node) ast.Visitor {
	block, ok := node.(*ast.BlockStmt)
	if !ok {
		return w
	}

	w.checkStatementList(block.List)
	return w
}

// lockCallInfo stores information about a lock/unlock call.
type lockCallInfo struct {
	receiver   string // string representation of the receiver expression
	method     string // "Lock", "Unlock", "RLock", "RUnlock"
	node       ast.Node
	isDeferred bool // whether this call is inside a defer statement
}

// checkStatementList examines a list of statements for bad lock patterns.
func (w *lintBadLockPattern) checkStatementList(stmts []ast.Stmt) {
	// Track the last lock operation per receiver.
	// Key: receiver string, Value: last lock call info.
	lastOp := map[string]*lockCallInfo{}

	for _, stmt := range stmts {
		info := extractLockCall(stmt)
		if info == nil {
			// Any non-lock/unlock statement clears the tracking for all receivers,
			// since we only care about consecutive lock/unlock calls.
			lastOp = map[string]*lockCallInfo{}
			continue
		}

		prev, hasPrev := lastOp[info.receiver]
		if hasPrev {
			// Check for bad patterns
			if prev.method == info.method {
				// Double lock or double unlock
				w.onFailure(lint.Failure{
					Category:   lint.FailureCategoryLogic,
					Confidence: 1,
					Node:       info.node,
					Failure:    fmt.Sprintf("suspicious lock sequence: %s called twice without %s", info.method, counterpart(info.method)),
				})
			} else if isLockMethod(prev.method) && isUnlockMethod(info.method) && isMatchingPair(prev.method, info.method) && !info.isDeferred {
				// Lock immediately followed by Unlock (empty critical section)
				// Skip if the unlock is deferred, as mu.Lock(); defer mu.Unlock() is a standard pattern.
				w.onFailure(lint.Failure{
					Category:   lint.FailureCategoryLogic,
					Confidence: 1,
					Node:       info.node,
					Failure:    fmt.Sprintf("suspicious lock sequence: %s immediately followed by %s with no work in critical section", prev.method, info.method),
				})
			}
		}

		lastOp[info.receiver] = info
	}
}

// extractLockCall extracts lock/unlock call information from a statement.
// Returns nil if the statement is not a lock/unlock call.
func extractLockCall(stmt ast.Stmt) *lockCallInfo {
	// Check for expression statements (most common: mu.Lock())
	exprStmt, ok := stmt.(*ast.ExprStmt)
	if !ok {
		// Also check for defer statements: defer mu.Unlock()
		deferStmt, ok := stmt.(*ast.DeferStmt)
		if !ok {
			return nil
		}
		info := extractLockCallFromCallExpr(deferStmt.Call)
		if info != nil {
			info.isDeferred = true
		}
		return info
	}

	call, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return nil
	}

	return extractLockCallFromCallExpr(call)
}

// extractLockCallFromCallExpr extracts lock/unlock info from a call expression.
func extractLockCallFromCallExpr(call *ast.CallExpr) *lockCallInfo {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	methodName := sel.Sel.Name
	if !lockMethods[methodName] {
		return nil
	}

	receiver := astutils.GoFmt(sel.X)
	return &lockCallInfo{
		receiver: receiver,
		method:   methodName,
		node:     call,
	}
}

// counterpart returns the opposite lock/unlock operation.
func counterpart(method string) string {
	if pair, ok := lockPair[method]; ok {
		return pair
	}
	if pair, ok := unlockPair[method]; ok {
		return pair
	}
	return ""
}

// isLockMethod returns true if the method is a locking method.
func isLockMethod(method string) bool {
	_, ok := lockPair[method]
	return ok
}

// isUnlockMethod returns true if the method is an unlocking method.
func isUnlockMethod(method string) bool {
	_, ok := unlockPair[method]
	return ok
}

// isMatchingPair returns true if the lock and unlock methods form a matching pair.
func isMatchingPair(lock, unlock string) bool {
	expected, ok := lockPair[lock]
	return ok && expected == unlock
}
