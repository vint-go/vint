package no_permissive_os_create

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoPermissiveOsCreateRule detects usage of os.Create which creates files
// with default 0666 permissions (before umask), potentially making them
// readable and writable by all users.
type NoPermissiveOsCreateRule struct{}

// Apply applies the rule to given file.
func (r *NoPermissiveOsCreateRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintPermissiveOsCreate{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoPermissiveOsCreateRule) Name() string {
	return "noPermissiveOsCreate"
}

// Group returns the rule group.
func (*NoPermissiveOsCreateRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoPermissiveOsCreateRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintPermissiveOsCreate struct {
	onFailure func(lint.Failure)
}

func (w *lintPermissiveOsCreate) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "os", "Create") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "os.Create uses default 0666 permissions; use os.OpenFile with explicit permissions instead",
	})

	return w
}
