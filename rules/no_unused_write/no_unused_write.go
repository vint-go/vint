package no_unused_write

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoUnusedWriteRule detects writes to elements of a struct or array
// where the struct or array is never used after the write. This often
// happens when a function copies a struct by value, modifies a field
// of the copy, and then never uses the copy again.
type NoUnusedWriteRule struct{}

// Apply applies the rule to given file.
func (r *NoUnusedWriteRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}
		failures = append(failures, checkFunction(funcDecl, file)...)
	}

	return failures
}

// Name returns the rule name.
func (*NoUnusedWriteRule) Name() string {
	return "noUnusedWrite"
}

// Group returns the rule group.
func (*NoUnusedWriteRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnusedWriteRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoUnusedWriteRule) RequiresTypecheck() bool {
	return true
}

// checkFunction analyzes a single function declaration for unused writes.
func checkFunction(funcDecl *ast.FuncDecl, file *lint.File) []lint.Failure {
	var failures []lint.Failure

	stmts := funcDecl.Body.List
	for i, stmt := range stmts {
		assignStmt, ok := stmt.(*ast.AssignStmt)
		if !ok {
			continue
		}
		// Only check plain assignments (=), not := or +=, etc.
		if assignStmt.Tok != token.ASSIGN {
			continue
		}

		for _, lhs := range assignStmt.Lhs {
			varName, fieldOrIndex := extractWriteTarget(lhs)
			if varName == "" || fieldOrIndex == "" {
				continue
			}

			// Check if the variable is a value type (struct or array)
			if !isValueType(varName, funcDecl, file) {
				continue
			}

			// Check if the variable is used after this assignment
			remainingStmts := stmts[i+1:]
			if !isVariableUsedAfter(varName, remainingStmts) {
				failures = append(failures, lint.Failure{
					Category:   lint.FailureCategoryLogic,
					Confidence: 1,
					Node:       assignStmt,
					Failure:    fmt.Sprintf("unused write to field or element of %s", varName),
				})
			}
		}
	}

	return failures
}

// extractWriteTarget extracts the root variable name and the field/index
// being written to from a selector or index expression.
// For example, for q.X it returns ("q", "X").
// For arr[0] it returns ("arr", "[0]").
func extractWriteTarget(expr ast.Expr) (varName string, fieldOrIndex string) {
	switch e := expr.(type) {
	case *ast.SelectorExpr:
		// q.X = ...
		if ident, ok := e.X.(*ast.Ident); ok {
			return ident.Name, e.Sel.Name
		}
	case *ast.IndexExpr:
		// arr[i] = ...
		if ident, ok := e.X.(*ast.Ident); ok {
			return ident.Name, "[index]"
		}
	}
	return "", ""
}

// isValueType checks if the named variable is a value type (struct or array),
// as opposed to a pointer, slice, map, or other reference type.
func isValueType(varName string, funcDecl *ast.FuncDecl, file *lint.File) bool {
	typesInfo := file.Pkg.TypesInfo()
	if typesInfo == nil {
		return false
	}

	// Search for the variable definition in the function body
	var varType types.Type
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		if varType != nil {
			return false
		}
		switch stmt := n.(type) {
		case *ast.AssignStmt:
			if stmt.Tok == token.DEFINE {
				for _, lhs := range stmt.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok && ident.Name == varName {
						if obj := typesInfo.ObjectOf(ident); obj != nil {
							varType = obj.Type()
						}
						return false
					}
				}
			}
		case *ast.DeclStmt:
			if genDecl, ok := stmt.Decl.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
				for _, spec := range genDecl.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						for _, name := range valueSpec.Names {
							if name.Name == varName {
								if obj := typesInfo.ObjectOf(name); obj != nil {
									varType = obj.Type()
								}
								return false
							}
						}
					}
				}
			}
		}
		return true
	})

	// Also check function parameters
	if varType == nil && funcDecl.Type.Params != nil {
		for _, field := range funcDecl.Type.Params.List {
			for _, name := range field.Names {
				if name.Name == varName {
					if obj := typesInfo.ObjectOf(name); obj != nil {
						varType = obj.Type()
					}
				}
			}
		}
	}

	if varType == nil {
		return false
	}

	// Check the underlying type - we only flag value types (structs and arrays)
	underlying := varType.Underlying()
	switch underlying.(type) {
	case *types.Struct:
		return true
	case *types.Array:
		return true
	default:
		return false
	}
}

// isVariableUsedAfter checks if the named variable appears in any of the
// subsequent statements in any capacity (read, passed as argument, returned, etc.).
func isVariableUsedAfter(varName string, stmts []ast.Stmt) bool {
	for _, stmt := range stmts {
		if containsVarUse(varName, stmt) {
			return true
		}
	}
	return false
}

// containsVarUse checks if the named variable is referenced in the given AST node.
func containsVarUse(varName string, node ast.Node) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if found {
			return false
		}
		if ident, ok := n.(*ast.Ident); ok && ident.Name == varName {
			found = true
			return false
		}
		return true
	})
	return found
}
