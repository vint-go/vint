package no_magic_number_in_argument

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoMagicNumberInArgumentRule detects magic numbers used as function or method
// call arguments. A magic number is a numeric literal (integer or float) that
// is not defined as a constant and whose meaning is not immediately clear.
type NoMagicNumberInArgumentRule struct {
	ignoredNumbers   map[string]bool
	ignoredFunctions []functionPattern
	ignoredFiles     []string
}

type functionPattern struct {
	pkg  string // package or "*" for wildcard
	name string // function name or "*" for wildcard
}

var defaultIgnoredNumbers = map[string]bool{
	"0":   true,
	"0.0": true,
	"1":   true,
	"1.0": true,
}

var defaultIgnoredFunctions = []functionPattern{
	{pkg: "math", name: "*"},
	{pkg: "http", name: "StatusText"},
	{pkg: "strconv", name: "ParseInt"},
	{pkg: "strconv", name: "ParseUint"},
	{pkg: "strconv", name: "ParseFloat"},
	{pkg: "strconv", name: "FormatInt"},
	{pkg: "strconv", name: "FormatUint"},
	{pkg: "strconv", name: "FormatFloat"},
	{pkg: "time", name: "Date"},
}

var defaultIgnoredFiles = []string{"_test.go"}

// Configure validates and applies the rule configuration.
//
// Configure implements the [lint.ConfigurableRule] interface.
func (r *NoMagicNumberInArgumentRule) Configure(arguments lint.Arguments) error {
	r.ignoredNumbers = copyMap(defaultIgnoredNumbers)
	r.ignoredFunctions = append([]functionPattern{}, defaultIgnoredFunctions...)
	r.ignoredFiles = append([]string{}, defaultIgnoredFiles...)

	if len(arguments) < 1 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noMagicNumberInArgument" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		switch {
		case isRuleOption(k, "ignored-numbers"):
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf(`invalid configuration value for ignored-numbers in "noMagicNumberInArgument" rule; need string but got %T`, v)
			}
			for k, v := range parseIgnoredNumbers(s) {
				r.ignoredNumbers[k] = v
			}
		case isRuleOption(k, "ignored-functions"):
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf(`invalid configuration value for ignored-functions in "noMagicNumberInArgument" rule; need string but got %T`, v)
			}
			r.ignoredFunctions = append(r.ignoredFunctions, parseIgnoredFunctions(s)...)
		case isRuleOption(k, "ignored-files"):
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf(`invalid configuration value for ignored-files in "noMagicNumberInArgument" rule; need string but got %T`, v)
			}
			r.ignoredFiles = append(r.ignoredFiles, parseIgnoredFiles(s)...)
		}
	}

	return nil
}

// Apply applies the rule to the given file.
func (r *NoMagicNumberInArgumentRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if r.ignoredNumbers == nil {
		r.ignoredNumbers = copyMap(defaultIgnoredNumbers)
		r.ignoredFunctions = append([]functionPattern{}, defaultIgnoredFunctions...)
		r.ignoredFiles = append([]string{}, defaultIgnoredFiles...)
	}

	// Check if the file should be ignored.
	for _, pattern := range r.ignoredFiles {
		if strings.HasSuffix(file.Name, pattern) {
			return nil
		}
	}

	var failures []lint.Failure

	w := &lintMagicNumberInArgument{
		ignoredNumbers:   r.ignoredNumbers,
		ignoredFunctions: r.ignoredFunctions,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoMagicNumberInArgumentRule) Name() string {
	return "noMagicNumberInArgument"
}

// Group returns the rule group.
func (*NoMagicNumberInArgumentRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoMagicNumberInArgumentRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintMagicNumberInArgument struct {
	ignoredNumbers   map[string]bool
	ignoredFunctions []functionPattern
	onFailure        func(lint.Failure)
}

func (w *lintMagicNumberInArgument) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check if this function call is in the ignored list.
	if w.isIgnoredFunction(call) {
		return w
	}

	// Check each argument for magic numbers.
	for _, arg := range call.Args {
		w.checkExprForMagicNumber(arg)
	}

	return w
}

func (w *lintMagicNumberInArgument) checkExprForMagicNumber(expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.INT || e.Kind == token.FLOAT {
			val := normalizeNumber(e.Value)
			if !w.ignoredNumbers[val] {
				w.onFailure(lint.Failure{
					Confidence: 1,
					Category:   lint.FailureCategoryStyle,
					Failure:    fmt.Sprintf("magic number: %s, in <argument> detected", val),
					Node:       e,
				})
			}
		}
	case *ast.BinaryExpr:
		// Check both sides of binary expressions like 5*time.Second
		w.checkExprForMagicNumber(e.X)
		w.checkExprForMagicNumber(e.Y)
	case *ast.UnaryExpr:
		// Check unary expressions like -5
		w.checkExprForMagicNumber(e.X)
	case *ast.ParenExpr:
		// Check parenthesized expressions like (5)
		w.checkExprForMagicNumber(e.X)
	}
}

func (w *lintMagicNumberInArgument) isIgnoredFunction(call *ast.CallExpr) bool {
	switch fun := call.Fun.(type) {
	case *ast.SelectorExpr:
		// Qualified call like pkg.Func()
		ident, ok := fun.X.(*ast.Ident)
		if !ok {
			return false
		}
		pkgName := ident.Name
		funcName := fun.Sel.Name
		return w.matchesIgnoredFunction(pkgName, funcName)
	case *ast.Ident:
		// Built-in or local function call - check with empty package
		return w.matchesIgnoredFunction("", fun.Name)
	}
	return false
}

func (w *lintMagicNumberInArgument) matchesIgnoredFunction(pkg, name string) bool {
	for _, pattern := range w.ignoredFunctions {
		if pattern.pkg == pkg || (pattern.pkg != "" && pattern.name == "*" && pattern.pkg == pkg) {
			if pattern.name == "*" || pattern.name == name {
				return true
			}
		}
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
			// Use %g to get a clean representation.
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

func parseIgnoredFunctions(s string) []functionPattern {
	var result []functionPattern
	for _, fn := range strings.Split(s, ",") {
		fn = strings.TrimSpace(fn)
		if fn == "" {
			continue
		}
		// Split on last dot to separate package from function.
		if idx := strings.LastIndex(fn, "."); idx != -1 {
			pkg := fn[:idx]
			name := fn[idx+1:]
			result = append(result, functionPattern{pkg: pkg, name: name})
		} else {
			// No dot means just a function name.
			result = append(result, functionPattern{pkg: "", name: fn})
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

