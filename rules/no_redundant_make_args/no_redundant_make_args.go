package no_redundant_make_args

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoRedundantMakeArgsRule detects redundant arguments in make() calls.
type NoRedundantMakeArgsRule struct{}

// Apply applies the rule to given file.
func (r *NoRedundantMakeArgsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRedundantMakeArgs{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoRedundantMakeArgsRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRedundantMakeArgs{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoRedundantMakeArgsRule) Name() string {
	return "noRedundantMakeArgs"
}

// Group returns the rule group.
func (*NoRedundantMakeArgsRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoRedundantMakeArgsRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintRedundantMakeArgs struct {
	onFailure func(lint.Failure)
}

func (w *lintRedundantMakeArgs) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	ident, ok := ce.Fun.(*ast.Ident)
	if !ok {
		return w
	}

	if ident.Name != "make" {
		return w
	}

	if len(ce.Args) < 1 {
		return w
	}

	switch ce.Args[0].(type) {
	case *ast.ArrayType:
		// make([]T, len, cap) — check if len == cap
		w.checkSliceMake(ce)
	case *ast.MapType:
		// make(map[K]V, 0) — check if size hint is 0
		w.checkMapMake(ce)
	}

	return w
}

func (w *lintRedundantMakeArgs) checkSliceMake(ce *ast.CallExpr) {
	// make([]T, len, cap) has 3 args
	if len(ce.Args) != 3 {
		return
	}

	lenArg := ce.Args[1]
	capArg := ce.Args[2]

	// Compare the string representations of len and cap arguments
	lenStr := astutils.GoFmt(lenArg)
	capStr := astutils.GoFmt(capArg)

	if lenStr == capStr {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryStyle,
			Failure:    "redundant capacity argument in make, length and capacity are the same",
		})
	}
}

func (w *lintRedundantMakeArgs) checkMapMake(ce *ast.CallExpr) {
	// make(map[K]V, 0) has 2 args
	if len(ce.Args) != 2 {
		return
	}

	sizeArg := ce.Args[1]

	// Check if the size hint is the literal 0
	lit, ok := sizeArg.(*ast.BasicLit)
	if !ok {
		return
	}

	if lit.Kind == token.INT && lit.Value == "0" {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryStyle,
			Failure:    "redundant size hint 0 in make, can be omitted",
		})
	}
}
