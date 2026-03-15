package use_assignment_operator

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseAssignmentOperatorRule detects assignments that can be simplified
// by using assignment operators (+=, -=, *=, etc.).
type UseAssignmentOperatorRule struct{}

// Apply applies the rule to given file.
func (r *UseAssignmentOperatorRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintUseAssignOp{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseAssignmentOperatorRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintUseAssignOp{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseAssignmentOperatorRule) Name() string {
	return "useAssignmentOperator"
}

// Group returns the rule group.
func (*UseAssignmentOperatorRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseAssignmentOperatorRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// binOpToAssignOp maps binary operators to their corresponding assignment operators.
var binOpToAssignOp = map[token.Token]string{
	token.ADD:     "+=",
	token.SUB:     "-=",
	token.MUL:     "*=",
	token.QUO:     "/=",
	token.REM:     "%=",
	token.AND:     "&=",
	token.OR:      "|=",
	token.XOR:     "^=",
	token.SHL:     "<<=",
	token.SHR:     ">>=",
	token.AND_NOT: "&^=",
}

type lintUseAssignOp struct {
	onFailure func(lint.Failure)
}

func (w *lintUseAssignOp) Visit(node ast.Node) ast.Visitor {
	assign, ok := node.(*ast.AssignStmt)
	if !ok {
		return w
	}

	// Only check plain assignments (=), not :=, +=, etc.
	if assign.Tok != token.ASSIGN {
		return w
	}

	// Must be a single assignment (not multiple)
	if len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
		return w
	}

	// RHS must be a binary expression
	binExpr, ok := assign.Rhs[0].(*ast.BinaryExpr)
	if !ok {
		return w
	}

	// The binary operator must have a corresponding assignment operator
	assignOp, ok := binOpToAssignOp[binExpr.Op]
	if !ok {
		return w
	}

	// The LHS of the assignment must match the left operand of the binary expression
	lhsStr := astutils.GoFmt(assign.Lhs[0])
	binLhsStr := astutils.GoFmt(binExpr.X)

	if lhsStr == "" || lhsStr != binLhsStr {
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       assign,
		Failure:    fmt.Sprintf("replace %s = %s %s %s with %s %s %s", lhsStr, binLhsStr, binExpr.Op.String(), astutils.GoFmt(binExpr.Y), lhsStr, assignOp, astutils.GoFmt(binExpr.Y)),
	})

	return w
}
