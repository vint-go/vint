package no_test_fatal_in_goroutine

import (
	"go/ast"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// fatalMethods lists testing.T/testing.B methods that call runtime.Goexit()
// and therefore must not be called from a goroutine other than the test goroutine.
var fatalMethods = map[string]bool{
	"Fatal":   true,
	"Fatalf":  true,
	"FailNow": true,
}

// NoTestFatalInGoroutineRule detects calls to Fatal, Fatalf, FailNow and similar
// methods from testing.T or testing.B inside goroutines spawned by the test.
type NoTestFatalInGoroutineRule struct{}

// Apply applies the rule to given file.
func (r *NoTestFatalInGoroutineRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// Check if this is a test or benchmark function
		testParam := testingParam(funcDecl)
		if testParam == "" {
			continue
		}

		// Walk the function body looking for go statements
		if funcDecl.Body != nil {
			w := &goStmtFinder{
				testParam: testParam,
				onFailure: func(f lint.Failure) {
					failures = append(failures, f)
				},
			}
			ast.Walk(w, funcDecl.Body)
		}
	}

	return failures
}

// testingParam returns the name of the *testing.T or *testing.B parameter
// of a test/benchmark function, or "" if this is not a test/benchmark function.
func testingParam(funcDecl *ast.FuncDecl) string {
	name := funcDecl.Name.Name
	if !strings.HasPrefix(name, "Test") && !strings.HasPrefix(name, "Benchmark") && !strings.HasPrefix(name, "Fuzz") {
		return ""
	}

	if funcDecl.Type.Params == nil || len(funcDecl.Type.Params.List) == 0 {
		return ""
	}

	for _, param := range funcDecl.Type.Params.List {
		starExpr, ok := param.Type.(*ast.StarExpr)
		if !ok {
			continue
		}
		selExpr, ok := starExpr.X.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		ident, ok := selExpr.X.(*ast.Ident)
		if !ok {
			continue
		}
		if ident.Name == "testing" && (selExpr.Sel.Name == "T" || selExpr.Sel.Name == "B" || selExpr.Sel.Name == "F") {
			if len(param.Names) > 0 {
				return param.Names[0].Name
			}
		}
	}

	return ""
}

// goStmtFinder walks the AST of a test function looking for go statements.
type goStmtFinder struct {
	testParam string
	onFailure func(lint.Failure)
}

func (w *goStmtFinder) Visit(node ast.Node) ast.Visitor {
	goStmt, ok := node.(*ast.GoStmt)
	if !ok {
		return w
	}

	// Found a go statement; now search within it for fatal calls on the test param
	fc := &fatalCallFinder{
		testParam: w.testParam,
		onFailure: w.onFailure,
	}
	ast.Walk(fc, goStmt.Call)

	// Don't descend into the go statement again from the outer walker
	return w
}

// fatalCallFinder walks within a goroutine body looking for t.Fatal(), t.Fatalf(), t.FailNow() etc.
type fatalCallFinder struct {
	testParam string
	onFailure func(lint.Failure)
}

func (w *fatalCallFinder) Visit(node ast.Node) ast.Visitor {
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	ident, ok := selExpr.X.(*ast.Ident)
	if !ok {
		return w
	}

	if ident.Name != w.testParam {
		return w
	}

	methodName := selExpr.Sel.Name
	if fatalMethods[methodName] {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       callExpr,
			Category:   lint.FailureCategoryLogic,
			Failure:    w.testParam + "." + methodName + " must not be called from a non-test goroutine",
		})
	}

	return w
}

// Name returns the rule name.
func (*NoTestFatalInGoroutineRule) Name() string {
	return "noTestFatalInGoroutine"
}

// Group returns the rule group.
func (*NoTestFatalInGoroutineRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoTestFatalInGoroutineRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
