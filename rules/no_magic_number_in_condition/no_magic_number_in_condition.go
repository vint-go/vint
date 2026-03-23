package no_magic_number_in_condition

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoMagicNumberInConditionRule detects magic numbers used in conditional
// expressions (if statements). A magic number is a numeric literal (integer
// or floating-point) that is not defined as a constant and whose meaning is
// not immediately clear. Using magic numbers in conditions makes code less
// readable and harder to maintain.
type NoMagicNumberInConditionRule struct {
	ignoredNumbers map[string]bool
	ignoredFiles   []string
}

var defaultIgnoredNumbers = map[string]bool{
	"0":   true,
	"0.0": true,
	"1":   true,
	"1.0": true,
}

var defaultIgnoredFiles = []string{"_test.go"}

// Configure validates and applies the rule configuration.
//
// Configure implements the [lint.ConfigurableRule] interface.
func (r *NoMagicNumberInConditionRule) Configure(arguments lint.Arguments) error {
	r.ignoredNumbers = copyMap(defaultIgnoredNumbers)
	r.ignoredFiles = append([]string{}, defaultIgnoredFiles...)

	if len(arguments) < 1 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noMagicNumberInCondition" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		switch {
		case isRuleOption(k, "ignored-numbers"):
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf(`invalid configuration value for ignored-numbers in "noMagicNumberInCondition" rule; need string but got %T`, v)
			}
			for k, v := range parseIgnoredNumbers(s) {
				r.ignoredNumbers[k] = v
			}
		case isRuleOption(k, "ignored-files"):
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf(`invalid configuration value for ignored-files in "noMagicNumberInCondition" rule; need string but got %T`, v)
			}
			r.ignoredFiles = append(r.ignoredFiles, parseIgnoredFiles(s)...)
		}
	}

	return nil
}

// Apply applies the rule to the given file.
func (r *NoMagicNumberInConditionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if r.ignoredNumbers == nil {
		r.ignoredNumbers = copyMap(defaultIgnoredNumbers)
		r.ignoredFiles = append([]string{}, defaultIgnoredFiles...)
	}

	// Check if the file should be ignored.
	for _, pattern := range r.ignoredFiles {
		if strings.HasSuffix(file.Name, pattern) {
			return nil
		}
	}

	var failures []lint.Failure

	w := &lintMagicNumberInCondition{
		ignoredNumbers: r.ignoredNumbers,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoMagicNumberInConditionRule) Name() string {
	return "noMagicNumberInCondition"
}

// Group returns the rule group.
func (*NoMagicNumberInConditionRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoMagicNumberInConditionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintMagicNumberInCondition struct {
	ignoredNumbers map[string]bool
	onFailure      func(lint.Failure)
}

func (w *lintMagicNumberInCondition) Visit(node ast.Node) ast.Visitor {
	ifStmt, ok := node.(*ast.IfStmt)
	if !ok {
		return w
	}

	// Check the condition expression of the if statement.
	w.checkCondition(ifStmt.Cond)

	return w
}

func (w *lintMagicNumberInCondition) checkCondition(expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		// For comparison operators, check both sides for magic numbers.
		if isComparisonOp(e.Op) {
			w.checkExprForMagicNumber(e.X)
			w.checkExprForMagicNumber(e.Y)
		}
		// For logical operators (&&, ||), recurse into both sides to find
		// nested comparisons.
		if e.Op == token.LAND || e.Op == token.LOR {
			w.checkCondition(e.X)
			w.checkCondition(e.Y)
		}
	case *ast.ParenExpr:
		w.checkCondition(e.X)
	}
}

func (w *lintMagicNumberInCondition) checkExprForMagicNumber(expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.INT || e.Kind == token.FLOAT {
			val := normalizeNumber(e.Value)
			if !w.ignoredNumbers[val] {
				w.onFailure(lint.Failure{
					Confidence: 1,
					Category:   lint.FailureCategoryStyle,
					Failure:    fmt.Sprintf("magic number: %s, in <condition> detected", val),
					Node:       e,
				})
			}
		}
	case *ast.UnaryExpr:
		w.checkExprForMagicNumber(e.X)
	case *ast.ParenExpr:
		w.checkExprForMagicNumber(e.X)
	}
}

func isComparisonOp(op token.Token) bool {
	switch op {
	case token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ:
		return true
	}
	return false
}

func normalizeNumber(s string) string {
	s = strings.ReplaceAll(s, "_", "")

	// Try to normalize integer representations to decimal.
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") ||
		strings.HasPrefix(s, "0o") || strings.HasPrefix(s, "0O") ||
		strings.HasPrefix(s, "0b") || strings.HasPrefix(s, "0B") {
		if n, err := strconv.ParseInt(s, 0, 64); err == nil {
			return strconv.FormatInt(n, 10)
		}
	}

	// For floats with trailing/leading zeros, normalize.
	if strings.Contains(s, ".") {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			normalized := strconv.FormatFloat(f, 'f', -1, 64)
			return normalized
		}
	}

	return s
}

func parseIgnoredNumbers(s string) map[string]bool {
	result := make(map[string]bool)
	for _, num := range strings.Split(s, ",") {
		num = strings.TrimSpace(num)
		if num != "" {
			result[num] = true
		}
	}
	return result
}

func parseIgnoredFiles(s string) []string {
	var result []string
	for _, f := range strings.Split(s, ",") {
		f = strings.TrimSpace(f)
		if f != "" {
			result = append(result, f)
		}
	}
	return result
}

func copyMap(m map[string]bool) map[string]bool {
	result := make(map[string]bool, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}
