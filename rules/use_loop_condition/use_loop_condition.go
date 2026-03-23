package use_loop_condition

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseLoopConditionRule detects for loops where the first statement is an
// if+break that could be lifted into the loop condition.
//
// Example:
//
//	for {
//	    if i >= len(items) {
//	        break
//	    }
//	    ...
//	}
//
// can be rewritten as:
//
//	for i < len(items) {
//	    ...
//	}
type UseLoopConditionRule struct{}

// Apply applies the rule to the given file.
func (r *UseLoopConditionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	w := &lintUseLoopCondition{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseLoopConditionRule) ApplyToNode(_ *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	w := &lintUseLoopCondition{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseLoopConditionRule) Name() string {
	return "useLoopCondition"
}

// Group returns the rule group.
func (*UseLoopConditionRule) Group() string {
	return "complexity"
}

// CacheTier returns the cache tier for this rule.
func (*UseLoopConditionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintUseLoopCondition struct {
	onFailure func(lint.Failure)
}

func (w *lintUseLoopCondition) Visit(node ast.Node) ast.Visitor {
	forStmt, ok := node.(*ast.ForStmt)
	if !ok {
		return w
	}

	// Only match for loops that have no condition (infinite loops: `for { ... }`)
	if forStmt.Cond != nil {
		return w
	}

	body := forStmt.Body
	if body == nil || len(body.List) < 1 {
		return w
	}

	// The first statement must be an if-statement
	ifStmt, ok := body.List[0].(*ast.IfStmt)
	if !ok {
		return w
	}

	// The if must have no init statement and no else clause
	if ifStmt.Init != nil {
		return w
	}
	if ifStmt.Else != nil {
		return w
	}

	// The if body must contain exactly one statement: a break (with no label)
	if ifStmt.Body == nil || len(ifStmt.Body.List) != 1 {
		return w
	}
	branchStmt, ok := ifStmt.Body.List[0].(*ast.BranchStmt)
	if !ok {
		return w
	}
	if branchStmt.Tok.String() != "break" || branchStmt.Label != nil {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Category:   lint.FailureCategoryComplexity,
		Failure:    "lift if+break into loop condition",
		Node:       ifStmt,
	})

	return w
}
