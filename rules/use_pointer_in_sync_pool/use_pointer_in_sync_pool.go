package use_pointer_in_sync_pool

import (
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UsePointerInSyncPoolRule detects non-pointer values stored in sync.Pool,
// which cause heap allocations due to interface boxing.
type UsePointerInSyncPoolRule struct{}

// Apply applies the rule to the given file.
func (r *UsePointerInSyncPoolRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()
	typesInfo := file.Pkg.TypesInfo()
	if typesInfo == nil {
		return nil
	}

	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	w := &lintSyncPool{
		pkg:       file.Pkg,
		typesInfo: typesInfo,
		onFailure: onFailure,
	}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*UsePointerInSyncPoolRule) Name() string {
	return "usePointerInSyncPool"
}

// Group returns the rule group.
func (*UsePointerInSyncPoolRule) Group() string {
	return "performance"
}

// CacheTier returns the cache tier for this rule.
func (*UsePointerInSyncPoolRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule needs type information.
func (*UsePointerInSyncPoolRule) RequiresTypecheck() bool {
	return true
}

type lintSyncPool struct {
	pkg       *lint.Package
	typesInfo *types.Info
	onFailure func(lint.Failure)
}

func (w *lintSyncPool) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	w.checkPutCall(call)

	return w
}

// checkPutCall checks if a sync.Pool.Put() call passes a non-pointer value.
func (w *lintSyncPool) checkPutCall(call *ast.CallExpr) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	if sel.Sel.Name != "Put" {
		return
	}

	if len(call.Args) != 1 {
		return
	}

	// Check that the receiver is a *sync.Pool
	if !w.isSyncPoolType(sel.X) {
		return
	}

	// Check if the argument is a non-pointer type
	arg := call.Args[0]
	argType := w.typesInfo.TypeOf(arg)
	if argType == nil {
		return
	}

	if isPointerLike(argType) {
		return
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryOptimization,
		Confidence: 1,
		Node:       call,
		Failure:    "non-pointer value stored in sync.Pool causes allocation; use a pointer instead",
	})
}

// isSyncPoolType checks if the expression refers to a sync.Pool value.
func (w *lintSyncPool) isSyncPoolType(expr ast.Expr) bool {
	t := w.pkg.TypeOf(expr)
	if t == nil {
		return false
	}

	// Dereference pointer if needed (sync.Pool is usually used as *sync.Pool)
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}

	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	obj := named.Obj()
	return obj != nil && obj.Pkg() != nil && obj.Pkg().Path() == "sync" && obj.Name() == "Pool"
}

// isPointerLike returns true if the type is a pointer, unsafe.Pointer, or interface
// (values that don't require boxing when stored in interface{}).
func isPointerLike(t types.Type) bool {
	switch t.Underlying().(type) {
	case *types.Pointer:
		return true
	case *types.Interface:
		// Already an interface, no boxing needed.
		return true
	}
	return false
}
