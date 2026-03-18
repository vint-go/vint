package use_parallel_assign_swap

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseParallelAssignSwapRule detects value swapping code that does not use
// parallel assignment. Go supports parallel assignment which makes value
// swaps cleaner and more idiomatic.
type UseParallelAssignSwapRule struct{}

// Apply applies the rule to given file.
func (r *UseParallelAssignSwapRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintParallelSwap{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseParallelAssignSwapRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	w := &lintParallelSwap{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseParallelAssignSwapRule) Name() string {
	return "useParallelAssignSwap"
}

// Group returns the rule group.
func (*UseParallelAssignSwapRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseParallelAssignSwapRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintParallelSwap struct {
	onFailure func(lint.Failure)
}

func (w *lintParallelSwap) Visit(node ast.Node) ast.Visitor {
	block, ok := node.(*ast.BlockStmt)
	if !ok {
		return w
	}

	w.checkBlock(block.List)
	return w
}

// checkBlock scans a list of statements for the pattern:
//
//	tmp := a
//	a = b
//	b = tmp
func (w *lintParallelSwap) checkBlock(stmts []ast.Stmt) {
	if len(stmts) < 3 {
		return
	}

	for i := 0; i <= len(stmts)-3; i++ {
		// Statement 1: tmp := a  (short variable declaration with exactly 1 lhs and 1 rhs)
		assign1, ok := stmts[i].(*ast.AssignStmt)
		if !ok || assign1.Tok != token.DEFINE {
			continue
		}
		if len(assign1.Lhs) != 1 || len(assign1.Rhs) != 1 {
			continue
		}
		tmpIdent, ok := assign1.Lhs[0].(*ast.Ident)
		if !ok {
			continue
		}
		aExpr := assign1.Rhs[0]
		aStr := astutils.GoFmt(aExpr)

		// Statement 2: a = b  (plain assignment)
		assign2, ok := stmts[i+1].(*ast.AssignStmt)
		if !ok || assign2.Tok != token.ASSIGN {
			continue
		}
		if len(assign2.Lhs) != 1 || len(assign2.Rhs) != 1 {
			continue
		}
		// The LHS of statement 2 must match the RHS of statement 1 (i.e., "a")
		if astutils.GoFmt(assign2.Lhs[0]) != aStr {
			continue
		}
		bExpr := assign2.Rhs[0]
		bStr := astutils.GoFmt(bExpr)

		// Statement 3: b = tmp  (plain assignment)
		assign3, ok := stmts[i+2].(*ast.AssignStmt)
		if !ok || assign3.Tok != token.ASSIGN {
			continue
		}
		if len(assign3.Lhs) != 1 || len(assign3.Rhs) != 1 {
			continue
		}
		// The LHS of statement 3 must match the RHS of statement 2 (i.e., "b")
		if astutils.GoFmt(assign3.Lhs[0]) != bStr {
			continue
		}
		// The RHS of statement 3 must be the temp variable
		rhsIdent, ok := assign3.Rhs[0].(*ast.Ident)
		if !ok || rhsIdent.Name != tmpIdent.Name {
			continue
		}

		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       assign1,
			Category:   lint.FailureCategoryStyle,
			Failure:    "use parallel assignment to swap values: " + aStr + ", " + bStr + " = " + bStr + ", " + aStr,
		})
	}
}
