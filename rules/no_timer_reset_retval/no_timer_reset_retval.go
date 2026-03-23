package no_timer_reset_retval

import (
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoTimerResetRetvalRule detects usage of the return value of (*time.Timer).Reset,
// which is unreliable and should not be used.
type NoTimerResetRetvalRule struct{}

// Apply applies the rule to given file.
func (r *NoTimerResetRetvalRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if file.Pkg.TypeCheck() != nil {
		return nil
	}

	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	w := &lintTimerResetRetval{
		file:      file,
		onFailure: onFailure,
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoTimerResetRetvalRule) Name() string {
	return "noTimerResetRetval"
}

// Group returns the rule group.
func (*NoTimerResetRetvalRule) Group() string {
	return "correctness"
}

// RequiresTypecheck returns true because this rule needs type info.
func (*NoTimerResetRetvalRule) RequiresTypecheck() bool {
	return true
}

// CacheTier returns the cache tier for this rule.
func (*NoTimerResetRetvalRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

type lintTimerResetRetval struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintTimerResetRetval) Visit(node ast.Node) ast.Visitor {
	// We need to find cases where the return value of t.Reset() is used.
	// The return value is used when the call is NOT a standalone ExprStmt.
	// We look for parent nodes that contain call expressions to Reset on *time.Timer.

	switch n := node.(type) {
	case *ast.AssignStmt:
		// Check if any RHS is a timer Reset call
		for _, rhs := range n.Rhs {
			if w.isTimerResetCall(rhs) {
				w.onFailure(lint.Failure{
					Confidence: 1,
					Node:       n,
					Category:   lint.FailureCategoryLogic,
					Failure:    "the return value of (*time.Timer).Reset is unreliable and should not be used",
				})
			}
		}
		return w

	case *ast.IfStmt:
		// Check if the condition contains a timer Reset call
		if w.containsTimerResetCall(n.Cond) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       n,
				Category:   lint.FailureCategoryLogic,
				Failure:    "the return value of (*time.Timer).Reset is unreliable and should not be used",
			})
		}
		// Also check Init statement if it's an assignment
		if init, ok := n.Init.(*ast.AssignStmt); ok {
			for _, rhs := range init.Rhs {
				if w.isTimerResetCall(rhs) {
					w.onFailure(lint.Failure{
						Confidence: 1,
						Node:       n,
						Category:   lint.FailureCategoryLogic,
						Failure:    "the return value of (*time.Timer).Reset is unreliable and should not be used",
					})
				}
			}
		}
		return w
	}

	return w
}

// isTimerResetCall checks if the given expression is a call to (*time.Timer).Reset.
func (w *lintTimerResetRetval) isTimerResetCall(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Reset" {
		return false
	}

	// Check that the receiver type is *time.Timer
	t := w.file.Pkg.TypeOf(sel.X)
	if t == nil {
		return false
	}

	return isTimeTimerPointer(t)
}

// containsTimerResetCall recursively checks if an expression contains a timer Reset call.
func (w *lintTimerResetRetval) containsTimerResetCall(expr ast.Expr) bool {
	if expr == nil {
		return false
	}

	if w.isTimerResetCall(expr) {
		return true
	}

	switch e := expr.(type) {
	case *ast.UnaryExpr:
		return w.containsTimerResetCall(e.X)
	case *ast.BinaryExpr:
		return w.containsTimerResetCall(e.X) || w.containsTimerResetCall(e.Y)
	case *ast.ParenExpr:
		return w.containsTimerResetCall(e.X)
	}

	return false
}

// isTimeTimerPointer checks if a type is *time.Timer.
func isTimeTimerPointer(t types.Type) bool {
	ptr, ok := t.(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := ptr.Elem().(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	return obj.Name() == "Timer" && obj.Pkg() != nil && obj.Pkg().Path() == "time"
}
