package no_redundant_label

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoRedundantLabelRule detects redundant statement labels.
// It identifies labels on statements where no nested break/continue
// statements reference that label.
type NoRedundantLabelRule struct{}

// Apply applies the rule to given file.
func (r *NoRedundantLabelRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	w := &lintRedundantLabel{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoRedundantLabelRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	w := &lintRedundantLabel{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoRedundantLabelRule) Name() string {
	return "noRedundantLabel"
}

// Group returns the rule group.
func (*NoRedundantLabelRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoRedundantLabelRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintRedundantLabel struct {
	onFailure func(lint.Failure)
}

func (w *lintRedundantLabel) Visit(node ast.Node) ast.Visitor {
	labelStmt, ok := node.(*ast.LabeledStmt)
	if !ok {
		return w
	}

	labelName := labelStmt.Label.Name

	// Check if the labeled statement is a loop or switch/select/type-switch
	// (these are the only statements where break/continue with a label makes sense).
	if !isLabelableStmt(labelStmt.Stmt) {
		// Labels on non-loop/non-switch statements can still be used by goto,
		// so we skip those to avoid false positives.
		return w
	}

	// Collect all break/continue label references within the labeled statement.
	refs := collectLabelRefs(labelStmt.Stmt)

	if !refs[labelName] {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Category:   lint.FailureCategoryStyle,
			Node:       labelStmt,
			Failure:    fmt.Sprintf("label %q is redundant", labelName),
		})
	}

	return w
}

// isLabelableStmt returns true if the statement is one where a label
// can be referenced by break or continue (loops, switch, select).
func isLabelableStmt(stmt ast.Stmt) bool {
	switch stmt.(type) {
	case *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
		return true
	}
	return false
}

// collectLabelRefs walks the statement tree and returns a set of label names
// referenced by break/continue statements. It does NOT descend into nested
// function literals (closures), since labels cannot cross function boundaries.
func collectLabelRefs(node ast.Node) map[string]bool {
	refs := map[string]bool{}
	ast.Inspect(node, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncLit:
			// Don't descend into closures; labels don't cross function boundaries.
			return false
		case *ast.BranchStmt:
			if n.Label != nil && (n.Tok == token.BREAK || n.Tok == token.CONTINUE) {
				refs[n.Label.Name] = true
			}
		}
		return true
	})
	return refs
}
