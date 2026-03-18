package no_benchmark_n_assignment

import (
	"go/ast"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoBenchmarkNAssignmentRule detects assignments to b.N in benchmark functions.
type NoBenchmarkNAssignmentRule struct{}

// Apply applies the rule to given file.
func (r *NoBenchmarkNAssignmentRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintBenchmarkNAssignment{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoBenchmarkNAssignmentRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintBenchmarkNAssignment{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoBenchmarkNAssignmentRule) Name() string {
	return "noBenchmarkNAssignment"
}

// Group returns the rule group.
func (*NoBenchmarkNAssignmentRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoBenchmarkNAssignmentRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintBenchmarkNAssignment struct {
	onFailure func(lint.Failure)
}

func (w *lintBenchmarkNAssignment) Visit(node ast.Node) ast.Visitor {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return w
	}

	// Check if this is a benchmark function: name starts with "Benchmark"
	if funcDecl.Name == nil || !strings.HasPrefix(funcDecl.Name.Name, "Benchmark") {
		return w
	}

	// Check that it has exactly one parameter of type *testing.B
	bParamName := benchmarkParamName(funcDecl)
	if bParamName == "" {
		return w
	}

	// Walk the function body looking for assignments to bParamName.N
	if funcDecl.Body != nil {
		finder := &nAssignmentFinder{
			paramName: bParamName,
			onFailure: w.onFailure,
		}
		ast.Walk(finder, funcDecl.Body)
	}

	return nil // don't recurse into the function again
}

// benchmarkParamName returns the name of the *testing.B parameter if the
// function has the correct benchmark signature, or "" otherwise.
func benchmarkParamName(funcDecl *ast.FuncDecl) string {
	if funcDecl.Type == nil || funcDecl.Type.Params == nil {
		return ""
	}

	params := funcDecl.Type.Params.List
	if len(params) != 1 {
		return ""
	}

	param := params[0]
	// Check that the type is *testing.B
	starExpr, ok := param.Type.(*ast.StarExpr)
	if !ok {
		return ""
	}

	selExpr, ok := starExpr.X.(*ast.SelectorExpr)
	if !ok {
		return ""
	}

	ident, ok := selExpr.X.(*ast.Ident)
	if !ok {
		return ""
	}

	if ident.Name != "testing" || selExpr.Sel.Name != "B" {
		return ""
	}

	if len(param.Names) == 0 {
		return ""
	}

	return param.Names[0].Name
}

// nAssignmentFinder walks a function body to find assignments to paramName.N.
type nAssignmentFinder struct {
	paramName string
	onFailure func(lint.Failure)
}

func (f *nAssignmentFinder) Visit(node ast.Node) ast.Visitor {
	assignStmt, ok := node.(*ast.AssignStmt)
	if !ok {
		return f
	}

	for _, lhs := range assignStmt.Lhs {
		selExpr, ok := lhs.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		ident, ok := selExpr.X.(*ast.Ident)
		if !ok {
			continue
		}

		if ident.Name == f.paramName && selExpr.Sel.Name == "N" {
			f.onFailure(lint.Failure{
				Confidence: 1,
				Node:       assignStmt,
				Category:   lint.FailureCategoryLogic,
				Failure:    "assignment to b.N in benchmark distorts the results",
			})
		}
	}

	return f
}
