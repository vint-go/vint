package no_duplicated_imports

import (
	"fmt"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// DuplicatedImportsRule looks for packages that are imported two or more times.
type DuplicatedImportsRule struct{}

// Apply applies the rule to given file.
func (*DuplicatedImportsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	impPaths := map[string]struct{}{}
	for _, imp := range file.AST.Imports {
		path := imp.Path.Value
		_, ok := impPaths[path]
		if ok {
			failures = append(failures, lint.Failure{
				Confidence: 1,
				Failure:    fmt.Sprintf("Package %s already imported", path),
				Node:       imp,
				Category:   lint.FailureCategoryImports,
			})
			continue
		}

		impPaths[path] = struct{}{}
	}

	return failures
}

// Name returns the rule name.
func (*DuplicatedImportsRule) Name() string {
	return "noDuplicatedImports"
}

// Group returns the rule group.
func (*DuplicatedImportsRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*DuplicatedImportsRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
