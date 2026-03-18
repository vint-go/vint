package use_compiled_regexp

import (
	"fmt"
	"go/ast"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseCompiledRegexpRule detects calls to regexp.Match, regexp.MatchString,
// or regexp.MatchReader inside loops, where the regex is recompiled on each
// iteration. The fix is to compile the regex once outside the loop.
type UseCompiledRegexpRule struct{}

// regexpFuncs lists the regexp package-level functions that compile
// a regex on every call and should be replaced with a pre-compiled regexp.
var regexpFuncs = map[string]bool{
	"Match":       true,
	"MatchString": true,
	"MatchReader": true,
}

// Apply applies the rule to the given file.
func (r *UseCompiledRegexpRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintCompiledRegexp{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseCompiledRegexpRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintCompiledRegexp{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseCompiledRegexpRule) Name() string {
	return "useCompiledRegexp"
}

// Group returns the rule group.
func (*UseCompiledRegexpRule) Group() string {
	return "performance"
}

// CacheTier returns the cache tier for this rule.
func (*UseCompiledRegexpRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintCompiledRegexp struct {
	onFailure func(lint.Failure)
}

func (w *lintCompiledRegexp) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.ForStmt:
		w.findRegexpCallInBlock(n.Body)
		return nil
	case *ast.RangeStmt:
		w.findRegexpCallInBlock(n.Body)
		return nil
	}
	return w
}

// findRegexpCallInBlock walks a block statement inside a loop looking for
// calls to regexp.Match, regexp.MatchString, or regexp.MatchReader.
// It does not descend into nested function literals (closures) since those
// might have their own scope and the call may not run on every loop iteration.
func (w *lintCompiledRegexp) findRegexpCallInBlock(block *ast.BlockStmt) {
	if block == nil {
		return
	}
	ast.Inspect(block, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.FuncLit:
			// Don't descend into closures.
			return false
		}

		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		pkgIdent, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}

		if pkgIdent.Name != "regexp" {
			return true
		}

		funcName := sel.Sel.Name
		if !regexpFuncs[funcName] {
			return true
		}

		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       call,
			Category:   lint.FailureCategoryOptimization,
			Failure:    fmt.Sprintf("calling regexp.%s in a loop; compile the regexp once outside the loop with regexp.Compile", funcName),
		})

		return true
	})
}
