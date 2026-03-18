package no_test_main_without_exit

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoTestMainWithoutExitRule detects TestMain functions that do not call os.Exit,
// which causes test failures to be hidden since the binary always exits with code 0.
type NoTestMainWithoutExitRule struct{}

// Apply applies the rule to given file.
func (r *NoTestMainWithoutExitRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintTestMainExit{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoTestMainWithoutExitRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintTestMainExit{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoTestMainWithoutExitRule) Name() string {
	return "noTestMainWithoutExit"
}

// Group returns the rule group.
func (*NoTestMainWithoutExitRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoTestMainWithoutExitRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintTestMainExit struct {
	onFailure func(lint.Failure)
}

func (w *lintTestMainExit) Visit(node ast.Node) ast.Visitor {
	fn, ok := node.(*ast.FuncDecl)
	if !ok {
		return w
	}

	// Check if this is func TestMain(m *testing.M)
	if !astutils.FuncSignatureIs(fn, "TestMain", []string{"*testing.M"}, nil) {
		return w
	}

	// Check if the function body contains a call to os.Exit
	if fn.Body != nil && !hasOsExit(fn.Body) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       fn,
			Category:   lint.FailureCategoryLogic,
			Failure:    "TestMain should call os.Exit to set exit code, otherwise test failures will be hidden",
		})
	}

	return nil // don't recurse into the function
}

// hasOsExit checks whether the given block contains a call to os.Exit.
func hasOsExit(node ast.Node) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if found {
			return false
		}
		ce, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if astutils.IsPkgDotName(ce.Fun, "os", "Exit") {
			found = true
			return false
		}
		return true
	})
	return found
}
