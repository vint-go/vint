package no_infinite_recursion

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoInfiniteRecursionRule detects functions that call themselves unconditionally
// in all code paths, which will cause a stack overflow.
type NoInfiniteRecursionRule struct{}

// Apply applies the rule to given file.
func (r *NoInfiniteRecursionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}

		var receiver *ast.Ident
		switch {
		case funcDecl.Recv == nil:
			receiver = nil
		case funcDecl.Recv.NumFields() < 1 || len(funcDecl.Recv.List[0].Names) < 1:
			receiver = &ast.Ident{Name: "_"}
		default:
			receiver = funcDecl.Recv.List[0].Names[0]
		}

		w := &lintInfiniteRecursion{
			onFailure:   onFailure,
			receiver:    receiver,
			funcName:    funcDecl.Name,
			seenCondExit: false,
		}

		ast.Walk(w, funcDecl.Body)
	}

	return failures
}

// Name returns the rule name.
func (*NoInfiniteRecursionRule) Name() string {
	return "noInfiniteRecursion"
}

// Group returns the rule group.
func (*NoInfiniteRecursionRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoInfiniteRecursionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintInfiniteRecursion struct {
	onFailure    func(lint.Failure)
	receiver     *ast.Ident
	funcName     *ast.Ident
	seenCondExit bool
	inGoStmt     bool
}

func (w *lintInfiniteRecursion) isRecursiveCall(call *ast.CallExpr) bool {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		// Direct function call: foo()
		if w.receiver != nil {
			return false
		}
		return fun.Name == w.funcName.Name
	case *ast.SelectorExpr:
		// Method call: r.Foo()
		ident, ok := fun.X.(*ast.Ident)
		if !ok {
			return false
		}
		if w.receiver == nil {
			return false
		}
		return ident.Name == w.receiver.Name && fun.Sel.Name == w.funcName.Name
	}
	return false
}

// Visit traverses the function body looking for unconditional recursive calls.
// It skips inside conditional control structures (if, for with condition, switch, select, range)
// but checks them for control-flow exits (return, panic, os.Exit etc.) that would
// prevent unconditional recursion. If a recursive call is found before any conditional exit,
// it is reported as infinite recursion.
func (w *lintInfiniteRecursion) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.CallExpr:
		// Check call arguments for recursive calls first
		for _, arg := range n.Args {
			ast.Walk(w, arg)
		}

		if _, ok := n.Fun.(*ast.FuncLit); ok {
			ast.Walk(w, n.Fun.(*ast.FuncLit).Body)
			return nil
		}

		if !w.seenCondExit && w.isRecursiveCall(n) {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 0.8,
				Node:       n,
				Failure:    "infinite recursive call",
			})
		}
		return nil

	case *ast.IfStmt:
		w.updateSeenCondExit(n.Body)
		w.updateSeenCondExit(n.Else)
		return nil

	case *ast.SelectStmt:
		w.updateSeenCondExit(n.Body)
		return nil

	case *ast.RangeStmt:
		w.updateSeenCondExit(n.Body)
		return nil

	case *ast.TypeSwitchStmt:
		w.updateSeenCondExit(n.Body)
		return nil

	case *ast.SwitchStmt:
		w.updateSeenCondExit(n.Body)
		return nil

	case *ast.GoStmt:
		w.inGoStmt = true
		ast.Walk(w, n.Call)
		w.inGoStmt = false
		return nil

	case *ast.ForStmt:
		if n.Cond != nil {
			// Conditional loop — check for exits inside
			w.updateSeenCondExit(n.Body)
			return nil
		}
		// Unconditional loop (for {}) — continue walking inside
		return w

	case *ast.FuncLit:
		if w.inGoStmt {
			return w
		}
		return nil // closure is a different scope
	}

	return w
}

func (w *lintInfiniteRecursion) updateSeenCondExit(node ast.Node) {
	if node == nil || w.seenCondExit {
		return
	}
	w.seenCondExit = hasControlExit(node)
}

// hasControlExit returns true if the node contains a control flow statement
// that exits the function (return, panic, os.Exit, log.Fatal, etc.).
func hasControlExit(node ast.Node) bool {
	isExit := func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.ReturnStmt:
			return true
		case *ast.CallExpr:
			if astutils.IsIdent(stmt.Fun, "panic") {
				return true
			}
			se, ok := stmt.Fun.(*ast.SelectorExpr)
			if !ok {
				return false
			}
			id, ok := se.X.(*ast.Ident)
			if !ok {
				return false
			}
			return isCallToExitFunc(id.Name, se.Sel.Name)
		}
		return false
	}

	return astutils.SeekNode[ast.Node](node, isExit) != nil
}

// isCallToExitFunc checks if a package.function call is one that terminates the program.
func isCallToExitFunc(pkg, fn string) bool {
	switch pkg {
	case "os":
		return fn == "Exit"
	case "log":
		return fn == "Fatal" || fn == "Fatalf" || fn == "Fatalln" ||
			fn == "Panic" || fn == "Panicf" || fn == "Panicln"
	}
	return false
}
