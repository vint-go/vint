package no_unrecognized_directive

import (
	"fmt"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// knownDirectives is the set of recognized Go compiler directive names.
var knownDirectives = map[string]bool{
	"build":                true,
	"embed":                true,
	"generate":             true,
	"linkname":             true,
	"noinline":             true,
	"nosplit":              true,
	"noescape":             true,
	"norace":               true,
	"nocheckptr":           true,
	"nointerface":          true,
	"nowritebarrier":       true,
	"nowritebarrierrec":    true,
	"yeswritebarrierrec":   true,
	"systemstack":          true,
	"uintptrescapes":       true,
	"uintptrkeepalive":     true,
	"notinheap":            true,
	"cgo_dynamic_linker":   true,
	"cgo_export_dynamic":   true,
	"cgo_export_static":    true,
	"cgo_import_dynamic":   true,
	"cgo_import_static":    true,
	"cgo_ldflag":           true,
	"cgo_unsafe_args":      true,
	"debug":                true,
	"wasmimport":           true,
	"wasmexport":           true,
}

// NoUnrecognizedDirectiveRule detects Go compiler directives that use an
// unrecognized directive name.
type NoUnrecognizedDirectiveRule struct{}

// Apply applies the rule to given file.
func (r *NoUnrecognizedDirectiveRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, group := range file.AST.Comments {
		for _, comment := range group.List {
			text := comment.Text
			if !strings.HasPrefix(text, "//go:") {
				continue
			}

			// Extract the directive name (everything after "//go:" up to
			// the first space or end of string).
			rest := strings.TrimPrefix(text, "//go:")
			directive := rest
			if idx := strings.IndexByte(rest, ' '); idx != -1 {
				directive = rest[:idx]
			}

			if directive == "" {
				continue
			}

			if !knownDirectives[directive] {
				failures = append(failures, lint.Failure{
					Category:   lint.FailureCategoryBadPractice,
					Confidence: 1,
					Node:       comment,
					Failure:    fmt.Sprintf("compiler directive unrecognized: %s", directive),
				})
			}
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoUnrecognizedDirectiveRule) Name() string {
	return "noUnrecognizedDirective"
}

// Group returns the rule group.
func (*NoUnrecognizedDirectiveRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnrecognizedDirectiveRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
