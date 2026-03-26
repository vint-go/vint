package use_filepath_join

import (
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseFilepathJoinRule detects path concatenation using string(os.PathSeparator)
// that can be replaced with filepath.Join.
type UseFilepathJoinRule struct{}

// Apply applies the rule to given file.
func (r *UseFilepathJoinRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintFilepathJoin{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseFilepathJoinRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintFilepathJoin{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseFilepathJoinRule) Name() string {
	return "useFilepathJoin"
}

// Group returns the rule group.
func (*UseFilepathJoinRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*UseFilepathJoinRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintFilepathJoin struct {
	onFailure func(lint.Failure)
}

func (w *lintFilepathJoin) Visit(node ast.Node) ast.Visitor {
	binExpr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	if binExpr.Op != token.ADD {
		return w
	}

	if containsPathSeparatorConversion(binExpr) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       binExpr,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "path concatenation using string(os.PathSeparator) can be replaced with filepath.Join",
		})
		return nil // don't recurse into children to avoid duplicate reports
	}

	return w
}

// containsPathSeparatorConversion checks if a binary ADD expression tree
// contains a string(os.PathSeparator) type conversion.
func containsPathSeparatorConversion(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		if e.Op != token.ADD {
			return false
		}
		return containsPathSeparatorConversion(e.X) || containsPathSeparatorConversion(e.Y)
	case *ast.CallExpr:
		return isStringOsPathSeparator(e)
	default:
		return false
	}
}

// isStringOsPathSeparator returns true if the call expression is string(os.PathSeparator).
func isStringOsPathSeparator(call *ast.CallExpr) bool {
	// Must be a type conversion: Fun is *ast.Ident with Name "string"
	ident, ok := call.Fun.(*ast.Ident)
	if !ok || ident.Name != "string" {
		return false
	}

	// Must have exactly 1 argument
	if len(call.Args) != 1 {
		return false
	}

	// Argument must be os.PathSeparator
	sel, ok := call.Args[0].(*ast.SelectorExpr)
	if !ok {
		return false
	}

	x, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	return x.Name == "os" && sel.Sel.Name == "PathSeparator"
}
