package use_filepath_join

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseFilepathJoinRule detects path concatenation using + and "/" that can be
// replaced with filepath.Join.
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

	if containsPathSeparatorLiteral(binExpr) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       binExpr,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "path concatenation can be replaced with filepath.Join",
		})
		return nil // don't recurse into children to avoid duplicate reports
	}

	return w
}

// containsPathSeparatorLiteral checks if a binary ADD expression tree contains
// a string literal with a path separator (/ or \).
func containsPathSeparatorLiteral(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		if e.Op != token.ADD {
			return false
		}
		return containsPathSeparatorLiteral(e.X) || containsPathSeparatorLiteral(e.Y)
	case *ast.BasicLit:
		if e.Kind != token.STRING {
			return false
		}
		// Extract the string value (remove quotes)
		val := e.Value
		if len(val) >= 2 {
			val = val[1 : len(val)-1]
		}
		return strings.ContainsAny(val, "/\\")
	default:
		return false
	}
}
