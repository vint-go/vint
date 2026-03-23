package no_reflect_value_compare

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoReflectValueCompareRule checks for accidentally using == or reflect.DeepEqual
// to compare reflect.Value values.
type NoReflectValueCompareRule struct{}

// Apply applies the rule to given file.
func (r *NoReflectValueCompareRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if file.Pkg.TypeCheck() != nil {
		return nil
	}

	w := &lintNoReflectValueCompare{
		file: file,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (*NoReflectValueCompareRule) Name() string {
	return "noReflectValueCompare"
}

// Group returns the rule group.
func (*NoReflectValueCompareRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoReflectValueCompareRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoReflectValueCompareRule) RequiresTypecheck() bool {
	return true
}

type lintNoReflectValueCompare struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintNoReflectValueCompare) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		w.checkBinaryExpr(n)
	case *ast.CallExpr:
		w.checkDeepEqualCall(n)
	}
	return w
}

func (w *lintNoReflectValueCompare) checkBinaryExpr(expr *ast.BinaryExpr) {
	if expr.Op != token.EQL && expr.Op != token.NEQ {
		return
	}

	typeOfX := w.file.Pkg.TypeOf(expr.X)
	typeOfY := w.file.Pkg.TypeOf(expr.Y)

	if isReflectValue(typeOfX) && isReflectValue(typeOfY) {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       expr,
			Failure:    "avoid using == with reflect.Value, use reflect.Value.Equal or compare with reflect.Value.Interface()",
		})
	}
}

func (w *lintNoReflectValueCompare) checkDeepEqualCall(call *ast.CallExpr) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "DeepEqual" {
		return
	}

	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok || pkgIdent.Name != "reflect" {
		return
	}

	// Verify it's actually the reflect package via type info
	typesInfo := w.file.Pkg.TypesInfo()
	if typesInfo != nil {
		obj, ok := typesInfo.Uses[pkgIdent].(*types.PkgName)
		if !ok || obj.Imported().Path() != "reflect" {
			return
		}
	}

	if len(call.Args) != 2 {
		return
	}

	typeOfA := w.file.Pkg.TypeOf(call.Args[0])
	typeOfB := w.file.Pkg.TypeOf(call.Args[1])

	if isReflectValue(typeOfA) || isReflectValue(typeOfB) {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       call,
			Failure:    "avoid using reflect.DeepEqual with reflect.Value, use reflect.Value.Equal or compare with reflect.Value.Interface()",
		})
	}
}

func isReflectValue(typ types.Type) bool {
	if typ == nil {
		return false
	}
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	return obj != nil && obj.Pkg() != nil && obj.Pkg().Path() == "reflect" && obj.Name() == "Value"
}
