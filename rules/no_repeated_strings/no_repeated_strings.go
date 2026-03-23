package no_repeated_strings

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoRepeatedStringsRule detects string literals that appear multiple times
// in a file and could be replaced by a named constant.
type NoRepeatedStringsRule struct {
	minOccurrences       int
	minLength            int
	ignoreStrings        *regexp.Regexp
	ignoreTests          bool
	ignoreCalls          bool
	ignorePattern        *regexp.Regexp
	evalConstExpressions bool
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
	r.ignoreCalls = true
	r.ignorePattern = nil
	r.evalConstExpressions = false

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
		case isRuleOption(k, "ignore-calls"):
			b, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for ignore-calls in "noRepeatedStrings" rule; need bool but got %T`, v)
			}
			r.ignoreCalls = b
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
		case isRuleOption(k, "eval-const-expressions"):
			b, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for eval-const-expressions in "noRepeatedStrings" rule; need bool but got %T`, v)
			}
			r.evalConstExpressions = b
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

	// Build a constant map if eval-const-expressions is enabled.
	var constMap map[string]string
	if r.evalConstExpressions {
		constMap = buildConstMap(file.AST)
	}

	// Collect all string literals and their positions.
	collector := &stringCollector{
		minLength: r.minLength,
		constMap:  constMap,
	}
	ast.Walk(collector, file.AST)

	// Group occurrences of each string value.
	groups := map[string][]stringEntry{}
	for _, entry := range collector.strings {
		groups[entry.value] = append(groups[entry.value], entry)
	}

	var failures []lint.Failure
	reported := map[string]bool{}

	for val, entries := range groups {
		if len(entries) < r.minOccurrences {
			continue
		}

		// Apply ignore-strings filtering if configured.
		if r.ignoreStrings != nil && r.ignoreStrings.MatchString(val) {
			continue
		}

		// Apply ignore-calls filtering: skip strings that only appear in call arguments.
		if r.ignoreCalls {
			allInCalls := true
			for _, e := range entries {
				if !e.inCall {
					allInCalls = false
					break
				}
			}
			if allInCalls {
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
			Failure:    fmt.Sprintf("string literal %q appears %d times, consider extracting it into a named constant", val, len(entries)),
			Node:       entries[0].node,
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
	value  string
	node   ast.Node
	inCall bool
}

type stringCollector struct {
	strings      []stringEntry
	minLength    int
	callArgDepth int
	constMap     map[string]string // constant name -> resolved string value (nil if eval-const-expressions is off)
}

func (c *stringCollector) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	// Track call expression arguments so we can mark string literals
	// that appear only within function call arguments.
	if callExpr, ok := node.(*ast.CallExpr); ok {
		// Walk the function expression (not a call argument).
		ast.Walk(c, callExpr.Fun)
		// Walk each argument with incremented callArgDepth.
		c.callArgDepth++
		for _, arg := range callExpr.Args {
			ast.Walk(c, arg)
		}
		c.callArgDepth--
		return nil // prevent ast.Walk from re-walking children
	}

	// When eval-const-expressions is enabled, try to resolve binary
	// concatenation expressions involving string constants.
	if c.constMap != nil {
		if binExpr, ok := node.(*ast.BinaryExpr); ok && binExpr.Op == token.ADD {
			if val, ok := c.resolveStringExpr(binExpr); ok && len(val) >= c.minLength {
				c.strings = append(c.strings, stringEntry{value: val, node: binExpr, inCall: c.callArgDepth > 0})
				return nil // don't walk children; we already resolved the whole expression
			}
		}
	}

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

	c.strings = append(c.strings, stringEntry{value: val, node: lit, inCall: c.callArgDepth > 0})

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

// buildConstMap walks the AST top-level declarations and collects string
// constants into a map of name -> resolved value.
func buildConstMap(file *ast.File) map[string]string {
	m := map[string]string{}
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}
		for _, spec := range genDecl.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, val := range vs.Values {
				if i >= len(vs.Names) {
					break
				}
				if resolved, ok := resolveConstExpr(val, m); ok {
					m[vs.Names[i].Name] = resolved
				}
			}
		}
	}
	return m
}

// resolveConstExpr tries to evaluate an expression as a string constant.
// It handles string literals, identifiers referring to known constants, and
// binary ADD expressions that concatenate strings.
func resolveConstExpr(expr ast.Expr, constMap map[string]string) (string, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			return unquoteString(e.Value), true
		}
		return "", false
	case *ast.Ident:
		if v, ok := constMap[e.Name]; ok {
			return v, true
		}
		return "", false
	case *ast.BinaryExpr:
		if e.Op != token.ADD {
			return "", false
		}
		left, lok := resolveConstExpr(e.X, constMap)
		if !lok {
			return "", false
		}
		right, rok := resolveConstExpr(e.Y, constMap)
		if !rok {
			return "", false
		}
		return left + right, true
	case *ast.ParenExpr:
		return resolveConstExpr(e.X, constMap)
	}
	return "", false
}

// resolveStringExpr tries to evaluate a binary expression to its concatenated
// string value using the collector's constant map.
func (c *stringCollector) resolveStringExpr(expr ast.Expr) (string, bool) {
	return resolveConstExpr(expr, c.constMap)
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name lowercased and without hyphens.
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}
