package no_magic_number_in_operation

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoMagicNumberInOperationRule detects magic numbers used in arithmetic and
// binary operations. A magic number is a numeric literal (integer or
// floating-point) that is not defined as a constant and whose meaning is not
// immediately clear. Using magic numbers directly in operations makes code
// less readable and harder to maintain.
type NoMagicNumberInOperationRule struct {
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
func (r *NoMagicNumberInOperationRule) Configure(arguments lint.Arguments) error {
	r.ignoredNumbers = copyMap(defaultIgnoredNumbers)
	r.ignoredFiles = append([]string{}, defaultIgnoredFiles...)

	if len(arguments) < 1 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noMagicNumberInOperation" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		switch {
		case isRuleOption(k, "ignored-numbers"):
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf(`invalid configuration value for ignored-numbers in "noMagicNumberInOperation" rule; need string but got %T`, v)
			}
			for k, v := range parseIgnoredNumbers(s) {
				r.ignoredNumbers[k] = v
			}
		case isRuleOption(k, "ignored-files"):
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf(`invalid configuration value for ignored-files in "noMagicNumberInOperation" rule; need string but got %T`, v)
			}
			r.ignoredFiles = append(r.ignoredFiles, parseIgnoredFiles(s)...)
		}
	}

	return nil
}

// Apply applies the rule to the given file.
func (r *NoMagicNumberInOperationRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
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

	w := &lintMagicNumberInOperation{
		ignoredNumbers: r.ignoredNumbers,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoMagicNumberInOperationRule) Name() string {
	return "noMagicNumberInOperation"
}

// Group returns the rule group.
func (*NoMagicNumberInOperationRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoMagicNumberInOperationRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintMagicNumberInOperation struct {
	ignoredNumbers map[string]bool
	onFailure      func(lint.Failure)
}

func (w *lintMagicNumberInOperation) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.AssignStmt:
		// Check the right-hand side of assignments for binary expressions
		// containing magic numbers.
		for _, rhs := range n.Rhs {
			w.checkExprForBinaryOps(rhs)
		}
		return nil // don't recurse further into the assign statement
	case *ast.IfStmt:
		// Check the condition of if statements for binary operations
		// containing magic numbers.
		w.checkConditionForOps(n.Cond)
		// Continue walking into the body, but not the condition again.
		ast.Walk(w, n.Body)
		if n.Else != nil {
			ast.Walk(w, n.Else)
		}
		return nil
	case *ast.ValueSpec:
		// Check variable declarations like `var x = 5 * y`.
		for _, val := range n.Values {
			w.checkExprForBinaryOps(val)
		}
		return nil
	}
	return w
}

// checkConditionForOps looks for binary operations in condition expressions
// and reports any magic number operands. Both arithmetic and comparison
// operators count as binary operations.
func (w *lintMagicNumberInOperation) checkConditionForOps(expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		if e.Op == token.LAND || e.Op == token.LOR {
			// For logical operators, recurse into both sides to find
			// nested binary operations.
			w.checkConditionForOps(e.X)
			w.checkConditionForOps(e.Y)
		} else {
			// For arithmetic and comparison operators, check operands.
			w.checkOperand(e.X)
			w.checkOperand(e.Y)
		}
	case *ast.ParenExpr:
		w.checkConditionForOps(e.X)
	}
}

// checkExprForBinaryOps walks an expression looking for binary expressions
// and reports magic number operands. Both arithmetic and comparison operators
// count as binary operations.
func (w *lintMagicNumberInOperation) checkExprForBinaryOps(expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		if e.Op == token.LAND || e.Op == token.LOR {
			// For logical operators, recurse into both sides.
			w.checkExprForBinaryOps(e.X)
			w.checkExprForBinaryOps(e.Y)
		} else {
			// For arithmetic and comparison operators, check operands.
			w.checkOperand(e.X)
			w.checkOperand(e.Y)
		}
	case *ast.ParenExpr:
		w.checkExprForBinaryOps(e.X)
	}
}

// checkOperand checks whether an operand of a binary expression is a magic
// number literal. It also recurses into nested binary expressions and
// parenthesized expressions.
func (w *lintMagicNumberInOperation) checkOperand(expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.INT || e.Kind == token.FLOAT {
			val := normalizeNumber(e.Value)
			if !w.ignoredNumbers[val] {
				w.onFailure(lint.Failure{
					Confidence: 1,
					Category:   lint.FailureCategoryStyle,
					Failure:    fmt.Sprintf("magic number: %s, in <operation> detected", val),
					Node:       e,
				})
			}
		}
	case *ast.BinaryExpr:
		if isArithmeticOp(e.Op) {
			w.checkOperand(e.X)
			w.checkOperand(e.Y)
		}
	case *ast.ParenExpr:
		w.checkOperand(e.X)
	case *ast.UnaryExpr:
		w.checkOperand(e.X)
	}
}

func isArithmeticOp(op token.Token) bool {
	switch op {
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM,
		token.AND, token.OR, token.XOR, token.SHL, token.SHR, token.AND_NOT:
		return true
	}
	return false
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
