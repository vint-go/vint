package no_filesystem_root_serving

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoFilesystemRootServingRule detects the use of http.Dir("/") as a potential directory traversal risk.
type NoFilesystemRootServingRule struct{}

// Apply applies the rule to given file.
func (r *NoFilesystemRootServingRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintFilesystemRootServing{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoFilesystemRootServingRule) Name() string {
	return "noFilesystemRootServing"
}

// Group returns the rule group.
func (*NoFilesystemRootServingRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoFilesystemRootServingRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintFilesystemRootServing struct {
	onFailure func(lint.Failure)
}

func (w *lintFilesystemRootServing) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check for http.Dir("/")
	if !astutils.IsPkgDotName(ce.Fun, "http", "Dir") {
		return w
	}

	// http.Dir takes exactly one argument
	if len(ce.Args) != 1 {
		return w
	}

	arg := ce.Args[0]
	lit, ok := arg.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return w
	}

	val, err := strconv.Unquote(lit.Value)
	if err != nil {
		return w
	}

	if val == "/" {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "serving the entire filesystem root via http.Dir is a security risk",
		})
	}

	return w
}
