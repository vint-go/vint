package no_duplicate_import

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoDuplicateImportRule detects multiple imports of the same package under different aliases.
type NoDuplicateImportRule struct{}

// Apply applies the rule to given file.
func (r *NoDuplicateImportRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	// Map from import path (unquoted) to the first import spec seen
	seen := map[string]*ast.ImportSpec{}

	for _, importSpec := range file.AST.Imports {
		importPath := strings.Trim(importSpec.Path.Value, `"`)
		if prev, ok := seen[importPath]; ok {
			_ = prev
			failures = append(failures, lint.Failure{
				Confidence: 1,
				Failure:    fmt.Sprintf("package %q is already imported", importPath),
				Node:       importSpec,
				Category:   lint.FailureCategoryImports,
			})
		} else {
			seen[importPath] = importSpec
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoDuplicateImportRule) Name() string {
	return "noDuplicateImport"
}

// Group returns the rule group.
func (*NoDuplicateImportRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoDuplicateImportRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
