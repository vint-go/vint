package no_impossible_nil_comparison

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoImpossibleNilComparisonRule detects return statements where a concrete
// (non-interface) typed value is returned in a position whose declared return
// type is an interface. Because the returned interface wraps a concrete type,
// comparing the result to nil at the call site will never be true, even when
// the underlying concrete value is the zero value (e.g. a nil pointer).
type NoImpossibleNilComparisonRule struct{}

// Apply applies the rule to given file.
func (r *NoImpossibleNilComparisonRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if file.Pkg.TypeCheck() != nil {
		return nil
	}

	info := file.Pkg.TypesInfo()
	if info == nil {
		return nil
	}

	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	w := &lintImpossibleNil{
		file:      file,
		info:      info,
		onFailure: onFailure,
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoImpossibleNilComparisonRule) Name() string {
	return "noImpossibleNilComparison"
}

// Group returns the rule group.
func (*NoImpossibleNilComparisonRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoImpossibleNilComparisonRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule uses type info.
func (*NoImpossibleNilComparisonRule) RequiresTypecheck() bool {
	return true
}

type lintImpossibleNil struct {
	file      *lint.File
	info      *types.Info
	onFailure func(lint.Failure)
	// stack of function result types (to handle nested func literals)
	resultStack []resultInfo
}

type resultInfo struct {
	results *types.Tuple
}

func (w *lintImpossibleNil) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		sig := w.funcDeclSignature(n)
		if sig == nil {
			return w
		}
		results := sig.Results()
		if results == nil || results.Len() == 0 {
			return w
		}
		// Check if any result type is an interface
		hasInterface := false
		for i := 0; i < results.Len(); i++ {
			if _, ok := results.At(i).Type().Underlying().(*types.Interface); ok {
				hasInterface = true
				break
			}
		}
		if !hasInterface {
			return w
		}
		w.resultStack = append(w.resultStack, resultInfo{results: results})
		w.walkBody(n.Body)
		w.resultStack = w.resultStack[:len(w.resultStack)-1]
		return nil // we already walked the body

	case *ast.FuncLit:
		sig := w.funcLitSignature(n)
		if sig == nil {
			return w
		}
		results := sig.Results()
		if results == nil || results.Len() == 0 {
			return w
		}
		hasInterface := false
		for i := 0; i < results.Len(); i++ {
			if _, ok := results.At(i).Type().Underlying().(*types.Interface); ok {
				hasInterface = true
				break
			}
		}
		if !hasInterface {
			return w
		}
		w.resultStack = append(w.resultStack, resultInfo{results: results})
		w.walkBody(n.Body)
		w.resultStack = w.resultStack[:len(w.resultStack)-1]
		return nil // we already walked the body

	case *ast.ReturnStmt:
		w.checkReturn(n)
	}

	return w
}

func (w *lintImpossibleNil) walkBody(body *ast.BlockStmt) {
	if body == nil {
		return
	}
	for _, stmt := range body.List {
		ast.Walk(w, stmt)
	}
}

func (w *lintImpossibleNil) funcDeclSignature(fn *ast.FuncDecl) *types.Signature {
	if fn.Name == nil {
		return nil
	}
	obj := w.info.ObjectOf(fn.Name)
	if obj == nil {
		return nil
	}
	sig, _ := obj.Type().(*types.Signature)
	return sig
}

func (w *lintImpossibleNil) funcLitSignature(fn *ast.FuncLit) *types.Signature {
	tv, ok := w.info.Types[fn]
	if !ok {
		return nil
	}
	sig, _ := tv.Type.(*types.Signature)
	return sig
}

func (w *lintImpossibleNil) checkReturn(ret *ast.ReturnStmt) {
	if len(w.resultStack) == 0 {
		return
	}

	current := w.resultStack[len(w.resultStack)-1]
	results := current.results
	if results == nil {
		return
	}

	// Handle bare returns (named results) - skip
	if len(ret.Results) == 0 {
		return
	}

	// Handle single expression returning multiple values (e.g., return someFunc())
	if len(ret.Results) == 1 && results.Len() > 1 {
		return
	}

	for i, expr := range ret.Results {
		if i >= results.Len() {
			break
		}

		declType := results.At(i).Type()
		// Check if declared return type is an interface
		_, isIface := declType.Underlying().(*types.Interface)
		if !isIface {
			continue
		}

		// Skip if the expression is untyped nil (explicit nil return is fine)
		if isNilIdent(expr) {
			continue
		}

		// Get the type of the returned expression
		exprType := w.file.Pkg.TypeOf(expr)
		if exprType == nil {
			continue
		}

		// Check if the expression type is a concrete (non-interface) type
		_, exprIsIface := exprType.Underlying().(*types.Interface)
		if exprIsIface {
			continue
		}

		// The expression is a concrete type being returned as an interface.
		// This means the interface value will never be nil even if the
		// concrete value is nil.
		typeName := exprType.String()
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 0.8,
			Node:       ret,
			Failure:    fmt.Sprintf("returning concrete type %s as interface will produce non-nil interface value", typeName),
		})
	}
}

// isNilIdent checks if an expression is the nil identifier.
func isNilIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "nil"
}
