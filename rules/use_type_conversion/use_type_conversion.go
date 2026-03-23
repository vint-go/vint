package use_type_conversion

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseTypeConversionRule detects field-by-field struct copying that can be
// replaced with a type conversion.
//
// When two struct types have identical field names and types, assigning
// field by field can be replaced with a type conversion, which is clearer
// and less error-prone.
//
// Source: https://staticcheck.dev/docs/checks/#S1016
type UseTypeConversionRule struct{}

// Apply applies the rule to given file.
func (r *UseTypeConversionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintTypeConversion{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseTypeConversionRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintTypeConversion{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*UseTypeConversionRule) Name() string {
	return "useTypeConversion"
}

// Group returns the rule group.
func (*UseTypeConversionRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseTypeConversionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck indicates this rule needs type information.
func (*UseTypeConversionRule) RequiresTypecheck() bool {
	return true
}

type lintTypeConversion struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintTypeConversion) Visit(node ast.Node) ast.Visitor {
	block, ok := node.(*ast.BlockStmt)
	if !ok {
		return w
	}

	w.checkBlock(block.List)
	return w
}

// fieldAssign represents a single assignment of the form dst.Field = src.Field.
type fieldAssign struct {
	dstObj   string // rendered destination variable (e.g. "c")
	srcObj   string // rendered source variable (e.g. "p")
	field    string // field name (e.g. "X")
	stmtIdx  int    // index in the block
	stmt     ast.Stmt
}

// checkBlock scans a block of statements for consecutive field-by-field
// copies that could be replaced with a type conversion.
func (w *lintTypeConversion) checkBlock(stmts []ast.Stmt) {
	typesInfo := w.pkg.TypesInfo()
	if typesInfo == nil {
		return
	}

	// Collect all field assignments in order.
	var assigns []fieldAssign
	for i, stmt := range stmts {
		fa, ok := w.parseFieldAssign(stmt)
		if !ok {
			continue
		}
		fa.stmtIdx = i
		fa.stmt = stmt
		assigns = append(assigns, fa)
	}

	if len(assigns) < 2 {
		return
	}

	// Group consecutive assignments by (dst, src) pair.
	type pairKey struct{ dst, src string }

	i := 0
	for i < len(assigns) {
		key := pairKey{assigns[i].dstObj, assigns[i].srcObj}
		start := i
		fields := map[string]bool{assigns[i].field: true}
		lastIdx := assigns[i].stmtIdx

		j := i + 1
		for j < len(assigns) {
			a := assigns[j]
			if a.dstObj != key.dst || a.srcObj != key.src {
				break
			}
			// The assignments must be consecutive in the original block
			// (no other statements in between).
			if a.stmtIdx != lastIdx+1 {
				break
			}
			if fields[a.field] {
				// Duplicate field assignment — not a clean copy pattern.
				break
			}
			fields[a.field] = true
			lastIdx = a.stmtIdx
			j++
		}

		count := j - start
		if count >= 2 {
			w.checkConvertible(assigns[start:j], typesInfo)
		}
		i = j
	}
}

// parseFieldAssign checks if a statement is of the form dst.Field = src.Field
// and returns the parsed components.
func (w *lintTypeConversion) parseFieldAssign(stmt ast.Stmt) (fieldAssign, bool) {
	assign, ok := stmt.(*ast.AssignStmt)
	if !ok {
		return fieldAssign{}, false
	}

	// Must be a simple assignment (=), not := or +=, etc.
	if assign.Tok.String() != "=" {
		return fieldAssign{}, false
	}

	// Must have exactly one LHS and one RHS.
	if len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
		return fieldAssign{}, false
	}

	// LHS must be a selector expression (dst.Field).
	lhsSel, ok := assign.Lhs[0].(*ast.SelectorExpr)
	if !ok {
		return fieldAssign{}, false
	}

	// RHS must be a selector expression (src.Field).
	rhsSel, ok := assign.Rhs[0].(*ast.SelectorExpr)
	if !ok {
		return fieldAssign{}, false
	}

	// The field names must match.
	if lhsSel.Sel.Name != rhsSel.Sel.Name {
		return fieldAssign{}, false
	}

	// The bases must be simple identifiers (not nested selectors or calls).
	lhsBase := astutils.GoFmt(lhsSel.X)
	rhsBase := astutils.GoFmt(rhsSel.X)
	if lhsBase == "" || rhsBase == "" {
		return fieldAssign{}, false
	}

	// dst and src must be different variables.
	if lhsBase == rhsBase {
		return fieldAssign{}, false
	}

	return fieldAssign{
		dstObj: lhsBase,
		srcObj: rhsBase,
		field:  lhsSel.Sel.Name,
	}, true
}

// checkConvertible determines whether the group of field assignments covers
// all fields of both struct types and the structs are convertible.
func (w *lintTypeConversion) checkConvertible(assigns []fieldAssign, typesInfo *types.Info) {
	if len(assigns) == 0 {
		return
	}

	// Get the first assignment to retrieve type info from the selector expressions.
	firstStmt := assigns[0].stmt.(*ast.AssignStmt)
	lhsSel := firstStmt.Lhs[0].(*ast.SelectorExpr)
	rhsSel := firstStmt.Rhs[0].(*ast.SelectorExpr)

	// Get the type of dst and src base expressions.
	dstType := w.pkg.TypeOf(lhsSel.X)
	srcType := w.pkg.TypeOf(rhsSel.X)
	if dstType == nil || srcType == nil {
		return
	}

	// Dereference pointers.
	dstType = derefPointer(dstType)
	srcType = derefPointer(srcType)

	// Get the underlying struct types.
	dstStruct, ok := dstType.Underlying().(*types.Struct)
	if !ok {
		return
	}
	srcStruct, ok := srcType.Underlying().(*types.Struct)
	if !ok {
		return
	}

	// Both structs must have the same number of fields.
	if dstStruct.NumFields() != srcStruct.NumFields() {
		return
	}

	// The number of assigned fields must match the struct field count.
	if len(assigns) != dstStruct.NumFields() {
		return
	}

	// Build a set of assigned field names.
	assignedFields := make(map[string]bool, len(assigns))
	for _, a := range assigns {
		assignedFields[a.field] = true
	}

	// Verify that every field of both structs is covered and types match.
	for i := 0; i < dstStruct.NumFields(); i++ {
		dstField := dstStruct.Field(i)
		if !assignedFields[dstField.Name()] {
			return
		}
	}

	// Verify the structs have identical field names and types (in any order).
	// Build maps from field name to type for both.
	srcFieldTypes := make(map[string]types.Type, srcStruct.NumFields())
	for i := 0; i < srcStruct.NumFields(); i++ {
		f := srcStruct.Field(i)
		srcFieldTypes[f.Name()] = f.Type()
	}

	for i := 0; i < dstStruct.NumFields(); i++ {
		dstField := dstStruct.Field(i)
		srcFieldType, ok := srcFieldTypes[dstField.Name()]
		if !ok {
			return // Source struct doesn't have a field with this name.
		}
		if !types.Identical(dstField.Type(), srcFieldType) {
			return // Field types don't match.
		}
	}

	// All checks passed — this is a field-by-field copy that can be a type conversion.
	qualifier := types.RelativeTo(w.pkg.TypesPkg())
	dstTypeName := types.TypeString(dstType, qualifier)
	srcTypeName := types.TypeString(srcType, qualifier)

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       assigns[0].stmt,
		Failure: fmt.Sprintf(
			"use type conversion %s(%s) instead of copying struct fields one by one from %s to %s",
			dstTypeName, assigns[0].srcObj, srcTypeName, dstTypeName,
		),
	})
}

// derefPointer returns the element type if t is a pointer; otherwise t itself.
func derefPointer(t types.Type) types.Type {
	if p, ok := t.(*types.Pointer); ok {
		return p.Elem()
	}
	return t
}
