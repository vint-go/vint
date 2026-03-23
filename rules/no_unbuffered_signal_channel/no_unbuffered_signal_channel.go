package no_unbuffered_signal_channel

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoUnbufferedSignalChannelRule detects misuse of unbuffered os.Signal channels
// as arguments to signal.Notify. The signal.Notify function sends signals in a
// non-blocking manner, so an unbuffered channel may miss signals.
type NoUnbufferedSignalChannelRule struct{}

// Apply applies the rule to given file.
func (r *NoUnbufferedSignalChannelRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUnbufferedSignalChannel{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoUnbufferedSignalChannelRule) Name() string {
	return "noUnbufferedSignalChannel"
}

// Group returns the rule group.
func (*NoUnbufferedSignalChannelRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnbufferedSignalChannelRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintUnbufferedSignalChannel struct {
	onFailure func(lint.Failure)
}

func (w *lintUnbufferedSignalChannel) Visit(node ast.Node) ast.Visitor {
	// We look at function bodies to track local variable assignments
	fn, ok := node.(*ast.FuncDecl)
	if !ok {
		fnLit, ok := node.(*ast.FuncLit)
		if !ok {
			return w
		}
		w.checkBlock(fnLit.Body)
		return nil // don't recurse into the function body again
	}

	if fn.Body != nil {
		w.checkBlock(fn.Body)
	}
	return nil // don't recurse into the function body again
}

// checkBlock scans a block statement to find signal.Notify calls with
// unbuffered channels.
func (w *lintUnbufferedSignalChannel) checkBlock(block *ast.BlockStmt) {
	if block == nil {
		return
	}

	// Collect variable names assigned from unbuffered make(chan ...) calls
	unbufferedChans := map[string]bool{}

	// Walk the entire function body collecting assignments and checking calls
	ast.Inspect(block, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.AssignStmt:
			w.collectUnbufferedChansFromAssign(v, unbufferedChans)
			return true
		case *ast.CallExpr:
			w.checkCallExpr(v, unbufferedChans)
			return false // don't recurse into args of this call
		}
		return true
	})
}

func (w *lintUnbufferedSignalChannel) collectUnbufferedChansFromAssign(assign *ast.AssignStmt, unbufferedChans map[string]bool) {
	for i, rhs := range assign.Rhs {
		if i >= len(assign.Lhs) {
			break
		}
		ident, ok := assign.Lhs[i].(*ast.Ident)
		if !ok {
			continue
		}
		if isUnbufferedMakeChan(rhs) {
			unbufferedChans[ident.Name] = true
		} else if isMakeChanCall(rhs) {
			// Reassignment with a buffered channel removes the flag
			delete(unbufferedChans, ident.Name)
		}
	}
}

func (w *lintUnbufferedSignalChannel) checkCallExpr(ce *ast.CallExpr, unbufferedChans map[string]bool) {
	if !astutils.IsPkgDotName(ce.Fun, "signal", "Notify") {
		return
	}

	if len(ce.Args) < 1 {
		return
	}

	firstArg := ce.Args[0]

	// Case 1: Direct unbuffered make call as argument: signal.Notify(make(chan os.Signal), ...)
	if isUnbufferedMakeChan(firstArg) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryLogic,
			Failure:    "unbuffered os.Signal channel passed to signal.Notify; use a buffered channel with at least capacity 1",
		})
		return
	}

	// Case 2: Variable that was assigned an unbuffered channel
	ident, ok := firstArg.(*ast.Ident)
	if !ok {
		return
	}
	if unbufferedChans[ident.Name] {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryLogic,
			Failure:    "unbuffered os.Signal channel passed to signal.Notify; use a buffered channel with at least capacity 1",
		})
	}
}

// isUnbufferedMakeChan returns true if the expression is a make(chan T) call
// with no buffer size argument (i.e., only 1 argument).
func isUnbufferedMakeChan(expr ast.Expr) bool {
	ce, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	ident, ok := ce.Fun.(*ast.Ident)
	if !ok {
		return false
	}
	if ident.Name != "make" {
		return false
	}

	// make(chan T) has exactly 1 argument; make(chan T, size) has 2
	if len(ce.Args) != 1 {
		return false
	}

	// Check that the first argument is a channel type
	_, ok = ce.Args[0].(*ast.ChanType)
	return ok
}

// isMakeChanCall returns true if the expression is any make(chan T, ...) call.
func isMakeChanCall(expr ast.Expr) bool {
	ce, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	ident, ok := ce.Fun.(*ast.Ident)
	if !ok {
		return false
	}
	if ident.Name != "make" {
		return false
	}

	if len(ce.Args) < 1 {
		return false
	}

	_, ok = ce.Args[0].(*ast.ChanType)
	return ok
}
