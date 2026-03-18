package no_invariant_loop_condition

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoInvariantLoopConditionRule detects for loops where the condition variable
// is never modified within the loop body, leading to an infinite or zero-iteration loop.
type NoInvariantLoopConditionRule struct{}

// Apply applies the rule to given file.
func (r *NoInvariantLoopConditionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintInvariantLoop{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoInvariantLoopConditionRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintInvariantLoop{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoInvariantLoopConditionRule) Name() string {
	return "noInvariantLoopCondition"
}

// Group returns the rule group.
func (*NoInvariantLoopConditionRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoInvariantLoopConditionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintInvariantLoop struct {
	onFailure func(lint.Failure)
}

func (w *lintInvariantLoop) Visit(node ast.Node) ast.Visitor {
	forStmt, ok := node.(*ast.ForStmt)
	if !ok {
		return w
	}

	// Only check for loops with a condition and no post statement
	// (loops with a post statement like i++ typically modify condition variables there)
	if forStmt.Cond == nil {
		return w
	}

	// Collect all identifiers referenced in the condition
	condVars := collectIdents(forStmt.Cond)
	if len(condVars) == 0 {
		return w
	}

	// Check if any condition variable is modified in:
	// 1. The post statement (e.g., i++ in for ; i < n; i++)
	// 2. The loop body
	if forStmt.Post != nil && modifiesAny(forStmt.Post, condVars) {
		return w
	}

	if forStmt.Body != nil && modifiesAny(forStmt.Body, condVars) {
		return w
	}

	// Check if the loop body contains any statements that could affect
	// control flow or have side effects that we can't analyze (channel receives,
	// function calls that might modify variables via pointers, etc.)
	if forStmt.Body != nil && hasUnanalyzableEffect(forStmt.Body, condVars) {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 0.8,
		Node:       forStmt,
		Category:   lint.FailureCategoryLogic,
		Failure:    "loop condition variable is never modified in the loop body",
	})

	return w
}

// collectIdents collects all identifier names referenced in an expression.
func collectIdents(expr ast.Expr) map[string]bool {
	idents := make(map[string]bool)
	ast.Inspect(expr, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			// Skip built-in constants and common function names
			if ident.Name != "true" && ident.Name != "false" && ident.Name != "nil" {
				idents[ident.Name] = true
			}
		}
		return true
	})
	return idents
}

// modifiesAny checks if a statement modifies any of the given variables.
func modifiesAny(node ast.Node, vars map[string]bool) bool {
	modified := false
	ast.Inspect(node, func(n ast.Node) bool {
		if modified {
			return false
		}
		switch stmt := n.(type) {
		case *ast.AssignStmt:
			for _, lhs := range stmt.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok {
					if vars[ident.Name] {
						modified = true
						return false
					}
				}
			}
		case *ast.IncDecStmt:
			if ident, ok := stmt.X.(*ast.Ident); ok {
				if vars[ident.Name] {
					modified = true
					return false
				}
			}
		case *ast.RangeStmt:
			// Range key/value can modify variables
			if stmt.Key != nil {
				if ident, ok := stmt.Key.(*ast.Ident); ok {
					if vars[ident.Name] {
						modified = true
						return false
					}
				}
			}
			if stmt.Value != nil {
				if ident, ok := stmt.Value.(*ast.Ident); ok {
					if vars[ident.Name] {
						modified = true
						return false
					}
				}
			}
		}
		return true
	})
	return modified
}

// hasUnanalyzableEffect checks if the loop body contains operations that could
// modify condition variables through indirect means (pointers, closures, etc.)
func hasUnanalyzableEffect(body *ast.BlockStmt, condVars map[string]bool) bool {
	hasEffect := false
	ast.Inspect(body, func(n ast.Node) bool {
		if hasEffect {
			return false
		}
		switch expr := n.(type) {
		case *ast.UnaryExpr:
			// Taking address of a condition variable means it could be modified indirectly
			if expr.Op == token.AND {
				if ident, ok := expr.X.(*ast.Ident); ok {
					if condVars[ident.Name] {
						hasEffect = true
						return false
					}
				}
			}
		case *ast.CallExpr:
			// Check if any condition variable is passed to a function call
			for _, arg := range expr.Args {
				if containsCondVar(arg, condVars) {
					hasEffect = true
					return false
				}
			}
		case *ast.SendStmt:
			// Channel sends with condition variables
			if containsCondVar(expr.Value, condVars) {
				hasEffect = true
				return false
			}
		}
		return true
	})
	return hasEffect
}

// containsCondVar checks if an expression contains a reference to any condition variable
// via address-of operator.
func containsCondVar(expr ast.Expr, condVars map[string]bool) bool {
	found := false
	ast.Inspect(expr, func(n ast.Node) bool {
		if found {
			return false
		}
		if unary, ok := n.(*ast.UnaryExpr); ok && unary.Op == token.AND {
			if ident, ok := unary.X.(*ast.Ident); ok {
				if condVars[ident.Name] {
					found = true
					return false
				}
			}
		}
		return true
	})
	return found
}
