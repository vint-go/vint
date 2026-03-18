package use_direct_method_call

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseDirectMethodCallRule detects method expression calls that can be replaced
// with direct method calls. For example, T.Method(receiver) can be written as
// receiver.Method().
type UseDirectMethodCallRule struct{}

// Apply applies the rule to given file.
func (r *UseDirectMethodCallRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintDirectMethodCall{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseDirectMethodCallRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintDirectMethodCall{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*UseDirectMethodCallRule) Name() string {
	return "useDirectMethodCall"
}

// Group returns the rule group.
func (*UseDirectMethodCallRule) Group() string {
	return "style"
}

// RequiresTypecheck returns true because the rule needs type information.
func (*UseDirectMethodCallRule) RequiresTypecheck() bool {
	return true
}

// CacheTier returns the cache tier for this rule.
func (*UseDirectMethodCallRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

type lintDirectMethodCall struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintDirectMethodCall) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	// Check if the selector expression is a method expression using type info.
	info := w.pkg.TypesInfo()
	if info == nil || info.Selections == nil {
		return w
	}

	selection, ok := info.Selections[sel]
	if !ok {
		return w
	}

	if selection.Kind() != types.MethodExpr {
		return w
	}

	// This is a method expression call like T.Method(receiver, args...)
	// Suggest using receiver.Method(args...) instead.
	methodName := sel.Sel.Name

	// The first argument is the receiver.
	if len(call.Args) == 0 {
		return w
	}

	receiverStr := astutils.GoFmt(call.Args[0])

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       call,
		Category:   lint.FailureCategoryStyle,
		Failure:    fmt.Sprintf("method expression call can be replaced with %s.%s()", receiverStr, methodName),
	})

	return w
}
