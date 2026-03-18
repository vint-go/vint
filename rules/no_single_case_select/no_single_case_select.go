package no_single_case_select

import (
	"go/ast"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoSingleCaseSelectRule detects select statements with a single case and no
// default clause. Such select statements are equivalent to a plain channel
// send or receive and the select adds unnecessary complexity.
type NoSingleCaseSelectRule struct{}

// Apply applies the rule to given file.
func (r *NoSingleCaseSelectRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoSingleCaseSelect{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoSingleCaseSelectRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoSingleCaseSelect{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoSingleCaseSelectRule) Name() string {
	return "noSingleCaseSelect"
}

// Group returns the rule group.
func (*NoSingleCaseSelectRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoSingleCaseSelectRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoSingleCaseSelect struct {
	onFailure func(lint.Failure)
}

func (w *lintNoSingleCaseSelect) Visit(node ast.Node) ast.Visitor {
	selectStmt, ok := node.(*ast.SelectStmt)
	if !ok {
		return w
	}

	if selectStmt.Body == nil {
		return w
	}

	clauses := selectStmt.Body.List
	if len(clauses) != 1 {
		return w
	}

	cc, ok := clauses[0].(*ast.CommClause)
	if !ok {
		return w
	}

	// Only flag if there is exactly one case (not a default clause).
	// A default clause has cc.Comm == nil.
	if cc.Comm == nil {
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       selectStmt,
		Failure:    "use plain channel send or receive instead of single-case select",
	})

	return w
}
