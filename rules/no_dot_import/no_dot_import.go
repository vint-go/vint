package no_dot_import

import (
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoDotImportRule flags dot imports as discouraged.
type NoDotImportRule struct{}

// Apply applies the rule to given file.
func (r *NoDotImportRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, importSpec := range file.AST.Imports {
		if importSpec.Name != nil && importSpec.Name.Name == "." {
			failures = append(failures, lint.Failure{
				Confidence: 1,
				Failure:    "should not use dot imports",
				Node:       importSpec,
				Category:   lint.FailureCategoryImports,
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoDotImportRule) Name() string {
	return "noDotImport"
}

// Group returns the rule group.
func (*NoDotImportRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoDotImportRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
