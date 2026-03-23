package no_malformed_directive

import (
	"fmt"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// knownDirectives is the set of recognized Go compiler directive names.
var knownDirectives = map[string]bool{
	"build":              true,
	"embed":              true,
	"generate":           true,
	"linkname":           true,
	"noinline":           true,
	"nosplit":            true,
	"noescape":           true,
	"norace":             true,
	"nocheckptr":         true,
	"nointerface":        true,
	"nowritebarrier":     true,
	"nowritebarrierrec":  true,
	"yeswritebarrierrec": true,
	"systemstack":        true,
	"uintptrescapes":     true,
	"uintptrkeepalive":   true,
	"notinheap":          true,
	"cgo_dynamic_linker": true,
	"cgo_export_dynamic": true,
	"cgo_export_static":  true,
	"cgo_import_dynamic": true,
	"cgo_import_static":  true,
	"cgo_ldflag":         true,
	"cgo_unsafe_args":    true,
	"debug":              true,
	"wasmimport":         true,
	"wasmexport":         true,
}

// NoMalformedDirectiveRule checks known Go toolchain directives for correctness.
// It validates that Go compiler directives (comments starting with //go:) are
// correctly formed and placed.
type NoMalformedDirectiveRule struct{}

// Apply applies the rule to given file.
func (r *NoMalformedDirectiveRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, group := range file.AST.Comments {
		for _, comment := range group.List {
			text := comment.Text

			// Check for space between // and go: (e.g. "// go:noinline")
			if strings.HasPrefix(text, "// go:") {
				rest := strings.TrimPrefix(text, "// ")
				directive := rest
				if idx := strings.IndexByte(rest, ' '); idx != -1 {
					directive = rest[:idx]
				}
				// Only flag if it looks like a known directive name
				directiveName := strings.TrimPrefix(directive, "go:")
				if knownDirectives[directiveName] || findClosestDirective(directiveName) != "" {
					failures = append(failures, lint.Failure{
						Category:   lint.FailureCategoryBadPractice,
						Confidence: 1,
						Node:       comment,
						Failure:    fmt.Sprintf("malformed directive: space between // and go: in %q; use %q", text, "//"+rest),
					})
				}
				continue
			}

			// Check actual //go: directives
			if !strings.HasPrefix(text, "//go:") {
				continue
			}

			// Extract the directive name (everything after "//go:" up to
			// the first space or end of string).
			rest := strings.TrimPrefix(text, "//go:")
			directiveName := rest
			if idx := strings.IndexByte(rest, ' '); idx != -1 {
				directiveName = rest[:idx]
			}

			if directiveName == "" {
				continue
			}

			// Skip known valid directives
			if knownDirectives[directiveName] {
				continue
			}

			// Unknown directive: flag it. If close to a known one, suggest the correction.
			closest := findClosestDirective(directiveName)
			if closest != "" {
				failures = append(failures, lint.Failure{
					Category:   lint.FailureCategoryBadPractice,
					Confidence: 1,
					Node:       comment,
					Failure:    fmt.Sprintf("malformed directive: possible misspelling of //go:%s in %q", closest, text),
				})
			} else {
				failures = append(failures, lint.Failure{
					Category:   lint.FailureCategoryBadPractice,
					Confidence: 1,
					Node:       comment,
					Failure:    fmt.Sprintf("unrecognized directive %q", text),
				})
			}
		}
	}

	return failures
}

// findClosestDirective returns the closest known directive name if the
// given name is within edit distance 2, or empty string if no close match.
func findClosestDirective(name string) string {
	bestMatch := ""
	bestDist := 3 // threshold: only suggest if distance <= 2

	for known := range knownDirectives {
		d := levenshtein(name, known)
		if d < bestDist {
			bestDist = d
			bestMatch = known
		}
	}

	return bestMatch
}

// levenshtein computes the Levenshtein distance between two strings.
func levenshtein(a, b string) int {
	la := len(a)
	lb := len(b)

	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	// Use a single row for the DP, reducing space from O(la*lb) to O(lb).
	prev := make([]int, lb+1)
	for j := 0; j <= lb; j++ {
		prev[j] = j
	}

	for i := 1; i <= la; i++ {
		curr := make([]int, lb+1)
		curr[0] = i
		for j := 1; j <= lb; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min3(
				curr[j-1]+1,   // insertion
				prev[j]+1,     // deletion
				prev[j-1]+cost, // substitution
			)
		}
		prev = curr
	}

	return prev[lb]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// Name returns the rule name.
func (*NoMalformedDirectiveRule) Name() string {
	return "noMalformedDirective"
}

// Group returns the rule group.
func (*NoMalformedDirectiveRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoMalformedDirectiveRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
