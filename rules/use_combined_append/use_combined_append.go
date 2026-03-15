package use_combined_append

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseCombinedAppendRule detects consecutive append operations on the same slice
// that could be combined into a single call for better efficiency.
type UseCombinedAppendRule struct{}

// Apply applies the rule to given file.
func (r *UseCombinedAppendRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintCombinedAppend{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseCombinedAppendRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintCombinedAppend{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseCombinedAppendRule) Name() string {
	return "useCombinedAppend"
}

// Group returns the rule group.
func (*UseCombinedAppendRule) Group() string {
	return "performance"
}

// CacheTier returns the cache tier for this rule.
func (*UseCombinedAppendRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintCombinedAppend struct {
	onFailure func(lint.Failure)
}

// appendInfo holds information about an append assignment statement.
type appendInfo struct {
	sliceName string
	stmt      *ast.AssignStmt
}

// extractAppendInfo checks if a statement is of the form `x = append(x, ...)` and
// returns the slice name and statement if so. It returns nil if the statement does
// not match or uses the ellipsis form (e.g. `x = append(x, y...)`).
func extractAppendInfo(stmt ast.Stmt) *appendInfo {
	assign, ok := stmt.(*ast.AssignStmt)
	if !ok {
		return nil
	}

	// Only single assignments: x = append(x, ...)
	if len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
		return nil
	}

	call, ok := assign.Rhs[0].(*ast.CallExpr)
	if !ok {
		return nil
	}

	ident, ok := call.Fun.(*ast.Ident)
	if !ok || ident.Name != "append" {
		return nil
	}

	// Need at least the slice argument plus one element
	if len(call.Args) < 2 {
		return nil
	}

	// Skip variadic form: x = append(x, y...)
	if call.Ellipsis.IsValid() {
		return nil
	}

	lhsStr := astutils.GoFmt(assign.Lhs[0])
	firstArgStr := astutils.GoFmt(call.Args[0])

	// The LHS must match the first argument (self-append pattern)
	if lhsStr == "" || lhsStr != firstArgStr {
		return nil
	}

	return &appendInfo{
		sliceName: lhsStr,
		stmt:      assign,
	}
}

func (w *lintCombinedAppend) Visit(node ast.Node) ast.Visitor {
	block, ok := node.(*ast.BlockStmt)
	if !ok {
		return w
	}

	w.checkStatements(block.List)

	return w
}

func (w *lintCombinedAppend) checkStatements(stmts []ast.Stmt) {
	if len(stmts) < 2 {
		return
	}

	i := 0
	for i < len(stmts) {
		info := extractAppendInfo(stmts[i])
		if info == nil {
			i++
			continue
		}

		// Found an append; check if the next statement is also an append to the same slice
		j := i + 1
		for j < len(stmts) {
			nextInfo := extractAppendInfo(stmts[j])
			if nextInfo == nil || nextInfo.sliceName != info.sliceName {
				break
			}
			j++
		}

		if j > i+1 {
			// We have consecutive appends from index i to j-1.
			// Report the failure on the second append statement (the first one that could be combined).
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       stmts[i+1].(*ast.AssignStmt),
				Category:   lint.FailureCategoryOptimization,
				Failure:    "consecutive append to " + info.sliceName + " can be combined into a single call",
			})
		}

		i = j
	}
}
