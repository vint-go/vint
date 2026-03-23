package no_immediate_new_deref

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoImmediateNewDerefRule detects immediate dereferencing of new expressions.
type NoImmediateNewDerefRule struct{}

// Apply applies the rule to given file.
func (r *NoImmediateNewDerefRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	file.Pkg.TypeCheck()

	w := &lintNoImmediateNewDeref{
		onFailure: onFailure,
		pkg:       file.Pkg,
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoImmediateNewDerefRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	file.Pkg.TypeCheck()

	w := &lintNoImmediateNewDeref{
		onFailure: onFailure,
		pkg:       file.Pkg,
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoImmediateNewDerefRule) Name() string {
	return "noImmediateNewDeref"
}

// Group returns the rule group.
func (*NoImmediateNewDerefRule) Group() string {
	return "style"
}

// RequiresTypecheck returns true because this rule uses type information.
func (*NoImmediateNewDerefRule) RequiresTypecheck() bool {
	return true
}

// CacheTier returns the cache tier for this rule.
func (*NoImmediateNewDerefRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

type lintNoImmediateNewDeref struct {
	onFailure func(lint.Failure)
	pkg       *lint.Package
}

func (w *lintNoImmediateNewDeref) Visit(node ast.Node) ast.Visitor {
	starExpr, ok := node.(*ast.StarExpr)
	if !ok {
		return w
	}

	callExpr, ok := starExpr.X.(*ast.CallExpr)
	if !ok {
		return w
	}

	ident, ok := callExpr.Fun.(*ast.Ident)
	if !ok || ident.Name != "new" {
		return w
	}

	// Check that this is the builtin new, not a user-defined function named "new"
	if ident.Obj != nil {
		return w
	}

	// Check if the argument is a type parameter
	if len(callExpr.Args) == 1 {
		argType := w.pkg.TypeOf(callExpr.Args[0])
		if argType != nil {
			if _, isTypeParam := argType.(*types.TypeParam); isTypeParam {
				return w
			}
		}
	}

	// Build a descriptive message
	var typeName string
	if len(callExpr.Args) == 1 {
		switch arg := callExpr.Args[0].(type) {
		case *ast.Ident:
			typeName = arg.Name
		case *ast.SelectorExpr:
			if x, ok := arg.X.(*ast.Ident); ok {
				typeName = x.Name + "." + arg.Sel.Name
			}
		}
	}

	msg := "immediate dereference of new expression"
	if typeName != "" {
		msg = fmt.Sprintf("immediate dereference of new(%s), use zero value instead", typeName)
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       starExpr,
		Category:   lint.FailureCategoryStyle,
		Failure:    msg,
	})

	return w
}
