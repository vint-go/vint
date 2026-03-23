package no_redundant_test_main_exit

import (
	"fmt"
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// RedundantTestMainExitRule suggests removing redundant [os.Exit] or [syscall.Exit] calls in TestMain function.
type RedundantTestMainExitRule struct{}

// Apply applies the rule to given file.
func (*RedundantTestMainExitRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if !file.IsTest() || !file.Pkg.IsAtLeastGoVersion(lint.Go115) {
		// skip analysis for non-test files or for Go versions before 1.15
		return failures
	}

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintRedundantTestMainExit{onFailure: onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (*RedundantTestMainExitRule) Name() string {
	return "noRedundantTestMainExit"
}

// Group returns the rule group.
func (*RedundantTestMainExitRule) Group() string {
	return "style"
}

type lintRedundantTestMainExit struct {
	onFailure func(lint.Failure)
}

func (w *lintRedundantTestMainExit) Visit(node ast.Node) ast.Visitor {
	if fd, ok := node.(*ast.FuncDecl); ok {
		if fd.Name.Name != "TestMain" {
			return nil // skip analysis for other functions than TestMain
		}

		return w
	}

	se, ok := node.(*ast.ExprStmt)
	if !ok {
		return w
	}
	ce, ok := se.X.(*ast.CallExpr)
	if !ok {
		return w
	}

	fc, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}
	id, ok := fc.X.(*ast.Ident)
	if !ok {
		return w
	}

	pkg := id.Name
	// skip flag calls because they are commonly used in TestMain
	if pkg == "flag" {
		return w
	}

	fn := fc.Sel.Name
	if isCallToExitFunction(pkg, fn, ce.Args) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryStyle,
			Failure:    fmt.Sprintf("redundant call to %s.%s in TestMain function, the test runner will handle it automatically as of Go 1.15", pkg, fn),
		})
	}

	return w
}

// CacheTier returns the cache tier for this rule.
func (*RedundantTestMainExitRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// exitFuncChecker is a function type that checks whether a function call is an exit function.
type exitFuncChecker func(args []ast.Expr) bool

var alwaysTrue exitFuncChecker = func([]ast.Expr) bool { return true }

// exitFunctions is a map of std packages and functions that are considered as exit functions.
var exitFunctions = map[string]map[string]exitFuncChecker{
	"os":      {"Exit": alwaysTrue},
	"syscall": {"Exit": alwaysTrue},
	"log": {
		"Fatal":   alwaysTrue,
		"Fatalf":  alwaysTrue,
		"Fatalln": alwaysTrue,
		"Panic":   alwaysTrue,
		"Panicf":  alwaysTrue,
		"Panicln": alwaysTrue,
	},
	"flag": {
		"Parse": func([]ast.Expr) bool { return true },
		"NewFlagSet": func(args []ast.Expr) bool {
			if len(args) != 2 {
				return false
			}
			return astutils.IsPkgDotName(args[1], "flag", "ExitOnError")
		},
	},
}

// isCallToExitFunction checks if the function call is a call to an exit function.
func isCallToExitFunction(pkgName, functionName string, callArgs []ast.Expr) bool {
	m, ok := exitFunctions[pkgName]
	if !ok {
		return false
	}

	check, ok := m[functionName]
	if !ok {
		return false
	}

	return check(callArgs)
}
