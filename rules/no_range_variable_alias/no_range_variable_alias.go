package no_range_variable_alias

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/strowk/vint/lint"
)

// NoRangeVariableAliasRule detects implicit memory aliasing in for...range statements
// (applicable to Go 1.21 or lower).
type NoRangeVariableAliasRule struct{}

// Apply applies the rule to given file.
func (r *NoRangeVariableAliasRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if file.Pkg.IsAtLeastGoVersion(lint.Go122) {
		return failures
	}

	file.Pkg.TypeCheck()

	w := &lintNoRangeVariableAlias{
		file: file,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoRangeVariableAliasRule) Name() string {
	return "noRangeVariableAlias"
}

// Group returns the rule group.
func (*NoRangeVariableAliasRule) Group() string {
	return "correctness"
}

func (*NoRangeVariableAliasRule) RequiresTypecheck() bool {
	return true
}

type lintNoRangeVariableAlias struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintNoRangeVariableAlias) Visit(node ast.Node) ast.Visitor {
	rangeStmt, ok := node.(*ast.RangeStmt)
	if !ok {
		return w
	}

	// Collect the range variables (key and value)
	var rangeVars []*ast.Ident
	if key, ok := rangeStmt.Key.(*ast.Ident); ok && key.Name != "_" {
		rangeVars = append(rangeVars, key)
	}
	if value, ok := rangeStmt.Value.(*ast.Ident); ok && value.Name != "_" {
		rangeVars = append(rangeVars, value)
	}

	if len(rangeVars) == 0 {
		return w
	}

	// Check if the range variable type is already a pointer (in which case &v is fine)
	pointerVars := map[*ast.Object]bool{}
	for _, v := range rangeVars {
		if t := w.file.Pkg.TypeOf(v); t != nil {
			if strings.HasPrefix(t.String(), "*") {
				pointerVars[v.Obj] = true
			}
		}
	}

	// Walk the body looking for address-of operations on range variables
	bodyWalker := &noRangeVarAliasBodyWalker{
		rangeVars:   rangeVars,
		pointerVars: pointerVars,
		onFailure:   w.onFailure,
	}
	ast.Walk(bodyWalker, rangeStmt.Body)

	return w
}

type noRangeVarAliasBodyWalker struct {
	rangeVars   []*ast.Ident
	pointerVars map[*ast.Object]bool
	onFailure   func(lint.Failure)
}

func (bw *noRangeVarAliasBodyWalker) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.UnaryExpr:
		if n.Op == token.AND {
			if bw.isRangeVar(n.X) {
				bw.onFailure(lint.Failure{
					Category:   lint.FailureCategoryLogic,
					Confidence: 1,
					Node:       n,
					Failure:    fmt.Sprintf("implicit memory aliasing in for...range: taking address of range variable '%s'", bw.identName(n.X)),
				})
			}
		}
	case *ast.RangeStmt:
		// Don't descend into nested range statements - they have their own variables
		return nil
	}
	return bw
}

// isRangeVar checks if the expression refers to a range variable (or a field of one).
func (bw *noRangeVarAliasBodyWalker) isRangeVar(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		for _, v := range bw.rangeVars {
			if e.Obj != nil && e.Obj == v.Obj && !bw.pointerVars[e.Obj] {
				return true
			}
		}
	case *ast.SelectorExpr:
		// For &v.Field - check if v is a range variable
		if ident, ok := e.X.(*ast.Ident); ok {
			for _, v := range bw.rangeVars {
				if ident.Obj != nil && ident.Obj == v.Obj && !bw.pointerVars[ident.Obj] {
					return true
				}
			}
		}
	}
	return false
}

// identName extracts the identifier name from an expression.
func (bw *noRangeVarAliasBodyWalker) identName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		if ident, ok := e.X.(*ast.Ident); ok {
			return ident.Name + "." + e.Sel.Name
		}
	}
	return ""
}
