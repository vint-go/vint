package use_merged_var_decl

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseMergedVarDeclRule detects variable declarations followed immediately by
// an assignment that can be combined into a short variable declaration.
type UseMergedVarDeclRule struct{}

// Apply applies the rule to given file.
func (r *UseMergedVarDeclRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	w := &lintMergedVarDecl{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseMergedVarDeclRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	w := &lintMergedVarDecl{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseMergedVarDeclRule) Name() string {
	return "useMergedVarDecl"
}

// Group returns the rule group.
func (*UseMergedVarDeclRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseMergedVarDeclRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintMergedVarDecl struct {
	onFailure func(lint.Failure)
}

func (w *lintMergedVarDecl) Visit(node ast.Node) ast.Visitor {
	block, ok := node.(*ast.BlockStmt)
	if !ok {
		return w
	}

	w.checkBlock(block.List)
	return w
}

// checkBlock inspects a list of statements for consecutive var-decl + assignment
// pairs that could be merged into a short variable declaration.
func (w *lintMergedVarDecl) checkBlock(stmts []ast.Stmt) {
	for i := 0; i < len(stmts)-1; i++ {
		// First statement must be a var declaration (ast.DeclStmt with *ast.GenDecl of token.VAR)
		declStmt, ok := stmts[i].(*ast.DeclStmt)
		if !ok {
			continue
		}
		genDecl, ok := declStmt.Decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}
		// Must have exactly one spec
		if len(genDecl.Specs) != 1 {
			continue
		}
		valueSpec, ok := genDecl.Specs[0].(*ast.ValueSpec)
		if !ok {
			continue
		}
		// Must declare exactly one variable with no initial value
		if len(valueSpec.Names) != 1 || len(valueSpec.Values) != 0 {
			continue
		}
		varName := valueSpec.Names[0].Name

		// Second statement must be an assignment to that same variable
		assignStmt, ok := stmts[i+1].(*ast.AssignStmt)
		if !ok {
			continue
		}
		if assignStmt.Tok != token.ASSIGN {
			continue
		}
		if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
			continue
		}
		lhsIdent, ok := assignStmt.Lhs[0].(*ast.Ident)
		if !ok {
			continue
		}
		if lhsIdent.Name != varName {
			continue
		}

		// Skip if the variable is assigned multiple times in the block.
		if hasMultipleAssignments(stmts, varName) {
			continue
		}

		// Skip if the RHS references the declared variable (self-referential).
		if refersToIdent(assignStmt.Rhs[0], varName) {
			continue
		}

		w.onFailure(lint.Failure{
			Confidence: 1,
			Category:   lint.FailureCategoryStyle,
			Node:       declStmt,
			Failure:    fmt.Sprintf("should merge variable declaration with assignment on next line to %s := ...", varName),
		})
	}
}

// hasMultipleAssignments checks whether varName appears on the LHS of
// assignments two or more times within the given statements. This prevents
// false positives when a variable is conditionally overridden or shared with
// a closure and reassigned inside it.
func hasMultipleAssignments(stmts []ast.Stmt, varName string) bool {
	count := 0
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			if count >= 2 {
				return false
			}
			assign, ok := n.(*ast.AssignStmt)
			if !ok {
				return true
			}
			for _, lhs := range assign.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok && ident.Name == varName {
					count++
				}
			}
			return true
		})
	}
	return count >= 2
}

// refersToIdent reports whether expr contains a reference to an identifier
// with the given name.
func refersToIdent(expr ast.Expr, name string) bool {
	found := false
	ast.Inspect(expr, func(n ast.Node) bool {
		if found {
			return false
		}
		if ident, ok := n.(*ast.Ident); ok && ident.Name == name {
			found = true
			return false
		}
		return true
	})
	return found
}
