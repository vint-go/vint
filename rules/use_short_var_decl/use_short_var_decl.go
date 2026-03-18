package use_short_var_decl

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseShortVarDeclRule detects suspicious re-assignments in if statement init
// blocks where a short variable declaration (:=) would be more appropriate.
type UseShortVarDeclRule struct{}

// Apply applies the rule to given file.
func (r *UseShortVarDeclRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	w := &lintShortVarDecl{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseShortVarDeclRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	w := &lintShortVarDecl{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseShortVarDeclRule) Name() string {
	return "useShortVarDecl"
}

// Group returns the rule group.
func (*UseShortVarDeclRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseShortVarDeclRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintShortVarDecl struct {
	onFailure func(lint.Failure)
}

func (w *lintShortVarDecl) Visit(node ast.Node) ast.Visitor {
	ifStmt, ok := node.(*ast.IfStmt)
	if !ok {
		return w
	}

	// Must have an init statement
	if ifStmt.Init == nil {
		return w
	}

	// The init must be an assignment statement with = (not :=)
	assignStmt, ok := ifStmt.Init.(*ast.AssignStmt)
	if !ok {
		return w
	}
	if assignStmt.Tok != token.ASSIGN {
		return w
	}

	// Must have exactly one LHS and one RHS
	if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		return w
	}

	// LHS must be an identifier
	lhsIdent, ok := assignStmt.Lhs[0].(*ast.Ident)
	if !ok {
		return w
	}

	// The condition must be a binary expression checking != nil
	binExpr, ok := ifStmt.Cond.(*ast.BinaryExpr)
	if !ok {
		return w
	}
	if binExpr.Op != token.NEQ {
		return w
	}

	// One side of the condition must be the same identifier, the other nil
	condIdent, condNil := extractIdentAndNil(binExpr)
	if condIdent == nil || !condNil {
		return w
	}
	if condIdent.Name != lhsIdent.Name {
		return w
	}

	// The body must contain a return statement that returns the same variable
	if ifStmt.Body == nil || len(ifStmt.Body.List) == 0 {
		return w
	}

	if !bodyReturnsIdent(ifStmt.Body, lhsIdent.Name) {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Category:   lint.FailureCategoryStyle,
		Node:       ifStmt,
		Failure:    "use short variable declaration (:=) instead of assignment (=) in if init",
	})

	return w
}

// extractIdentAndNil checks if a binary expression has an identifier on one side
// and nil on the other. Returns the identifier and true if nil was found.
func extractIdentAndNil(expr *ast.BinaryExpr) (*ast.Ident, bool) {
	xIdent, xIsIdent := expr.X.(*ast.Ident)
	yIdent, yIsIdent := expr.Y.(*ast.Ident)

	if xIsIdent && yIsIdent && yIdent.Name == "nil" {
		return xIdent, true
	}
	if yIsIdent && xIsIdent && xIdent.Name == "nil" {
		return yIdent, true
	}
	// Handle case where X is ident and Y is nil (non-ident nil is unusual but check anyway)
	if xIsIdent && !yIsIdent {
		return nil, false
	}
	if yIsIdent && !xIsIdent {
		return nil, false
	}
	return nil, false
}

// bodyReturnsIdent checks if the block statement contains a return statement
// that returns the named identifier.
func bodyReturnsIdent(body *ast.BlockStmt, name string) bool {
	for _, stmt := range body.List {
		retStmt, ok := stmt.(*ast.ReturnStmt)
		if !ok {
			continue
		}
		for _, result := range retStmt.Results {
			ident, ok := result.(*ast.Ident)
			if ok && ident.Name == name {
				return true
			}
		}
	}
	return false
}
