package no_writer_buffer_modification

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoWriterBufferModificationRule detects modifications to the byte slice parameter
// in io.Writer implementations. The io.Writer contract states that Write must not
// modify the slice data, even temporarily.
type NoWriterBufferModificationRule struct{}

// Apply applies the rule to given file.
func (r *NoWriterBufferModificationRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		checkWriteMethod(funcDecl, onFailure)
	}

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoWriterBufferModificationRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return nil
	}

	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	checkWriteMethod(funcDecl, onFailure)
	return failures
}

// checkWriteMethod checks if the function declaration is a Write method with the
// io.Writer signature, and if so, scans its body for modifications to the byte slice parameter.
func checkWriteMethod(funcDecl *ast.FuncDecl, onFailure func(lint.Failure)) {
	// Must be a method (has a receiver).
	if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
		return
	}

	// Must be named "Write".
	if funcDecl.Name.Name != "Write" {
		return
	}

	// Must have the io.Writer signature: Write(p []byte) (int, error).
	if !isWriterSignature(funcDecl.Type) {
		return
	}

	// Must have a body.
	if funcDecl.Body == nil {
		return
	}

	// Get the name of the byte slice parameter.
	paramName := getByteSliceParamName(funcDecl.Type.Params)
	if paramName == "" {
		return // unnamed parameter, can't track modifications
	}

	// Walk the function body looking for modifications to the parameter.
	finder := &bufferModFinder{
		paramName: paramName,
		onFailure: onFailure,
	}
	ast.Walk(finder, funcDecl.Body)
}

// isWriterSignature checks if a function type matches the io.Writer signature:
// one parameter of type []byte, two results of types (int, error).
func isWriterSignature(ft *ast.FuncType) bool {
	// Check exactly one parameter.
	if ft.Params == nil || countFields(ft.Params) != 1 {
		return false
	}

	// Check parameter type is []byte.
	paramTypes := flattenFieldTypes(ft.Params)
	if len(paramTypes) != 1 || typeExprToString(paramTypes[0]) != "[]byte" {
		return false
	}

	// Check exactly two results.
	if ft.Results == nil || countFields(ft.Results) != 2 {
		return false
	}

	// Check result types are (int, error).
	resultTypes := flattenFieldTypes(ft.Results)
	if len(resultTypes) != 2 {
		return false
	}
	if typeExprToString(resultTypes[0]) != "int" || typeExprToString(resultTypes[1]) != "error" {
		return false
	}

	return true
}

// getByteSliceParamName returns the name of the first parameter (the []byte slice).
func getByteSliceParamName(params *ast.FieldList) string {
	if params == nil || len(params.List) == 0 {
		return ""
	}
	field := params.List[0]
	if len(field.Names) == 0 {
		return ""
	}
	return field.Names[0].Name
}

// countFields counts the total number of individual fields (expanding grouped names).
func countFields(fl *ast.FieldList) int {
	if fl == nil {
		return 0
	}
	count := 0
	for _, f := range fl.List {
		if len(f.Names) == 0 {
			count++
		} else {
			count += len(f.Names)
		}
	}
	return count
}

// flattenFieldTypes returns one ast.Expr per field entry (expanding grouped names).
func flattenFieldTypes(fl *ast.FieldList) []ast.Expr {
	if fl == nil {
		return nil
	}
	var types []ast.Expr
	for _, f := range fl.List {
		if len(f.Names) == 0 {
			types = append(types, f.Type)
		} else {
			for range f.Names {
				types = append(types, f.Type)
			}
		}
	}
	return types
}

// typeExprToString converts a simple AST type expression to its string representation.
func typeExprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		if x, ok := t.X.(*ast.Ident); ok {
			return x.Name + "." + t.Sel.Name
		}
		return ""
	case *ast.ArrayType:
		if t.Len == nil {
			return "[]" + typeExprToString(t.Elt)
		}
		return "[...]" + typeExprToString(t.Elt)
	default:
		return ""
	}
}

// bufferModFinder walks an AST to find modifications to a named byte slice parameter.
type bufferModFinder struct {
	paramName string
	onFailure func(lint.Failure)
}

func (f *bufferModFinder) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *ast.AssignStmt:
		// Check for p[i] = ... or p[i:j] = ... (index/slice assignment)
		for _, lhs := range n.Lhs {
			if f.isParamIndexExpr(lhs) {
				f.onFailure(lint.Failure{
					Confidence: 1,
					Node:       n,
					Category:   lint.FailureCategoryLogic,
					Failure:    "an io.Writer must not modify the provided buffer",
				})
			}
		}

	case *ast.IncDecStmt:
		// Check for p[i]++ or p[i]--
		if f.isParamIndexExpr(n.X) {
			f.onFailure(lint.Failure{
				Confidence: 1,
				Node:       n,
				Category:   lint.FailureCategoryLogic,
				Failure:    "an io.Writer must not modify the provided buffer",
			})
		}
	}

	return f
}

// isParamIndexExpr checks if the expression is an indexing operation on the parameter,
// such as p[i] or p[i:j].
func (f *bufferModFinder) isParamIndexExpr(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.IndexExpr:
		// p[i]
		if ident, ok := e.X.(*ast.Ident); ok && ident.Name == f.paramName {
			return true
		}
	case *ast.SliceExpr:
		// p[i:j]
		if ident, ok := e.X.(*ast.Ident); ok && ident.Name == f.paramName {
			return true
		}
	}
	return false
}

// Name returns the rule name.
func (*NoWriterBufferModificationRule) Name() string {
	return "noWriterBufferModification"
}

// Group returns the rule group.
func (*NoWriterBufferModificationRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoWriterBufferModificationRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
