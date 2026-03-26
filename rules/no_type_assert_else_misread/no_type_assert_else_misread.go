package no_type_assert_else_misread

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoTypeAssertElseMisreadRule detects else branches of type assertions
// that are probably not reading the right value, since the zero value
// is used instead of the original.
type NoTypeAssertElseMisreadRule struct{}

// Apply applies the rule to given file.
func (r *NoTypeAssertElseMisreadRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	ast.Inspect(file.AST, func(n ast.Node) bool {
		checkNode(n, onFailure)
		return true
	})

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoTypeAssertElseMisreadRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	checkNode(node, onFailure)
	return failures
}

// Name returns the rule name.
func (*NoTypeAssertElseMisreadRule) Name() string {
	return "noTypeAssertElseMisread"
}

// Group returns the rule group.
func (*NoTypeAssertElseMisreadRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoTypeAssertElseMisreadRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

func checkNode(n ast.Node, onFailure func(lint.Failure)) {
	// Check if-statements with init type assertion: if v, ok := x.(T); ok { ... } else { uses v }
	if ifStmt, ok := n.(*ast.IfStmt); ok {
		if ifStmt.Init != nil {
			valueIdent, okIdent := extractTypeAssertAssign(ifStmt.Init)
			if valueIdent != nil && okIdent != nil {
				checkIfWithTypeAssert(ifStmt, valueIdent, okIdent, onFailure)
			}
		}
	}

	// Check blocks for pattern: v, ok := x.(T); followed by if ok { ... } else { uses v }
	block, ok := n.(*ast.BlockStmt)
	if !ok {
		return
	}

	for i := 0; i < len(block.List)-1; i++ {
		valueIdent, okIdent := extractTypeAssertAssign(block.List[i])
		if valueIdent == nil || okIdent == nil {
			continue
		}

		ifStmt, ok := block.List[i+1].(*ast.IfStmt)
		if !ok {
			continue
		}

		// Skip if the if statement has its own Init (could be a different assertion)
		if ifStmt.Init != nil {
			continue
		}

		checkIfWithTypeAssert(ifStmt, valueIdent, okIdent, onFailure)
	}
}

func checkIfWithTypeAssert(ifStmt *ast.IfStmt, valueIdent, okIdent *ast.Ident, onFailure func(lint.Failure)) {
	isPositive := isIdentCondition(ifStmt.Cond, okIdent.Name)
	isNegated := isNegatedIdentCondition(ifStmt.Cond, okIdent.Name)

	if !isPositive && !isNegated {
		return
	}

	// Determine which branch is the "else" (failure) branch
	var elseBranch ast.Node
	if isPositive {
		elseBranch = ifStmt.Else
	} else {
		// condition is !ok, so the "if" body is the else-equivalent
		elseBranch = ifStmt.Body
	}

	if elseBranch == nil {
		return
	}

	// Check if the value variable is used in the else branch
	refs := findIdentRefs(elseBranch, valueIdent.Name)
	for _, ref := range refs {
		onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       ref,
			Failure: fmt.Sprintf(
				"type assertion else branch reads %s which is the zero value of the asserted type, not the original value",
				valueIdent.Name,
			),
		})
	}
}

// extractTypeAssertAssign checks if a statement is a two-value type assertion
// assignment (v, ok := x.(T)) and returns the value and ok identifiers.
func extractTypeAssertAssign(stmt ast.Stmt) (*ast.Ident, *ast.Ident) {
	assign, ok := stmt.(*ast.AssignStmt)
	if !ok {
		return nil, nil
	}
	if len(assign.Lhs) != 2 || len(assign.Rhs) != 1 {
		return nil, nil
	}

	// Check RHS is a type assertion
	_, ok = assign.Rhs[0].(*ast.TypeAssertExpr)
	if !ok {
		return nil, nil
	}

	valueIdent, ok := assign.Lhs[0].(*ast.Ident)
	if !ok {
		return nil, nil
	}

	okIdent, ok := assign.Lhs[1].(*ast.Ident)
	if !ok {
		return nil, nil
	}

	// Skip if the ok variable is blank
	if okIdent.Name == "_" {
		return nil, nil
	}

	// Skip if the value variable is blank
	if valueIdent.Name == "_" {
		return nil, nil
	}

	return valueIdent, okIdent
}

// isIdentCondition checks if expr is a simple identifier with the given name.
func isIdentCondition(expr ast.Expr, name string) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == name
}

// isNegatedIdentCondition checks if expr is !name.
func isNegatedIdentCondition(expr ast.Expr, name string) bool {
	unary, ok := expr.(*ast.UnaryExpr)
	if !ok || unary.Op != token.NOT {
		return false
	}
	return isIdentCondition(unary.X, name)
}

// findIdentRefs finds all references to an identifier with the given name
// within the given AST node. It is scope-aware: if the variable is
// redeclared via := in an inner scope, references within that scope are
// skipped because they refer to the new variable, not the outer zero-valued one.
func findIdentRefs(node ast.Node, name string) []*ast.Ident {
	var refs []*ast.Ident
	findIdentRefsInner(node, name, &refs)
	return refs
}

func findIdentRefsInner(node ast.Node, name string, refs *[]*ast.Ident) {
	// For block statements, walk statements sequentially and stop once
	// the variable is redeclared, since subsequent statements in the
	// same block refer to the new variable.
	if block, ok := node.(*ast.BlockStmt); ok {
		for _, stmt := range block.List {
			if stmtDefinesName(stmt, name) {
				return // stop: variable redeclared in this scope
			}
			findIdentRefsInner(stmt, name, refs)
		}
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if n == node {
			return true // always descend into the root
		}
		switch n := n.(type) {
		case *ast.BlockStmt:
			// Handle blocks with our scope-aware logic instead.
			findIdentRefsInner(n, name, refs)
			return false
		case *ast.IfStmt:
			// If the IfStmt's Init redeclares the variable, skip the whole IfStmt.
			if n.Init != nil {
				if assign, ok := n.Init.(*ast.AssignStmt); ok &&
					assign.Tok == token.DEFINE && assignDefinesName(assign, name) {
					return false
				}
			}
		case *ast.Ident:
			if n.Name == name {
				*refs = append(*refs, n)
			}
		}
		return true
	})
}

// stmtDefinesName reports whether the statement defines the given name via :=.
// This covers both direct assignments and if-statements with init clauses.
func stmtDefinesName(stmt ast.Stmt, name string) bool {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		return s.Tok == token.DEFINE && assignDefinesName(s, name)
	case *ast.IfStmt:
		if s.Init != nil {
			if assign, ok := s.Init.(*ast.AssignStmt); ok {
				return assign.Tok == token.DEFINE && assignDefinesName(assign, name)
			}
		}
	}
	return false
}

// assignDefinesName reports whether the assignment defines a variable with the given name on its LHS.
func assignDefinesName(assign *ast.AssignStmt, name string) bool {
	for _, lhs := range assign.Lhs {
		if id, ok := lhs.(*ast.Ident); ok && id.Name == name {
			return true
		}
	}
	return false
}
