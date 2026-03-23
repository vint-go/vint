package use_merged_conditional_decl

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseMergedConditionalDeclRule detects variable declarations followed by a
// conditional assignment that can be merged into a single declaration with a
// conditional expression.
//
// Source: https://staticcheck.dev/docs/checks/#QF1007
type UseMergedConditionalDeclRule struct{}

// Apply applies the rule to the given file.
func (r *UseMergedConditionalDeclRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	w := &lintMergedConditionalDecl{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseMergedConditionalDeclRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	w := &lintMergedConditionalDecl{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseMergedConditionalDeclRule) Name() string {
	return "useMergedConditionalDecl"
}

// Group returns the rule group.
func (*UseMergedConditionalDeclRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseMergedConditionalDeclRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintMergedConditionalDecl struct {
	onFailure func(lint.Failure)
}

func (w *lintMergedConditionalDecl) Visit(node ast.Node) ast.Visitor {
	block, ok := node.(*ast.BlockStmt)
	if !ok {
		return w
	}

	w.checkBlock(block.List)
	return w
}

// checkBlock inspects a list of statements for consecutive variable-declaration +
// conditional-reassignment pairs that could be merged.
func (w *lintMergedConditionalDecl) checkBlock(stmts []ast.Stmt) {
	for i := 0; i < len(stmts)-1; i++ {
		varName, ok := extractSingleVarAssign(stmts[i])
		if !ok {
			continue
		}

		ifStmt, ok := stmts[i+1].(*ast.IfStmt)
		if !ok {
			continue
		}

		// The if must not have an init statement
		if ifStmt.Init != nil {
			continue
		}

		// The if must not have an else branch
		if ifStmt.Else != nil {
			continue
		}

		// The if body must contain exactly one statement
		if len(ifStmt.Body.List) != 1 {
			continue
		}

		// That single statement must be an assignment to the same variable
		assignStmt, ok := ifStmt.Body.List[0].(*ast.AssignStmt)
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

		w.onFailure(lint.Failure{
			Confidence: 1,
			Category:   lint.FailureCategoryStyle,
			Node:       stmts[i],
			Failure:    fmt.Sprintf("merge conditional assignment into variable declaration of %s", varName),
		})
	}
}

// extractSingleVarAssign returns the variable name if the statement is a short
// variable declaration (:=) or var declaration with exactly one name and one value.
func extractSingleVarAssign(stmt ast.Stmt) (string, bool) {
	// Check for short variable declaration: x := "value"
	if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
		if assignStmt.Tok != token.DEFINE {
			return "", false
		}
		if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
			return "", false
		}
		ident, ok := assignStmt.Lhs[0].(*ast.Ident)
		if !ok {
			return "", false
		}
		return ident.Name, true
	}

	// Check for var declaration: var x = "value"
	declStmt, ok := stmt.(*ast.DeclStmt)
	if !ok {
		return "", false
	}
	genDecl, ok := declStmt.Decl.(*ast.GenDecl)
	if !ok || genDecl.Tok != token.VAR {
		return "", false
	}
	if len(genDecl.Specs) != 1 {
		return "", false
	}
	valueSpec, ok := genDecl.Specs[0].(*ast.ValueSpec)
	if !ok {
		return "", false
	}
	// Must declare exactly one variable WITH an initial value
	if len(valueSpec.Names) != 1 || len(valueSpec.Values) != 1 {
		return "", false
	}
	return valueSpec.Names[0].Name, true
}
