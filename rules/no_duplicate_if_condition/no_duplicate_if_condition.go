package no_duplicate_if_condition

import (
	"fmt"
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoDuplicateIfConditionRule detects duplicate conditions in if/else if chains.
type NoDuplicateIfConditionRule struct{}

// Apply applies the rule to given file.
func (r *NoDuplicateIfConditionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoDupIfCond{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
		visited: map[*ast.IfStmt]bool{},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoDuplicateIfConditionRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoDupIfCond{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
		visited: map[*ast.IfStmt]bool{},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoDuplicateIfConditionRule) Name() string {
	return "noDuplicateIfCondition"
}

// Group returns the rule group.
func (*NoDuplicateIfConditionRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoDuplicateIfConditionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoDupIfCond struct {
	onFailure func(lint.Failure)
	visited   map[*ast.IfStmt]bool
}

func (w *lintNoDupIfCond) Visit(node ast.Node) ast.Visitor {
	ifStmt, ok := node.(*ast.IfStmt)
	if !ok {
		return w
	}

	// Skip if this if-stmt was already processed as part of
	// an earlier chain (i.e., it is an else-if branch).
	if w.visited[ifStmt] {
		return w
	}

	// Walk the full if/else if chain and report duplicates.
	w.checkIfChain(ifStmt)

	return w
}

// checkIfChain walks the if/else if chain starting at ifStmt and reports
// any condition that was already seen earlier in the chain.
func (w *lintNoDupIfCond) checkIfChain(ifStmt *ast.IfStmt) {
	seen := map[string]*ast.IfStmt{}

	for current := ifStmt; current != nil; {
		if current.Cond != nil {
			key := astutils.GoFmt(current.Cond)
			if key != "" {
				if _, exists := seen[key]; exists {
					w.onFailure(lint.Failure{
						Category:   lint.FailureCategoryLogic,
						Confidence: 1,
						Node:       current,
						Failure:    fmt.Sprintf("duplicate condition %s in if/else if chain", key),
					})
				} else {
					seen[key] = current
				}
			}
		}

		// Follow the else branch if it's another if statement.
		next, ok := current.Else.(*ast.IfStmt)
		if !ok {
			break
		}
		// Mark the else-if as visited so we don't re-process it as a chain start.
		w.visited[next] = true
		current = next
	}
}
