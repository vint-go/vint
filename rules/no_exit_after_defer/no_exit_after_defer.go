package no_exit_after_defer

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoExitAfterDeferRule detects calls to os.Exit, log.Fatal, log.Fatalf, and
// log.Fatalln inside functions that use defer.
type NoExitAfterDeferRule struct{}

// Apply applies the rule to given file.
func (r *NoExitAfterDeferRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintExitAfterDefer{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoExitAfterDeferRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintExitAfterDefer{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoExitAfterDeferRule) Name() string {
	return "noExitAfterDefer"
}

// Group returns the rule group.
func (*NoExitAfterDeferRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoExitAfterDeferRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintExitAfterDefer struct {
	onFailure func(lint.Failure)
}

func (w *lintExitAfterDefer) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		if n.Body != nil {
			w.checkFuncBody(n.Body)
		}
		return nil // don't recurse into function declarations again
	case *ast.FuncLit:
		if n.Body != nil {
			w.checkFuncBody(n.Body)
		}
		return nil // don't recurse into function literals again
	}
	return w
}

// checkFuncBody checks whether a function body contains both defer statements
// and exit/fatal calls. If it does, it reports failures for the exit/fatal calls.
func (w *lintExitAfterDefer) checkFuncBody(body *ast.BlockStmt) {
	if !hasDefer(body) {
		return
	}

	// Find all exit/fatal calls and report them
	exitFinder := &exitCallFinder{onFailure: w.onFailure}
	ast.Walk(exitFinder, body)
}

// hasDefer checks whether a block statement (or any nested blocks that are NOT
// nested function literals) contains a defer statement.
func hasDefer(node ast.Node) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if found {
			return false
		}
		// Don't look inside nested function literals; they have their own scope.
		if _, ok := n.(*ast.FuncLit); ok {
			return false
		}
		if _, ok := n.(*ast.DeferStmt); ok {
			found = true
			return false
		}
		return true
	})
	return found
}

// exitCallFinder walks a function body to find calls to os.Exit, log.Fatal,
// log.Fatalf, and log.Fatalln. It skips nested function literals.
type exitCallFinder struct {
	onFailure func(lint.Failure)
}

func (f *exitCallFinder) Visit(node ast.Node) ast.Visitor {
	// Don't descend into nested function literals; they have their own scope.
	if _, ok := node.(*ast.FuncLit); ok {
		return nil
	}

	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return f
	}

	if astutils.IsPkgDotName(ce.Fun, "os", "Exit") {
		f.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryLogic,
			Failure:    "calling os.Exit in a function that uses defer will skip deferred cleanup",
		})
	} else if astutils.IsPkgDotName(ce.Fun, "log", "Fatal") {
		f.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryLogic,
			Failure:    "calling log.Fatal in a function that uses defer will skip deferred cleanup",
		})
	} else if astutils.IsPkgDotName(ce.Fun, "log", "Fatalf") {
		f.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryLogic,
			Failure:    "calling log.Fatalf in a function that uses defer will skip deferred cleanup",
		})
	} else if astutils.IsPkgDotName(ce.Fun, "log", "Fatalln") {
		f.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryLogic,
			Failure:    "calling log.Fatalln in a function that uses defer will skip deferred cleanup",
		})
	}

	return f
}
