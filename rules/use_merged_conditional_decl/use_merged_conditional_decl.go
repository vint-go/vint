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
		varName, initBool, ok := extractSingleBoolVarAssign(stmts[i])
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

		// The reassignment must be the opposite boolean literal
		reassignBool, ok := isBoolLiteral(assignStmt.Rhs[0])
		if !ok || reassignBool != oppositeBool(initBool) {
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

// extractSingleBoolVarAssign returns the variable name and the boolean literal value
// if the statement is a short variable declaration (:=) or var declaration with
// exactly one name and one boolean literal value (true or false).
func extractSingleBoolVarAssign(stmt ast.Stmt) (string, string, bool) {
	// Check for short variable declaration: x := true
	if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
		if assignStmt.Tok != token.DEFINE {
			return "", "", false
		}
		if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
			return "", "", false
		}
		ident, ok := assignStmt.Lhs[0].(*ast.Ident)
		if !ok {
			return "", "", false
		}
		boolVal, ok := isBoolLiteral(assignStmt.Rhs[0])
		if !ok {
			return "", "", false
		}
		return ident.Name, boolVal, true
	}

	// Check for var declaration: var x = true
	declStmt, ok := stmt.(*ast.DeclStmt)
	if !ok {
		return "", "", false
	}
	genDecl, ok := declStmt.Decl.(*ast.GenDecl)
	if !ok || genDecl.Tok != token.VAR {
		return "", "", false
	}
	if len(genDecl.Specs) != 1 {
		return "", "", false
	}
	valueSpec, ok := genDecl.Specs[0].(*ast.ValueSpec)
	if !ok {
		return "", "", false
	}
	// Must declare exactly one variable WITH an initial value
	if len(valueSpec.Names) != 1 || len(valueSpec.Values) != 1 {
		return "", "", false
	}
	boolVal, ok := isBoolLiteral(valueSpec.Values[0])
	if !ok {
		return "", "", false
	}
	return valueSpec.Names[0].Name, boolVal, true
}

// isBoolLiteral checks if the expression is a boolean literal (true or false)
// and returns the literal name.
func isBoolLiteral(expr ast.Expr) (string, bool) {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return "", false
	}
	if ident.Name == "true" || ident.Name == "false" {
		return ident.Name, true
	}
	return "", false
}

// oppositeBool returns the opposite boolean literal name.
func oppositeBool(b string) string {
	if b == "true" {
		return "false"
	}
	return "true"
}
