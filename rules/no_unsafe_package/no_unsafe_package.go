package no_unsafe_package

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// unsafeFunctions lists the unsafe package functions and types to detect.
var unsafeFunctions = []string{
	"Pointer",
	"String",
	"StringData",
	"Slice",
	"SliceData",
}

// NoUnsafePackageRule detects usage of the unsafe package.
type NoUnsafePackageRule struct{}

// Apply applies the rule to given file.
func (*NoUnsafePackageRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoUnsafePackage{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoUnsafePackageRule) Name() string {
	return "noUnsafePackage"
}

// Group returns the rule group.
func (*NoUnsafePackageRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnsafePackageRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoUnsafePackage struct {
	onFailure func(lint.Failure)
}

func (w *lintNoUnsafePackage) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	for _, fn := range unsafeFunctions {
		if astutils.IsPkgDotName(ce.Fun, "unsafe", fn) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       ce,
				Category:   lint.FailureCategoryBadPractice,
				Failure:    "use of unsafe." + fn + " bypasses Go type safety",
			})
			return w
		}
	}

	return w
}
