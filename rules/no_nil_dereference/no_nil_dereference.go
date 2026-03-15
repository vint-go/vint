package no_nil_dereference

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoNilDereferenceRule reports nil pointer dereferences and degenerate nil comparisons.
// It detects cases where a pointer is dereferenced after being confirmed nil, and
// comparisons against nil where the result is always the same.
type NoNilDereferenceRule struct{}

// Apply applies the rule to given file.
func (r *NoNilDereferenceRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNilDeref{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}
		w.checkFunction(funcDecl)
	}

	return failures
}

// Name returns the rule name.
func (*NoNilDereferenceRule) Name() string {
	return "noNilDereference"
}

// Group returns the rule group.
func (*NoNilDereferenceRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoNilDereferenceRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNilDeref struct {
	onFailure func(lint.Failure)
}

// checkFunction examines a function for nil dereference and degenerate nil comparison issues.
func (w *lintNilDeref) checkFunction(funcDecl *ast.FuncDecl) {
	// Check for nil dereferences inside nil-confirmed branches
	w.checkNilDerefInBranches(funcDecl.Body)

	// Check for degenerate nil comparisons (variables that are always nil)
	w.checkDegenerateNilComparisons(funcDecl)
}

// checkNilDerefInBranches looks for patterns like:
//
//	if p == nil { ... *p ... }
//
// where a pointer is dereferenced after being confirmed nil.
func (w *lintNilDeref) checkNilDerefInBranches(body *ast.BlockStmt) {
	ast.Inspect(body, func(n ast.Node) bool {
		ifStmt, ok := n.(*ast.IfStmt)
		if !ok {
			return true
		}

		// Check if the condition is a nil comparison: p == nil
		varName, isNilCheck := extractNilEqualityCheck(ifStmt.Cond)
		if !isNilCheck {
			return true
		}

		// Now look for dereferences of this variable in the if-body (the nil-confirmed branch)
		w.findDerefsInBlock(ifStmt.Body, varName)

		return true
	})
}

// checkDegenerateNilComparisons looks for patterns like:
//
//	var p *int
//	if p != nil { ... }
//
// where a pointer variable is declared as nil (never assigned a non-nil value)
// and then compared to nil.
func (w *lintNilDeref) checkDegenerateNilComparisons(funcDecl *ast.FuncDecl) {
	// Collect all pointer variables declared with var (no initializer) in the function body
	nilVars := map[string]bool{}

	for _, stmt := range funcDecl.Body.List {
		declStmt, ok := stmt.(*ast.DeclStmt)
		if !ok {
			continue
		}
		genDecl, ok := declStmt.Decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}
		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			// Only consider pointer types with no initializer values
			if len(valueSpec.Values) > 0 {
				continue
			}
			if !isPointerType(valueSpec.Type) {
				continue
			}
			for _, name := range valueSpec.Names {
				nilVars[name.Name] = true
			}
		}
	}

	if len(nilVars) == 0 {
		return
	}

	// Now remove variables that get assigned a non-nil value before a nil comparison
	w.removeAssignedVars(funcDecl.Body, nilVars)

	if len(nilVars) == 0 {
		return
	}

	// Find nil comparisons for these always-nil variables
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		ifStmt, ok := n.(*ast.IfStmt)
		if !ok {
			return true
		}

		varName, isNotNilCheck := extractNilNotEqualCheck(ifStmt.Cond)
		if !isNotNilCheck {
			return true
		}

		if nilVars[varName] {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 1,
				Node:       ifStmt.Cond,
				Failure:    fmt.Sprintf("degenerate nil comparison: %s is always nil", varName),
			})
		}

		return true
	})
}

// findDerefsInBlock looks for dereferences of the named variable in a block.
func (w *lintNilDeref) findDerefsInBlock(block *ast.BlockStmt, varName string) {
	ast.Inspect(block, func(n ast.Node) bool {
		switch expr := n.(type) {
		case *ast.StarExpr:
			// *p
			if ident, ok := expr.X.(*ast.Ident); ok && ident.Name == varName {
				w.onFailure(lint.Failure{
					Category:   lint.FailureCategoryLogic,
					Confidence: 1,
					Node:       expr,
					Failure:    fmt.Sprintf("nil pointer dereference: %s is nil in this branch", varName),
				})
			}
		case *ast.SelectorExpr:
			// p.Field or p.Method()
			if ident, ok := expr.X.(*ast.Ident); ok && ident.Name == varName {
				w.onFailure(lint.Failure{
					Category:   lint.FailureCategoryLogic,
					Confidence: 1,
					Node:       expr,
					Failure:    fmt.Sprintf("nil pointer dereference: %s is nil in this branch", varName),
				})
			}
		}
		return true
	})
}

// extractNilEqualityCheck checks if the expression is of the form `x == nil`
// and returns the variable name.
func extractNilEqualityCheck(expr ast.Expr) (string, bool) {
	binExpr, ok := expr.(*ast.BinaryExpr)
	if !ok || binExpr.Op != token.EQL {
		return "", false
	}
	return extractNilComparisonVarName(binExpr)
}

// extractNilNotEqualCheck checks if the expression is of the form `x != nil`
// and returns the variable name.
func extractNilNotEqualCheck(expr ast.Expr) (string, bool) {
	binExpr, ok := expr.(*ast.BinaryExpr)
	if !ok || binExpr.Op != token.NEQ {
		return "", false
	}
	return extractNilComparisonVarName(binExpr)
}

// extractNilComparisonVarName extracts the variable name from a binary expression
// where one side is nil and the other is an identifier.
func extractNilComparisonVarName(binExpr *ast.BinaryExpr) (string, bool) {
	// Check x == nil
	if ident, ok := binExpr.X.(*ast.Ident); ok {
		if isNilIdent(binExpr.Y) {
			return ident.Name, true
		}
	}
	// Check nil == x
	if ident, ok := binExpr.Y.(*ast.Ident); ok {
		if isNilIdent(binExpr.X) {
			return ident.Name, true
		}
	}
	return "", false
}

// isNilIdent checks if the expression is the nil identifier.
func isNilIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "nil"
}

// isPointerType checks if the type expression represents a pointer type.
func isPointerType(expr ast.Expr) bool {
	_, ok := expr.(*ast.StarExpr)
	return ok
}

// removeAssignedVars removes variable names from the map if they are assigned
// any value in the block (i.e., they are no longer guaranteed nil).
func (w *lintNilDeref) removeAssignedVars(block *ast.BlockStmt, nilVars map[string]bool) {
	ast.Inspect(block, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}
		for _, lhs := range assign.Lhs {
			if ident, ok := lhs.(*ast.Ident); ok {
				delete(nilVars, ident.Name)
			}
		}
		return true
	})
}
