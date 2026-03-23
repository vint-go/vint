package no_cgi_import

import (
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoCgiImportRule detects the import of the net/http/cgi package, which is blocklisted.
type NoCgiImportRule struct{}

// Apply applies the rule to given file.
func (*NoCgiImportRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, importSpec := range file.AST.Imports {
		if importSpec.Path == nil {
			continue
		}
		importPath := strings.Trim(importSpec.Path.Value, `"`)
		if importPath == "net/http/cgi" {
			failures = append(failures, lint.Failure{
				Confidence: 1,
				Failure:    "import of net/http/cgi package is not allowed due to known security issues",
				Node:       importSpec,
				Category:   lint.FailureCategoryImports,
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoCgiImportRule) Name() string {
	return "noCgiImport"
}

// Group returns the rule group.
func (*NoCgiImportRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoCgiImportRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
