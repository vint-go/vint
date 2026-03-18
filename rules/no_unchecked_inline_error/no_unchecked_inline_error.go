package no_unchecked_inline_error

import (
	"go/ast"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoUncheckedInlineErrorRule detects unchecked errors in if statement
// initialization clauses. It flags cases where an error variable is assigned
// in the init clause but not referenced in the condition.
type NoUncheckedInlineErrorRule struct{}

// Apply applies the rule to the given file.
func (r *NoUncheckedInlineErrorRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintUncheckedInlineError{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoUncheckedInlineErrorRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	w := &lintUncheckedInlineError{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoUncheckedInlineErrorRule) Name() string {
	return "noUncheckedInlineError"
}

// Group returns the rule group.
func (*NoUncheckedInlineErrorRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoUncheckedInlineErrorRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintUncheckedInlineError struct {
	onFailure func(lint.Failure)
}

func (w *lintUncheckedInlineError) Visit(node ast.Node) ast.Visitor {
	ifStmt, ok := node.(*ast.IfStmt)
	if !ok {
		return w
	}

	if ifStmt.Init == nil {
		return w
	}

	assign, ok := ifStmt.Init.(*ast.AssignStmt)
	if !ok {
		return w
	}

	// Collect names of error-like variables from the LHS of the assignment.
	// We look for identifiers whose name contains "err" (case-insensitive match
	// on the common convention of naming error variables "err" or "Err...").
	var errVarNames []string
	for _, lhs := range assign.Lhs {
		ident, ok := lhs.(*ast.Ident)
		if !ok {
			continue
		}
		name := ident.Name
		if name == "_" {
			continue
		}
		if isErrorName(name) {
			errVarNames = append(errVarNames, name)
		}
	}

	if len(errVarNames) == 0 {
		return w
	}

	// Check whether any of the error variable names appear in the condition.
	condIdents := collectIdents(ifStmt.Cond)
	for _, errName := range errVarNames {
		if !condIdents[errName] {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       ifStmt,
				Category:   lint.FailureCategoryErrors,
				Failure:    "error variable '" + errName + "' assigned in if initialization but not checked in the condition",
			})
			break
		}
	}

	return w
}

// isErrorName returns true if the identifier name looks like an error variable.
// It matches "err" exactly or names starting with "err" (like "err2", "errFoo").
func isErrorName(name string) bool {
	if name == "err" {
		return true
	}
	if len(name) > 3 && name[:3] == "err" {
		// e.g. err2, errFoo, errBar
		return true
	}
	return false
}

// collectIdents walks an expression and returns a set of all identifier names.
func collectIdents(expr ast.Expr) map[string]bool {
	idents := map[string]bool{}
	ast.Inspect(expr, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			idents[ident.Name] = true
		}
		return true
	})
	return idents
}
