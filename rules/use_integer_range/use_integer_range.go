package use_integer_range

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseIntegerRangeRule detects C-style for loops that can be replaced with Go 1.22+ integer range syntax.
type UseIntegerRangeRule struct{}

// Apply applies the rule to given file.
func (r *UseIntegerRangeRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintIntRange{file: file, onFailure: onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseIntegerRangeRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintIntRange{file: file, onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseIntegerRangeRule) Name() string {
	return "useIntegerRange"
}

// Group returns the rule group.
func (*UseIntegerRangeRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseIntegerRangeRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintIntRange struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintIntRange) Visit(node ast.Node) ast.Visitor {
	forStmt, ok := node.(*ast.ForStmt)
	if !ok {
		return w
	}

	// Check init: must be `i := 0`
	varName, ok := isZeroInit(forStmt.Init)
	if !ok {
		return w
	}

	// Check condition: must be `i < n` or `i <= n-1` etc.
	upperBound, ok := isLessThanCond(forStmt.Cond, varName)
	if !ok {
		return w
	}

	// Check post: must be i++, i += 1, i = i + 1, or i = 1 + i
	if !isIncrementByOne(forStmt.Post, varName) {
		return w
	}

	// Check that the loop variable is not modified in the body
	if isVarModifiedInBody(forStmt.Body, varName) {
		return w
	}

	boundStr := w.file.Render(upperBound)

	// Check if loop variable is used in the body
	if !isVarUsedInBody(forStmt.Body, varName) {
		w.onFailure(lint.Failure{
			Failure:    fmt.Sprintf("for loop can be simplified to `for range %s`", boundStr),
			Confidence: 1,
			Node:       forStmt,
			Category:   lint.FailureCategoryStyle,
		})
	} else {
		w.onFailure(lint.Failure{
			Failure:    fmt.Sprintf("for loop can be simplified to `for %s := range %s`", varName, boundStr),
			Confidence: 1,
			Node:       forStmt,
			Category:   lint.FailureCategoryStyle,
		})
	}

	return w
}

// isZeroInit checks if the init statement is `varName := 0` and returns the variable name.
func isZeroInit(stmt ast.Stmt) (string, bool) {
	assign, ok := stmt.(*ast.AssignStmt)
	if !ok || assign.Tok != token.DEFINE || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
		return "", false
	}

	ident, ok := assign.Lhs[0].(*ast.Ident)
	if !ok {
		return "", false
	}

	lit, ok := assign.Rhs[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.INT || lit.Value != "0" {
		return "", false
	}

	return ident.Name, true
}

// isLessThanCond checks if the condition is `varName < expr` and returns the upper bound expression.
func isLessThanCond(expr ast.Expr, varName string) (ast.Expr, bool) {
	bin, ok := expr.(*ast.BinaryExpr)
	if !ok {
		return nil, false
	}

	switch bin.Op {
	case token.LSS: // i < n
		if !astutils.IsIdent(bin.X, varName) {
			return nil, false
		}
		return bin.Y, true
	case token.GTR: // n > i
		if !astutils.IsIdent(bin.Y, varName) {
			return nil, false
		}
		return bin.X, true
	default:
		return nil, false
	}
}

// isIncrementByOne checks if the post statement is i++, i += 1, i = i + 1, or i = 1 + i.
func isIncrementByOne(stmt ast.Stmt, varName string) bool {
	switch s := stmt.(type) {
	case *ast.IncDecStmt:
		// i++
		return s.Tok == token.INC && astutils.IsIdent(s.X, varName)
	case *ast.AssignStmt:
		if len(s.Lhs) != 1 || !astutils.IsIdent(s.Lhs[0], varName) {
			return false
		}
		switch s.Tok {
		case token.ADD_ASSIGN: // i += 1
			if len(s.Rhs) != 1 {
				return false
			}
			return isIntLit(s.Rhs[0], "1")
		case token.ASSIGN: // i = i + 1 or i = 1 + i
			if len(s.Rhs) != 1 {
				return false
			}
			bin, ok := s.Rhs[0].(*ast.BinaryExpr)
			if !ok || bin.Op != token.ADD {
				return false
			}
			return (astutils.IsIdent(bin.X, varName) && isIntLit(bin.Y, "1")) ||
				(isIntLit(bin.X, "1") && astutils.IsIdent(bin.Y, varName))
		}
	}
	return false
}

func isIntLit(expr ast.Expr, value string) bool {
	lit, ok := expr.(*ast.BasicLit)
	return ok && lit.Kind == token.INT && lit.Value == value
}

// isVarModifiedInBody checks if the variable is assigned or incremented in the loop body.
func isVarModifiedInBody(body *ast.BlockStmt, varName string) bool {
	modified := false
	ast.Inspect(body, func(n ast.Node) bool {
		if modified {
			return false
		}
		switch s := n.(type) {
		case *ast.AssignStmt:
			for _, lhs := range s.Lhs {
				if astutils.IsIdent(lhs, varName) {
					modified = true
					return false
				}
			}
		case *ast.IncDecStmt:
			if astutils.IsIdent(s.X, varName) {
				modified = true
				return false
			}
		}
		return true
	})
	return modified
}

// isVarUsedInBody checks if the variable is referenced in the loop body.
func isVarUsedInBody(body *ast.BlockStmt, varName string) bool {
	used := false
	ast.Inspect(body, func(n ast.Node) bool {
		if used {
			return false
		}
		ident, ok := n.(*ast.Ident)
		if ok && ident.Name == varName {
			used = true
			return false
		}
		return true
	})
	return used
}
