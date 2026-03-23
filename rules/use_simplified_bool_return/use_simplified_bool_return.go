package use_simplified_bool_return

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseSimplifiedBoolReturnRule detects if/else blocks that return
// true/false and can be simplified to returning the condition directly.
type UseSimplifiedBoolReturnRule struct{}

// Apply applies the rule to given file.
func (r *UseSimplifiedBoolReturnRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintSimplifiedBoolReturn{onFailure: onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseSimplifiedBoolReturnRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintSimplifiedBoolReturn{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseSimplifiedBoolReturnRule) Name() string {
	return "useSimplifiedBoolReturn"
}

// Group returns the rule group.
func (*UseSimplifiedBoolReturnRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseSimplifiedBoolReturnRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintSimplifiedBoolReturn struct {
	onFailure func(lint.Failure)
}

func (w *lintSimplifiedBoolReturn) Visit(node ast.Node) ast.Visitor {
	block, ok := node.(*ast.BlockStmt)
	if !ok {
		return w
	}

	for i := 0; i < len(block.List)-1; i++ {
		ifStmt, ok := block.List[i].(*ast.IfStmt)
		if !ok {
			continue
		}

		// Skip if statements with init or else clauses for this pattern
		if ifStmt.Init != nil || ifStmt.Else != nil {
			continue
		}

		// The if body must be a single return statement returning a bool literal
		ifRetVal, ok := singleReturnBoolLiteral(ifStmt.Body)
		if !ok {
			continue
		}

		// The next statement must be a return statement returning a bool literal
		nextRet, ok := block.List[i+1].(*ast.ReturnStmt)
		if !ok {
			continue
		}

		nextRetVal, ok := isBoolLiteral(nextRet)
		if !ok {
			continue
		}

		// The two return values must be opposite
		if ifRetVal == nextRetVal {
			continue
		}

		w.onFailure(lint.Failure{
			Confidence: 1,
			Category:   lint.FailureCategoryStyle,
			Node:       ifStmt,
			Failure:    "redundant if/else returning bool; simplify to return the condition directly",
		})
	}

	return w
}

// singleReturnBoolLiteral checks if a block contains exactly one statement
// that is a return statement with a single boolean literal result.
// Returns the boolean value and true if it matches.
func singleReturnBoolLiteral(block *ast.BlockStmt) (bool, bool) {
	if block == nil || len(block.List) != 1 {
		return false, false
	}

	ret, ok := block.List[0].(*ast.ReturnStmt)
	if !ok {
		return false, false
	}

	return isBoolLiteral(ret)
}

// isBoolLiteral checks if a return statement returns exactly one boolean literal.
// Returns the boolean value and true if it matches.
func isBoolLiteral(ret *ast.ReturnStmt) (bool, bool) {
	if len(ret.Results) != 1 {
		return false, false
	}

	ident, ok := ret.Results[0].(*ast.Ident)
	if !ok {
		return false, false
	}

	switch ident.Name {
	case "true":
		return true, true
	case "false":
		return false, true
	default:
		return false, false
	}
}
