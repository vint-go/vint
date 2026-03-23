package no_repeated_numbers

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoRepeatedNumbersRule detects numeric literals (integers and floats) that appear
// multiple times in a file and could be replaced by named constants.
type NoRepeatedNumbersRule struct {
	minOccurrences int
	minValue       float64
	maxValue       float64
	hasMin         bool
	hasMax         bool
}

const defaultMinOccurrences = 3

// Configure validates and applies the rule configuration.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *NoRepeatedNumbersRule) Configure(arguments lint.Arguments) error {
	r.minOccurrences = defaultMinOccurrences
	r.minValue = 0
	r.maxValue = 0
	r.hasMin = false
	r.hasMax = false

	if len(arguments) < 1 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noRepeatedNumbers" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		switch {
		case isRuleOption(k, "min-occurrences"):
			n, ok := v.(int64)
			if !ok {
				return fmt.Errorf(`invalid configuration value for min-occurrences in "noRepeatedNumbers" rule; need int64 but got %T`, v)
			}
			r.minOccurrences = int(n)
		case isRuleOption(k, "min"):
			switch val := v.(type) {
			case int64:
				r.minValue = float64(val)
				r.hasMin = true
			case float64:
				r.minValue = val
				r.hasMin = true
			default:
				return fmt.Errorf(`invalid configuration value for min in "noRepeatedNumbers" rule; need numeric but got %T`, v)
			}
		case isRuleOption(k, "max"):
			switch val := v.(type) {
			case int64:
				r.maxValue = float64(val)
				r.hasMax = true
			case float64:
				r.maxValue = val
				r.hasMax = true
			default:
				return fmt.Errorf(`invalid configuration value for max in "noRepeatedNumbers" rule; need numeric but got %T`, v)
			}
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoRepeatedNumbersRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if r.minOccurrences == 0 {
		r.minOccurrences = defaultMinOccurrences
	}

	// Collect all numeric literals and their positions.
	collector := &numberCollector{}
	ast.Walk(collector, file.AST)

	// Count occurrences of each numeric value.
	counts := map[string][]ast.Node{}
	for _, entry := range collector.numbers {
		val := entry.value
		counts[val] = append(counts[val], entry.node)
	}

	var failures []lint.Failure
	reported := map[string]bool{}

	for val, nodes := range counts {
		if len(nodes) < r.minOccurrences {
			continue
		}

		// Apply min/max filtering if configured.
		if r.hasMin || r.hasMax {
			numVal := parseNumericValue(val)
			if r.hasMin && numVal < r.minValue {
				continue
			}
			if r.hasMax && r.maxValue != 0 && numVal > r.maxValue {
				continue
			}
		}

		if reported[val] {
			continue
		}
		reported[val] = true

		// Report failure on the first occurrence.
		failures = append(failures, lint.Failure{
			Confidence: 1,
			Category:   lint.FailureCategoryStyle,
			Failure:    fmt.Sprintf("numeric literal %s appears %d times, consider extracting it into a named constant", val, len(nodes)),
			Node:       nodes[0],
		})
	}

	return failures
}

// Name returns the rule name.
func (*NoRepeatedNumbersRule) Name() string {
	return "noRepeatedNumbers"
}

// Group returns the rule group.
func (*NoRepeatedNumbersRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoRepeatedNumbersRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type numberEntry struct {
	value string
	node  ast.Node
}

type numberCollector struct {
	numbers []numberEntry
}

func (c *numberCollector) Visit(node ast.Node) ast.Visitor {
	lit, ok := node.(*ast.BasicLit)
	if !ok {
		return c
	}

	if lit.Kind != token.INT && lit.Kind != token.FLOAT {
		return c
	}

	// Normalize the value string.
	val := normalizeNumber(lit.Value)
	c.numbers = append(c.numbers, numberEntry{value: val, node: lit})

	return c
}

// normalizeNumber normalizes a numeric literal for comparison purposes.
// It strips underscores and lowercases hex letters for consistent comparison.
func normalizeNumber(s string) string {
	s = strings.ReplaceAll(s, "_", "")
	return s
}

// parseNumericValue parses a numeric string into a float64 for min/max comparison.
func parseNumericValue(s string) float64 {
	var val float64
	// Try parsing as float first, then int.
	_, err := fmt.Sscanf(s, "%g", &val)
	if err != nil {
		// Try hex/octal/binary.
		var intVal int64
		_, err = fmt.Sscanf(s, "%v", &intVal)
		if err == nil {
			return float64(intVal)
		}
	}
	return val
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name lowercased and without hyphens.
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}
