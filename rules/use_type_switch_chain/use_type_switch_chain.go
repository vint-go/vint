package use_type_switch_chain

import (
	"fmt"
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseTypeSwitchChainRule detects repeated type assertions in if-else chains
// and suggests replacing them with type switch statements.
type UseTypeSwitchChainRule struct{}

// Apply applies the rule to given file.
func (r *UseTypeSwitchChainRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintTypeSwitchChain{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseTypeSwitchChainRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintTypeSwitchChain{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*UseTypeSwitchChainRule) Name() string {
	return "useTypeSwitchChain"
}

// Group returns the rule group.
func (*UseTypeSwitchChainRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseTypeSwitchChainRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintTypeSwitchChain struct {
	onFailure func(lint.Failure)
}

func (w *lintTypeSwitchChain) Visit(node ast.Node) ast.Visitor {
	ifStmt, ok := node.(*ast.IfStmt)
	if !ok {
		return w
	}

	// Check if this if-else chain contains type assertions on the same variable
	varName, count := countTypeAssertionChain(ifStmt)
	if count >= 2 {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Category:   lint.FailureCategoryStyle,
			Failure:    fmt.Sprintf("could replace type assertion if-else chain on %s with type switch", varName),
			Node:       ifStmt,
		})
		// Don't recurse into the else branches to avoid double-reporting
		return nil
	}

	return w
}

// countTypeAssertionChain counts the number of consecutive type assertions
// on the same variable in an if-else chain. Returns the variable name and count.
func countTypeAssertionChain(ifStmt *ast.IfStmt) (string, int) {
	// Check if the first if has a type assertion in its init
	varName := extractTypeAssertionVar(ifStmt)
	if varName == "" {
		return "", 0
	}

	count := 1
	current := ifStmt

	for {
		elseStmt := current.Else
		if elseStmt == nil {
			break
		}

		elseIf, ok := elseStmt.(*ast.IfStmt)
		if !ok {
			// Plain else block, stop counting but still include what we have
			break
		}

		elseVarName := extractTypeAssertionVar(elseIf)
		if elseVarName != varName {
			// Not a type assertion or different variable
			break
		}

		count++
		current = elseIf
	}

	return varName, count
}

// extractTypeAssertionVar extracts the variable being type-asserted in an if statement
// of the form: if v, ok := x.(Type); ok { ... }
// Returns the string representation of x, or empty string if the pattern doesn't match.
func extractTypeAssertionVar(ifStmt *ast.IfStmt) string {
	// Must have an init statement
	if ifStmt.Init == nil {
		return ""
	}

	// Init must be an assignment: v, ok := x.(Type)
	assign, ok := ifStmt.Init.(*ast.AssignStmt)
	if !ok {
		return ""
	}

	// Must be := (short variable declaration)
	if assign.Tok.String() != ":=" {
		return ""
	}

	// Must have 2 LHS variables (v, ok)
	if len(assign.Lhs) != 2 {
		return ""
	}

	// Must have 1 RHS expression
	if len(assign.Rhs) != 1 {
		return ""
	}

	// RHS must be a type assertion
	typeAssert, ok := assign.Rhs[0].(*ast.TypeAssertExpr)
	if !ok {
		return ""
	}

	// Must not be a type switch (x.(type))
	if typeAssert.Type == nil {
		return ""
	}

	// The condition must check the ok variable
	condIdent, ok := ifStmt.Cond.(*ast.Ident)
	if !ok {
		return ""
	}

	// The ok variable should be the second LHS variable
	okIdent, ok := assign.Lhs[1].(*ast.Ident)
	if !ok {
		return ""
	}

	if condIdent.Name != okIdent.Name {
		return ""
	}

	// Return the string representation of the expression being asserted
	return astutils.GoFmt(typeAssert.X)
}
