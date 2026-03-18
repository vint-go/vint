package no_redundant_var_type

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoRedundantVarTypeRule detects variable declarations where the type is
// explicitly specified but is redundant because it matches the type that
// would be inferred from the right-hand side expression.
type NoRedundantVarTypeRule struct{}

// Apply applies the rule to given file.
func (r *NoRedundantVarTypeRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintRedundantVarType{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
		pkg: file.Pkg,
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoRedundantVarTypeRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintRedundantVarType{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
		pkg: file.Pkg,
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoRedundantVarTypeRule) Name() string {
	return "noRedundantVarType"
}

// Group returns the rule group.
func (*NoRedundantVarTypeRule) Group() string {
	return "style"
}

// RequiresTypecheck returns true because this rule uses type information.
func (*NoRedundantVarTypeRule) RequiresTypecheck() bool {
	return true
}

// CacheTier returns the cache tier for this rule.
func (*NoRedundantVarTypeRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

type lintRedundantVarType struct {
	onFailure func(lint.Failure)
	pkg       *lint.Package
}

func (w *lintRedundantVarType) Visit(node ast.Node) ast.Visitor {
	genDecl, ok := node.(*ast.GenDecl)
	if !ok || genDecl.Tok != token.VAR {
		return w
	}

	typesInfo := w.pkg.TypesInfo()
	if typesInfo == nil {
		return w
	}

	for _, spec := range genDecl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		// Must have an explicit type and at least one value
		if valueSpec.Type == nil || len(valueSpec.Values) == 0 {
			continue
		}

		// Must have matching number of names and values for simple checking
		if len(valueSpec.Names) != len(valueSpec.Values) {
			continue
		}

		// Get the declared type
		declaredType := w.pkg.TypeOf(valueSpec.Type)
		if declaredType == nil {
			continue
		}

		// Check each value's type against the declared type
		allMatch := true
		for _, val := range valueSpec.Values {
			// If the RHS involves an untyped expression (literal, untyped constant
			// identifier, untyped binary expression), the type annotation gives
			// the expression its type. Skip unless we can determine that the
			// natural default type matches the declared type.
			if isUntypedExpr(val, typesInfo) {
				// The type annotation may be meaningful. The expression's default
				// type might differ from the declared type.
				naturalType := defaultTypeOfLiteral(val)
				if naturalType == nil || !types.Identical(declaredType, naturalType) {
					allMatch = false
					break
				}
				continue
			}

			rhsType := w.pkg.TypeOf(val)
			if rhsType == nil {
				allMatch = false
				break
			}

			if !types.Identical(declaredType, rhsType) {
				allMatch = false
				break
			}
		}

		if !allMatch {
			continue
		}

		typeName := astutils.GoFmt(valueSpec.Type)
		w.onFailure(lint.Failure{
			Confidence: 1,
			Category:   lint.FailureCategoryStyle,
			Node:       valueSpec,
			Failure:    fmt.Sprintf("redundant type %s in variable declaration", typeName),
		})
	}

	return w
}

// isUntypedExpr checks whether an expression is an untyped constant expression.
// This includes basic literals, identifiers referencing untyped constants,
// and expressions composed of untyped operands.
func isUntypedExpr(expr ast.Expr, info *types.Info) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		// All basic literals (int, float, string, char, imaginary) are untyped.
		return true
	case *ast.Ident:
		// Check if this is an untyped constant (like true, false, iota, or
		// user-defined untyped constants).
		if e.Name == "true" || e.Name == "false" || e.Name == "nil" {
			return true
		}
		tv, ok := info.Types[e]
		if ok && tv.Value != nil {
			if basic, ok := tv.Type.Underlying().(*types.Basic); ok {
				// The type checker contextually types constants, so we need
				// to check if the original definition was untyped.
				// Check if the object is a builtin or untyped constant.
				obj := info.ObjectOf(e)
				if obj != nil {
					if c, ok := obj.(*types.Const); ok {
						if basic, ok := c.Type().(*types.Basic); ok {
							return basic.Info()&types.IsUntyped != 0
						}
					}
				}
				_ = basic
			}
		}
		return false
	case *ast.UnaryExpr:
		return isUntypedExpr(e.X, info)
	case *ast.BinaryExpr:
		return isUntypedExpr(e.X, info) && isUntypedExpr(e.Y, info)
	case *ast.ParenExpr:
		return isUntypedExpr(e.X, info)
	default:
		return false
	}
}

// defaultTypeOfLiteral returns the default Go type for a basic literal
// or simple untyped constant expression based on its AST form.
func defaultTypeOfLiteral(expr ast.Expr) types.Type {
	switch e := expr.(type) {
	case *ast.BasicLit:
		switch e.Kind {
		case token.INT:
			return types.Typ[types.Int]
		case token.FLOAT:
			return types.Typ[types.Float64]
		case token.IMAG:
			return types.Typ[types.Complex128]
		case token.CHAR:
			return types.Typ[types.Rune]
		case token.STRING:
			return types.Typ[types.String]
		}
	case *ast.Ident:
		if e.Name == "true" || e.Name == "false" {
			return types.Typ[types.Bool]
		}
		if e.Name == "nil" {
			return nil // nil has no default type
		}
		// For other identifiers referencing untyped constants, we can't easily
		// determine the default type from the AST alone.
		return nil
	case *ast.UnaryExpr:
		return defaultTypeOfLiteral(e.X)
	case *ast.ParenExpr:
		return defaultTypeOfLiteral(e.X)
	}
	return nil
}
