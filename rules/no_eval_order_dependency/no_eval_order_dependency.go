package no_eval_order_dependency

import (
	"go/ast"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoEvalOrderDependencyRule detects unwanted dependencies on evaluation order
// in return statements. It flags cases where a return value is a variable that
// may be modified by a function call in another return value (e.g., via address-of
// or pointer receiver).
type NoEvalOrderDependencyRule struct{}

// Apply applies the rule to given file.
func (r *NoEvalOrderDependencyRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoEvalOrderDep{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoEvalOrderDependencyRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoEvalOrderDep{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoEvalOrderDependencyRule) Name() string {
	return "noEvalOrderDependency"
}

// Group returns the rule group.
func (*NoEvalOrderDependencyRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoEvalOrderDependencyRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoEvalOrderDep struct {
	onFailure func(lint.Failure)
}

func (w *lintNoEvalOrderDep) Visit(node ast.Node) ast.Visitor {
	ret, ok := node.(*ast.ReturnStmt)
	if !ok {
		return w
	}

	if len(ret.Results) < 2 {
		return w
	}

	// Collect all variable names that appear as plain identifiers in return values
	// and all variable names whose addresses are taken by function calls in other return values.
	// If a variable appears both as a plain return value and has its address taken
	// in a function call in another return value, flag the return statement.

	// First pass: collect all plain identifier names used as return values and their indices
	type identInfo struct {
		name  string
		index int
	}
	var returnedIdents []identInfo
	for i, result := range ret.Results {
		if ident, ok := result.(*ast.Ident); ok {
			returnedIdents = append(returnedIdents, identInfo{name: ident.Name, index: i})
		}
	}

	if len(returnedIdents) == 0 {
		return w
	}

	// Second pass: for each return value that is a call expression,
	// check if any argument takes the address of a returned variable
	for i, result := range ret.Results {
		addrTakenVars := collectAddrTakenVars(result)
		for _, ri := range returnedIdents {
			if ri.index == i {
				continue // same position, skip
			}
			if addrTakenVars[ri.name] {
				w.onFailure(lint.Failure{
					Category:   lint.FailureCategoryLogic,
					Confidence: 0.8,
					Node:       ret,
					Failure:    "return value may depend on evaluation order: variable " + ri.name + " may be modified by a function call in another return value",
				})
				return w // only report once per return statement
			}
		}
	}

	return w
}

// collectAddrTakenVars walks an expression and returns the set of variable names
// whose addresses are taken (via &var) inside function call arguments.
func collectAddrTakenVars(expr ast.Expr) map[string]bool {
	result := make(map[string]bool)
	ast.Inspect(expr, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		for _, arg := range call.Args {
			unary, ok := arg.(*ast.UnaryExpr)
			if !ok {
				continue
			}
			if unary.Op.String() != "&" {
				continue
			}
			if ident, ok := unary.X.(*ast.Ident); ok {
				result[ident.Name] = true
			}
		}
		return true
	})
	return result
}
