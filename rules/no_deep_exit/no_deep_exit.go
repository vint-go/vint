package no_deep_exit

import (
	"fmt"
	"go/ast"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// DeepExitRule lints program exit in functions other than main or init.
type DeepExitRule struct{}

// Apply applies the rule to given file.
func (*DeepExitRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintDeepExit{onFailure: onFailure, isTestFile: file.IsTest()}
	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (*DeepExitRule) Name() string {
	return "noDeepExit"
}

// Group returns the rule group.
func (*DeepExitRule) Group() string {
	return "correctness"
}

type lintDeepExit struct {
	onFailure  func(lint.Failure)
	isTestFile bool
}

func (w *lintDeepExit) Visit(node ast.Node) ast.Visitor {
	if fd, ok := node.(*ast.FuncDecl); ok {
		if w.mustIgnore(fd) {
			return nil // skip analysis of this function
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
	fn := fc.Sel.Name
	if isCallToExitFunction(pkg, fn, ce.Args) {
		msg := fmt.Sprintf("calls to %s.%s only in main() or init() functions", pkg, fn)

		if pkg == "flag" && fn == "NewFlagSet" &&
			len(ce.Args) == 2 && astutils.IsPkgDotName(ce.Args[1], "flag", "ExitOnError") {
			msg = "calls to flag.NewFlagSet with flag.ExitOnError only in main() or init() functions"
		}
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    msg,
		})
	}

	return w
}

func (w *lintDeepExit) mustIgnore(fd *ast.FuncDecl) bool {
	fn := fd.Name.Name

	return fn == "init" || fn == "main" || w.isTestMain(fd) || w.isTestExample(fd)
}

func (w *lintDeepExit) isTestMain(fd *ast.FuncDecl) bool {
	return w.isTestFile && fd.Name.Name == "TestMain"
}

// isTestExample returns true if the function is a testable example function.
// See https://go.dev/blog/examples#examples-are-tests for more information.
//
// Inspired by https://github.com/golang/go/blob/go1.23.0/src/go/doc/example.go#L72-L77
func (w *lintDeepExit) isTestExample(fd *ast.FuncDecl) bool {
	if !w.isTestFile {
		return false
	}
	name := fd.Name.Name
	const prefix = "Example"
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	if len(name) == len(prefix) { // "Example" is a package level example
		return len(fd.Type.Params.List) == 0
	}
	r, _ := utf8.DecodeRuneInString(name[len(prefix):])
	if unicode.IsLower(r) {
		return false
	}
	return len(fd.Type.Params.List) == 0
}

// CacheTier returns the cache tier for this rule.
func (*DeepExitRule) CacheTier() rulecache.CacheTier {
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
