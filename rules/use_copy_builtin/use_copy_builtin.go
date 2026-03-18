package use_copy_builtin

import (
	"fmt"
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseCopyBuiltinRule detects for loops that copy elements from one slice to
// another and suggests using the built-in copy function instead.
type UseCopyBuiltinRule struct{}

// Apply applies the rule to given file.
func (r *UseCopyBuiltinRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintCopyBuiltin{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseCopyBuiltinRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintCopyBuiltin{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseCopyBuiltinRule) Name() string {
	return "useCopyBuiltin"
}

// Group returns the rule group.
func (*UseCopyBuiltinRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseCopyBuiltinRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintCopyBuiltin struct {
	onFailure func(lint.Failure)
}

func (w *lintCopyBuiltin) Visit(node ast.Node) ast.Visitor {
	rangeStmt, ok := node.(*ast.RangeStmt)
	if !ok {
		return w
	}

	// The body must contain exactly one statement: an assignment.
	if len(rangeStmt.Body.List) != 1 {
		return w
	}

	assignStmt, ok := rangeStmt.Body.List[0].(*ast.AssignStmt)
	if !ok {
		return w
	}

	// Must be a simple assignment with exactly one LHS and one RHS.
	if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		return w
	}

	// The assignment operator must be "=" (not ":=" or "+=", etc.)
	if assignStmt.Tok.String() != "=" {
		return w
	}

	// LHS must be an index expression: dst[i]
	lhsIndex, ok := assignStmt.Lhs[0].(*ast.IndexExpr)
	if !ok {
		return w
	}

	// Get the range source expression as a string.
	rangeXStr := astutils.GoFmt(rangeStmt.X)
	if rangeXStr == "" {
		return w
	}

	// Get the LHS destination slice name.
	lhsDstStr := astutils.GoFmt(lhsIndex.X)
	if lhsDstStr == "" {
		return w
	}

	// Check: for i, v := range src { dst[i] = v }
	// In this pattern, the range has both Key (i) and Value (v).
	if rangeStmt.Key != nil && rangeStmt.Value != nil {
		keyIdent, keyOk := rangeStmt.Key.(*ast.Ident)
		valIdent, valOk := rangeStmt.Value.(*ast.Ident)
		if keyOk && valOk && keyIdent.Name != "_" && valIdent.Name != "_" {
			// LHS index must be the key variable: dst[i]
			lhsIdxStr := astutils.GoFmt(lhsIndex.Index)
			if lhsIdxStr == keyIdent.Name {
				// RHS must be the value variable: v
				rhsIdent, rhsOk := assignStmt.Rhs[0].(*ast.Ident)
				if rhsOk && rhsIdent.Name == valIdent.Name {
					w.onFailure(lint.Failure{
						Confidence: 1,
						Node:       rangeStmt,
						Category:   lint.FailureCategoryStyle,
						Failure:    fmt.Sprintf("should replace loop with copy(%s, %s)", lhsDstStr, rangeXStr),
					})
					return w
				}
			}
		}
	}

	// Check: for i := range src { dst[i] = src[i] }
	// In this pattern, the range has only Key (i) and Value is nil or blank.
	if rangeStmt.Key != nil {
		keyIdent, keyOk := rangeStmt.Key.(*ast.Ident)
		if keyOk && keyIdent.Name != "_" {
			// Value must be nil or blank.
			valueIsBlankOrNil := rangeStmt.Value == nil
			if !valueIsBlankOrNil {
				if valIdent, ok := rangeStmt.Value.(*ast.Ident); ok && valIdent.Name == "_" {
					valueIsBlankOrNil = true
				}
			}
			if valueIsBlankOrNil {
				// LHS index must be the key variable: dst[i]
				lhsIdxStr := astutils.GoFmt(lhsIndex.Index)
				if lhsIdxStr == keyIdent.Name {
					// RHS must be src[i] - index expression into the range source.
					rhsIndex, rhsOk := assignStmt.Rhs[0].(*ast.IndexExpr)
					if rhsOk {
						rhsSrcStr := astutils.GoFmt(rhsIndex.X)
						rhsIdxStr := astutils.GoFmt(rhsIndex.Index)
						if rhsSrcStr == rangeXStr && rhsIdxStr == keyIdent.Name {
							w.onFailure(lint.Failure{
								Confidence: 1,
								Node:       rangeStmt,
								Category:   lint.FailureCategoryStyle,
								Failure:    fmt.Sprintf("should replace loop with copy(%s, %s)", lhsDstStr, rangeXStr),
							})
							return w
						}
					}
				}
			}
		}
	}

	return w
}
