package no_init_function

import (
	"go/ast"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoInitFunctionRule disallows the use of init() functions in Go packages.
type NoInitFunctionRule struct{}

// Apply applies the rule to given file.
func (r *NoInitFunctionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if funcDecl.Name.Name == "init" && funcDecl.Recv == nil {
			failures = append(failures, lint.Failure{
				Category:   lint.FailureCategoryBadPractice,
				Confidence: 1,
				Node:       funcDecl,
				Failure:    "init function found, consider using explicit initialization",
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoInitFunctionRule) Name() string {
	return "noInitFunction"
}

// Group returns the rule group.
func (*NoInitFunctionRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoInitFunctionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
