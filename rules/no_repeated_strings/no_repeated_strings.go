package no_repeated_strings

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoRepeatedStringsRule detects string literals that appear multiple times
// in a file and could be replaced by a named constant.
type NoRepeatedStringsRule struct {
	minOccurrences int
	minLength      int
	ignoreStrings  *regexp.Regexp
	ignoreTests    bool
	ignorePattern  *regexp.Regexp
}

const (
	defaultStringMinOccurrences = 3
	defaultStringMinLength      = 3
)

// Configure validates and applies the rule configuration.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *NoRepeatedStringsRule) Configure(arguments lint.Arguments) error {
	r.minOccurrences = defaultStringMinOccurrences
	r.minLength = defaultStringMinLength
	r.ignoreStrings = nil
	r.ignoreTests = true
	r.ignorePattern = nil

	if len(arguments) < 1 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noRepeatedStrings" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		switch {
		case isRuleOption(k, "min-occurrences"):
			n, ok := v.(int64)
			if !ok {
				return fmt.Errorf(`invalid configuration value for min-occurrences in "noRepeatedStrings" rule; need int64 but got %T`, v)
			}
			r.minOccurrences = int(n)
		case isRuleOption(k, "min-length"):
			n, ok := v.(int64)
			if !ok {
				return fmt.Errorf(`invalid configuration value for min-length in "noRepeatedStrings" rule; need int64 but got %T`, v)
			}
			r.minLength = int(n)
		case isRuleOption(k, "ignore-strings"):
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf(`invalid configuration value for ignore-strings in "noRepeatedStrings" rule; need string but got %T`, v)
			}
			if s != "" {
				re, err := regexp.Compile(s)
				if err != nil {
					return fmt.Errorf(`invalid regex for ignore-strings in "noRepeatedStrings" rule: %w`, err)
				}
				r.ignoreStrings = re
			}
		case isRuleOption(k, "ignore-tests"):
			b, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for ignore-tests in "noRepeatedStrings" rule; need bool but got %T`, v)
			}
			r.ignoreTests = b
		case isRuleOption(k, "ignore"):
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf(`invalid configuration value for ignore in "noRepeatedStrings" rule; need string but got %T`, v)
			}
			if s != "" {
				re, err := regexp.Compile(s)
				if err != nil {
					return fmt.Errorf(`invalid regex for ignore in "noRepeatedStrings" rule: %w`, err)
				}
				r.ignorePattern = re
			}
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoRepeatedStringsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if r.minOccurrences == 0 {
		r.minOccurrences = defaultStringMinOccurrences
	}
	if r.minLength == 0 {
		r.minLength = defaultStringMinLength
	}

	// Skip test files if configured.
	if r.ignoreTests && file.IsTest() {
		return nil
	}

	// Skip files matching the ignore pattern.
	if r.ignorePattern != nil && r.ignorePattern.MatchString(file.Name) {
		return nil
	}

	// Collect all string literals and their positions.
	collector := &stringCollector{
		minLength: r.minLength,
	}
	ast.Walk(collector, file.AST)

	// Count occurrences of each string value.
	counts := map[string][]ast.Node{}
	for _, entry := range collector.strings {
		counts[entry.value] = append(counts[entry.value], entry.node)
	}

	var failures []lint.Failure
	reported := map[string]bool{}

	for val, nodes := range counts {
		if len(nodes) < r.minOccurrences {
			continue
		}

		// Apply ignore-strings filtering if configured.
		if r.ignoreStrings != nil && r.ignoreStrings.MatchString(val) {
			continue
		}

		if reported[val] {
			continue
		}
		reported[val] = true

		// Report failure on the first occurrence.
		failures = append(failures, lint.Failure{
			Confidence: 1,
			Category:   lint.FailureCategoryStyle,
			Failure:    fmt.Sprintf("string literal %q appears %d times, consider extracting it into a named constant", val, len(nodes)),
			Node:       nodes[0],
		})
	}

	return failures
}

// Name returns the rule name.
func (*NoRepeatedStringsRule) Name() string {
	return "noRepeatedStrings"
}

// Group returns the rule group.
func (*NoRepeatedStringsRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoRepeatedStringsRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type stringEntry struct {
	value string
	node  ast.Node
}

type stringCollector struct {
	strings   []stringEntry
	minLength int
}

func (c *stringCollector) Visit(node ast.Node) ast.Visitor {
	lit, ok := node.(*ast.BasicLit)
	if !ok {
		return c
	}

	if lit.Kind != token.STRING {
		return c
	}

	// Extract the actual string value (strip quotes).
	val := unquoteString(lit.Value)

	// Skip strings below the minimum length.
	if len(val) < c.minLength {
		return c
	}

	c.strings = append(c.strings, stringEntry{value: val, node: lit})

	return c
}

// unquoteString removes surrounding quotes from a string literal.
// It handles both double-quoted and backtick-quoted strings.
func unquoteString(s string) string {
	if len(s) < 2 {
		return s
	}
	if s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	if s[0] == '`' && s[len(s)-1] == '`' {
		return s[1 : len(s)-1]
	}
	return s
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name lowercased and without hyphens.
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}
