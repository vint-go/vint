package no_exposed_pprof

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoExposedPprofRule detects when net/http/pprof is imported as a blank import,
// which automatically exposes profiling endpoints on the default HTTP mux.
type NoExposedPprofRule struct{}

// Apply applies the rule to given file.
func (r *NoExposedPprofRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, imp := range file.AST.Imports {
		if imp.Path == nil {
			continue
		}

		// Check if this is a blank import of net/http/pprof
		if imp.Path.Value == `"net/http/pprof"` && isBlank(imp.Name) {
			failures = append(failures, lint.Failure{
				Confidence: 1,
				Node:       imp,
				Category:   lint.FailureCategoryBadPractice,
				Failure:    "blank import of net/http/pprof automatically exposes profiling endpoints on the default mux",
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoExposedPprofRule) Name() string {
	return "noExposedPprof"
}

// Group returns the rule group.
func (*NoExposedPprofRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoExposedPprofRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// isBlank returns whether id is the blank identifier "_".
// If id == nil, the answer is false.
func isBlank(id *ast.Ident) bool { return id != nil && id.Name == "_" }
