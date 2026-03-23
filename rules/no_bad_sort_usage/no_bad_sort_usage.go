package no_bad_sort_usage

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoBadSortUsageRule detects suspicious usage of sort.IntSlice, sort.Float64Slice,
// and sort.StringSlice type conversions when the simpler sort.Ints, sort.Float64s,
// or sort.Strings convenience functions should be used instead.
type NoBadSortUsageRule struct{}

// sortTypeMapping maps sort wrapper types to their simpler convenience functions
// and the expected underlying slice type.
var sortTypeMapping = map[string]struct {
	suggestion string
	sliceType  string
}{
	"IntSlice":     {suggestion: "sort.Ints", sliceType: "[]int"},
	"Float64Slice": {suggestion: "sort.Float64s", sliceType: "[]float64"},
	"StringSlice":  {suggestion: "sort.Strings", sliceType: "[]string"},
}

// Apply applies the rule to given file.
func (r *NoBadSortUsageRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if file.Pkg.TypeCheck() != nil {
		return nil
	}

	w := &lintBadSortUsage{
		file: file,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoBadSortUsageRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintBadSortUsage{file: file, onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoBadSortUsageRule) Name() string {
	return "noBadSortUsage"
}

// Group returns the rule group.
func (*NoBadSortUsageRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoBadSortUsageRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck indicates this rule needs type information.
func (*NoBadSortUsageRule) RequiresTypecheck() bool {
	return true
}

type lintBadSortUsage struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintBadSortUsage) Visit(node ast.Node) ast.Visitor {
	assign, ok := node.(*ast.AssignStmt)
	if !ok {
		return w
	}

	// Only check plain assignments (=)
	if assign.Tok.String() != "=" {
		return w
	}

	// Must have matching LHS/RHS counts
	if len(assign.Lhs) != len(assign.Rhs) {
		return w
	}

	for i := range assign.Lhs {
		w.checkAssignment(assign, assign.Lhs[i], assign.Rhs[i])
	}

	return w
}

func (w *lintBadSortUsage) checkAssignment(assign *ast.AssignStmt, lhs, rhs ast.Expr) {
	// RHS must be a call expression like sort.IntSlice(x)
	call, ok := rhs.(*ast.CallExpr)
	if !ok {
		return
	}

	// The function must be a selector expression like sort.IntSlice
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	// Check if the selector name is one of our target types
	mapping, ok := sortTypeMapping[sel.Sel.Name]
	if !ok {
		return
	}

	// Verify it's actually the sort package via type info
	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok {
		return
	}

	typesInfo := w.file.Pkg.TypesInfo()
	if typesInfo != nil {
		obj, ok := typesInfo.Uses[pkgIdent].(*types.PkgName)
		if !ok || obj.Imported().Path() != "sort" {
			return
		}
	}

	// Must have exactly one argument
	if len(call.Args) != 1 {
		return
	}

	// Check that the argument type matches the expected slice type
	argType := w.file.Pkg.TypeOf(call.Args[0])
	if argType == nil {
		return
	}

	if argType.String() != mapping.sliceType {
		return
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryLogic,
		Confidence: 1,
		Node:       assign,
		Failure:    fmt.Sprintf("suspicious sort.%s usage, maybe %s was intended", sel.Sel.Name, mapping.suggestion),
	})
}
