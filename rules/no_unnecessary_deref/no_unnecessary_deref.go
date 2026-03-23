package no_unnecessary_deref

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoUnnecessaryDerefRule detects dereference expressions that can be omitted.
// Go automatically dereferences pointers for field access and array indexing,
// so explicit dereferencing is unnecessary in many cases.
type NoUnnecessaryDerefRule struct {
	skipRecvDeref bool
	configured    bool
}

const defaultSkipRecvDeref = true

// Configure validates and applies the rule configuration.
func (r *NoUnnecessaryDerefRule) Configure(arguments lint.Arguments) error {
	r.skipRecvDeref = defaultSkipRecvDeref
	r.configured = true

	if len(arguments) == 0 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noUnnecessaryDeref" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		if normalizeOption(k) == "skiprecvderef" {
			boolVal, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for skipRecvDeref in "noUnnecessaryDeref" rule; need bool but got %T`, v)
			}
			r.skipRecvDeref = boolVal
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoUnnecessaryDerefRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if !r.configured {
		r.skipRecvDeref = defaultSkipRecvDeref
		r.configured = true
	}

	file.Pkg.TypeCheck()

	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUnnecessaryDeref{
		onFailure:     onFailure,
		pkg:           file.Pkg,
		skipRecvDeref: r.skipRecvDeref,
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoUnnecessaryDerefRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	if !r.configured {
		r.skipRecvDeref = defaultSkipRecvDeref
		r.configured = true
	}

	file.Pkg.TypeCheck()

	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUnnecessaryDeref{
		onFailure:     onFailure,
		pkg:           file.Pkg,
		skipRecvDeref: r.skipRecvDeref,
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoUnnecessaryDerefRule) Name() string {
	return "noUnnecessaryDeref"
}

// Group returns the rule group.
func (*NoUnnecessaryDerefRule) Group() string {
	return "style"
}

// RequiresTypecheck returns true because this rule uses type information.
func (*NoUnnecessaryDerefRule) RequiresTypecheck() bool {
	return true
}

// CacheTier returns the cache tier for this rule.
func (*NoUnnecessaryDerefRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

type lintUnnecessaryDeref struct {
	onFailure     func(lint.Failure)
	pkg           *lint.Package
	skipRecvDeref bool
}

func (w *lintUnnecessaryDeref) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.SelectorExpr:
		w.checkSelector(n)
	case *ast.IndexExpr:
		w.checkIndex(n)
	}
	return w
}

// checkSelector checks for patterns like (*x).field where x is a pointer.
func (w *lintUnnecessaryDeref) checkSelector(sel *ast.SelectorExpr) {
	parenExpr, ok := sel.X.(*ast.ParenExpr)
	if !ok {
		return
	}

	starExpr, ok := parenExpr.X.(*ast.StarExpr)
	if !ok {
		return
	}

	// Verify the inner expression is actually a pointer type
	innerType := w.pkg.TypeOf(starExpr.X)
	if innerType == nil {
		return
	}

	_, isPtr := innerType.Underlying().(*types.Pointer)
	if !isPtr {
		return
	}

	// If skipRecvDeref is enabled, check if this is a method call on a pointer receiver
	if w.skipRecvDeref {
		if w.isPointerReceiverMethodCall(sel) {
			return
		}
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       sel,
		Category:   lint.FailureCategoryStyle,
		Failure:    fmt.Sprintf("expression can be simplified to %s.%s", exprString(starExpr.X), sel.Sel.Name),
	})
}

// checkIndex checks for patterns like (*a)[i] where a is a pointer to an array.
func (w *lintUnnecessaryDeref) checkIndex(idx *ast.IndexExpr) {
	parenExpr, ok := idx.X.(*ast.ParenExpr)
	if !ok {
		return
	}

	starExpr, ok := parenExpr.X.(*ast.StarExpr)
	if !ok {
		return
	}

	// Verify the inner expression is a pointer to an array
	innerType := w.pkg.TypeOf(starExpr.X)
	if innerType == nil {
		return
	}

	ptrType, isPtr := innerType.Underlying().(*types.Pointer)
	if !isPtr {
		return
	}

	// Go only auto-dereferences pointers to arrays for indexing, not slices
	_, isArray := ptrType.Elem().Underlying().(*types.Array)
	if !isArray {
		return
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       idx,
		Category:   lint.FailureCategoryStyle,
		Failure:    fmt.Sprintf("expression can be simplified to %s[%s]", exprString(starExpr.X), exprString(idx.Index)),
	})
}

// isPointerReceiverMethodCall checks if the selector expression is part of a
// method call where the method has a pointer receiver.
func (w *lintUnnecessaryDeref) isPointerReceiverMethodCall(sel *ast.SelectorExpr) bool {
	// Look up the selection info from types.Info
	typesInfo := w.pkg.TypesInfo()
	if typesInfo == nil {
		return false
	}

	selection, ok := typesInfo.Selections[sel]
	if !ok {
		return false
	}

	// Check if this is a method value/call (not a field access)
	if selection.Kind() != types.MethodVal {
		return false
	}

	// Check if the method has a pointer receiver
	fn, ok := selection.Obj().(*types.Func)
	if !ok {
		return false
	}

	sig, ok := fn.Type().(*types.Signature)
	if !ok {
		return false
	}

	recv := sig.Recv()
	if recv == nil {
		return false
	}

	_, isPtr := recv.Type().(*types.Pointer)
	return isPtr
}

// exprString returns a simple string representation of an expression.
func exprString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return exprString(e.X) + "." + e.Sel.Name
	case *ast.BasicLit:
		return e.Value
	case *ast.CallExpr:
		return exprString(e.Fun) + "(...)"
	case *ast.IndexExpr:
		return exprString(e.X) + "[" + exprString(e.Index) + "]"
	case *ast.StarExpr:
		return "*" + exprString(e.X)
	case *ast.ParenExpr:
		return "(" + exprString(e.X) + ")"
	default:
		return "..."
	}
}

// normalizeOption normalizes a configuration option name by removing hyphens,
// underscores, and lowering case.
func normalizeOption(name string) string {
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "_", "")
	return strings.ToLower(name)
}
