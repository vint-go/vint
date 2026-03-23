package no_misplaced_default_case

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoMisplacedDefaultCaseRule detects when default case in switch isn't on first or last position.
type NoMisplacedDefaultCaseRule struct{}

// Apply applies the rule to given file.
func (r *NoMisplacedDefaultCaseRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintMisplacedDefault{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoMisplacedDefaultCaseRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintMisplacedDefault{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoMisplacedDefaultCaseRule) Name() string {
	return "noMisplacedDefaultCase"
}

// Group returns the rule group.
func (*NoMisplacedDefaultCaseRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoMisplacedDefaultCaseRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintMisplacedDefault struct {
	onFailure func(lint.Failure)
}

func (w *lintMisplacedDefault) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.SwitchStmt:
		w.checkCaseList(n.Body.List)
	case *ast.TypeSwitchStmt:
		w.checkCaseList(n.Body.List)
	}
	return w
}

func (w *lintMisplacedDefault) checkCaseList(cases []ast.Stmt) {
	if len(cases) <= 2 {
		// With 0, 1 or 2 cases the default can only be first or last.
		return
	}

	for i, stmt := range cases {
		cc, ok := stmt.(*ast.CaseClause)
		if !ok {
			continue
		}
		// A nil List means this is a default case.
		if cc.List == nil {
			if i != 0 && i != len(cases)-1 {
				w.onFailure(lint.Failure{
					Confidence: 1,
					Node:       cc,
					Category:   lint.FailureCategoryStyle,
					Failure:    "default case should be the first or last case in the switch",
				})
			}
		}
	}
}
