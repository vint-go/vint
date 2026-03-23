package use_string_map_key

import (
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseStringMapKeyRule detects missed optimization opportunities when indexing
// maps by byte slices. The Go compiler can optimize m[string(b)] to avoid
// allocation, but not when an intermediate variable is used.
type UseStringMapKeyRule struct{}

// Apply applies the rule to given file.
func (r *UseStringMapKeyRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	file.Pkg.TypeCheck()
	typesInfo := file.Pkg.TypesInfo()
	if typesInfo == nil {
		return nil
	}

	var failures []lint.Failure
	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	// Walk the AST looking for function bodies.
	ast.Inspect(file.AST, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.FuncDecl:
			if node.Body != nil {
				checkBlock(node.Body, typesInfo, onFailure)
			}
			return false
		case *ast.FuncLit:
			if node.Body != nil {
				checkBlock(node.Body, typesInfo, onFailure)
			}
			return false
		}
		return true
	})

	return failures
}

// Name returns the rule name.
func (*UseStringMapKeyRule) Name() string {
	return "useStringMapKey"
}

// Group returns the rule group.
func (*UseStringMapKeyRule) Group() string {
	return "performance"
}

// CacheTier returns the cache tier for this rule.
func (*UseStringMapKeyRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule uses type information.
func (*UseStringMapKeyRule) RequiresTypecheck() bool {
	return true
}

// byteToStringVar records a local variable that holds the result of string([]byte).
type byteToStringVar struct {
	obj *ast.Object
}

// checkBlock scans a function body for the pattern:
//
//	s := string(byteSlice)
//	_ = m[s]  // missed optimization
func checkBlock(body *ast.BlockStmt, typesInfo *types.Info, onFailure func(lint.Failure)) {
	// Collect all variables assigned from string([]byte) conversion.
	vars := collectByteToStringVars(body, typesInfo)
	if len(vars) == 0 {
		return
	}

	// Now walk the body looking for map index expressions using those variables.
	ast.Inspect(body, func(n ast.Node) bool {
		ie, ok := n.(*ast.IndexExpr)
		if !ok {
			return true
		}

		// Check if the index is one of our tracked variables.
		ident, ok := ie.Index.(*ast.Ident)
		if !ok {
			return true
		}

		if !isTrackedVar(ident, vars) {
			return true
		}

		// Check that the expression being indexed is a map with string keys.
		mapType := typesInfo.TypeOf(ie.X)
		if mapType == nil {
			return true
		}

		mapUnderlying, ok := mapType.Underlying().(*types.Map)
		if !ok {
			return true
		}

		if basic, ok := mapUnderlying.Key().(*types.Basic); ok && basic.Kind() == types.String {
			onFailure(lint.Failure{
				Confidence: 1,
				Node:       ie,
				Category:   lint.FailureCategoryOptimization,
				Failure:    "use string(b) directly as map key instead of intermediate variable to allow compiler optimization",
			})
		}

		return true
	})
}

// collectByteToStringVars finds all variables assigned from string([]byte) conversions.
func collectByteToStringVars(body *ast.BlockStmt, typesInfo *types.Info) []byteToStringVar {
	var vars []byteToStringVar

	ast.Inspect(body, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		// Only consider := and = assignments.
		if assign.Tok.String() != ":=" && assign.Tok.String() != "=" {
			return true
		}

		for i, rhs := range assign.Rhs {
			if i >= len(assign.Lhs) {
				break
			}

			// Check if RHS is a string(expr) call.
			call, ok := rhs.(*ast.CallExpr)
			if !ok {
				continue
			}

			funIdent, ok := call.Fun.(*ast.Ident)
			if !ok {
				continue
			}

			if funIdent.Name != "string" || len(call.Args) != 1 {
				continue
			}

			// Check that the argument is of type []byte using type info.
			argType := typesInfo.TypeOf(call.Args[0])
			if argType == nil {
				continue
			}

			if !isByteSlice(argType) {
				continue
			}

			// The LHS must be an identifier.
			lhsIdent, ok := assign.Lhs[i].(*ast.Ident)
			if !ok {
				continue
			}

			if lhsIdent.Obj != nil {
				vars = append(vars, byteToStringVar{obj: lhsIdent.Obj})
			}
		}

		return true
	})

	return vars
}

// isTrackedVar checks if an identifier matches one of the tracked variables.
func isTrackedVar(ident *ast.Ident, vars []byteToStringVar) bool {
	if ident.Obj == nil {
		return false
	}
	for _, v := range vars {
		if ident.Obj == v.obj {
			return true
		}
	}
	return false
}

// isByteSlice returns true if the type is []byte.
func isByteSlice(t types.Type) bool {
	slice, ok := t.Underlying().(*types.Slice)
	if !ok {
		return false
	}
	basic, ok := slice.Elem().(*types.Basic)
	return ok && basic.Kind() == types.Byte
}
