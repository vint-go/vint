package no_empty_fallthrough

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoEmptyFallthroughRule detects switch case clauses that contain only a
// fallthrough statement with no other logic. Such cases can be simplified
// by combining the case values into a comma-separated list.
type NoEmptyFallthroughRule struct{}

// Apply applies the rule to given file.
func (r *NoEmptyFallthroughRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoEmptyFallthrough{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoEmptyFallthroughRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoEmptyFallthrough{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoEmptyFallthroughRule) Name() string {
	return "noEmptyFallthrough"
}

// Group returns the rule group.
func (*NoEmptyFallthroughRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoEmptyFallthroughRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoEmptyFallthrough struct {
	onFailure func(lint.Failure)
}

func (w *lintNoEmptyFallthrough) Visit(node ast.Node) ast.Visitor {
	switchStmt, ok := node.(*ast.SwitchStmt)
	if !ok {
		return w
	}

	if switchStmt.Body == nil {
		return w
	}

	for _, item := range switchStmt.Body.List {
		cc, ok := item.(*ast.CaseClause)
		if !ok {
			continue
		}

		// Check if this case clause has exactly one statement and it is a fallthrough
		if len(cc.Body) == 1 {
			branchStmt, ok := cc.Body[0].(*ast.BranchStmt)
			if ok && branchStmt.Tok == token.FALLTHROUGH {
				w.onFailure(lint.Failure{
					Category:   lint.FailureCategoryStyle,
					Confidence: 1,
					Node:       cc,
					Failure:    "case clause with only a fallthrough can be combined with the next case",
				})
			}
		}
	}

	return w
}
