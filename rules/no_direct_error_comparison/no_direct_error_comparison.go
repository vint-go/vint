package no_direct_error_comparison

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path"
	"strings"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/lint"
)

// allowedError represents an error/function pair where direct comparison is acceptable.
type allowedError struct {
	Err string // fully qualified error name, e.g. "io.EOF"
	Fun string // fully qualified function name, e.g. "example.com/pkg.Read"
}

// NoDirectErrorComparisonRule flags direct comparisons of error values using == or != and recommends errors.Is().
type NoDirectErrorComparisonRule struct {
	allowedErrors         []allowedError
	allowedErrorsWildcard []allowedError
}

// Configure validates and applies the rule configuration.
func (r *NoDirectErrorComparisonRule) Configure(arguments lint.Arguments) error {
	r.allowedErrors = nil
	r.allowedErrorsWildcard = nil
	if len(arguments) == 0 {
		return nil
	}
	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noDirectErrorComparison" rule, expecting a k,v map, got %T`, arguments[0])
	}
	for k, v := range argKV {
		switch normalizeRuleOption(k) {
		case "allowederrors":
			entries, err := parseAllowedErrors(v)
			if err != nil {
				return fmt.Errorf(`invalid "allowed-errors" config: %w`, err)
			}
			r.allowedErrors = entries
		case "allowederrorswildcard":
			entries, err := parseAllowedErrors(v)
			if err != nil {
				return fmt.Errorf(`invalid "allowed-errors-wildcard" config: %w`, err)
			}
			r.allowedErrorsWildcard = entries
		}
	}
	return nil
}

func parseAllowedErrors(v any) ([]allowedError, error) {
	list, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("expected a list, got %T", v)
	}
	var result []allowedError
	for _, item := range list {
		m, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected a map with err/fun keys, got %T", item)
		}
		entry := allowedError{}
		if errVal, ok := m["err"]; ok {
			entry.Err, _ = errVal.(string)
		}
		if funVal, ok := m["fun"]; ok {
			entry.Fun, _ = funVal.(string)
		}
		result = append(result, entry)
	}
	return result, nil
}

// Apply applies the rule to given file.
func (r *NoDirectErrorComparisonRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintNoDirectErrorComparison{
		file:                  file,
		onFailure:             onFailure,
		allowedErrors:         r.allowedErrors,
		allowedErrorsWildcard: r.allowedErrorsWildcard,
	}
	if w.file.Pkg.TypeCheck() != nil {
		return nil
	}

	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (*NoDirectErrorComparisonRule) Name() string {
	return "noDirectErrorComparison"
}

// Group returns the rule group.
func (*NoDirectErrorComparisonRule) Group() string {
	return "correctness"
}

func (*NoDirectErrorComparisonRule) RequiresTypecheck() bool {
	return true
}

type lintNoDirectErrorComparison struct {
	file                  *lint.File
	onFailure             func(lint.Failure)
	allowedErrors         []allowedError
	allowedErrorsWildcard []allowedError
	enclosingFunc         string // current enclosing function's fully qualified name
}

func (w *lintNoDirectErrorComparison) Visit(node ast.Node) ast.Visitor {
	// Track enclosing function for allowed-errors matching.
	if fn, ok := node.(*ast.FuncDecl); ok {
		prevFunc := w.enclosingFunc
		w.enclosingFunc = w.qualifiedFuncName(fn)
		ast.Walk(&lintNoDirectErrorComparison{
			file:                  w.file,
			onFailure:             w.onFailure,
			allowedErrors:         w.allowedErrors,
			allowedErrorsWildcard: w.allowedErrorsWildcard,
			enclosingFunc:         w.enclosingFunc,
		}, fn.Body)
		w.enclosingFunc = prevFunc
		return nil // skip body, already walked
	}

	expr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return w
	}

	switch expr.Op {
	case token.EQL, token.NEQ:
	default:
		return w
	}

	// Check if either side is nil — always allowed
	if isNilIdent(expr.X) || isNilIdent(expr.Y) {
		return w
	}

	// Check if either side is io.EOF — always allowed
	if astutils.IsPkgDotName(expr.X, "io", "EOF") || astutils.IsPkgDotName(expr.Y, "io", "EOF") {
		return w
	}

	typeOfX := w.file.Pkg.TypeOf(expr.X)
	typeOfY := w.file.Pkg.TypeOf(expr.Y)

	if typeOfX == nil || typeOfY == nil {
		return w
	}

	// Check if at least one side implements the error interface
	errorIface := types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
	xIsError := types.Implements(typeOfX, errorIface)
	yIsError := types.Implements(typeOfY, errorIface)

	if !xIsError && !yIsError {
		return w
	}

	// Check allowed-errors configuration
	if w.isAllowed(expr.X, expr.Y) {
		return w
	}

	xStr := astutils.GoFmt(expr.X)
	yStr := astutils.GoFmt(expr.Y)

	var failureMsg string
	var replacementLine string
	if expr.Op == token.EQL {
		failureMsg = fmt.Sprintf("avoid direct error comparison, use errors.Is(%s, %s) instead", xStr, yStr)
		replacementLine = fmt.Sprintf("errors.Is(%s, %s)", xStr, yStr)
	} else {
		failureMsg = fmt.Sprintf("avoid direct error comparison, use !errors.Is(%s, %s) instead", xStr, yStr)
		replacementLine = fmt.Sprintf("!errors.Is(%s, %s)", xStr, yStr)
	}

	w.onFailure(lint.Failure{
		Category:        lint.FailureCategoryErrors,
		Confidence:      1,
		Node:            node,
		Failure:         failureMsg,
		ReplacementLine: replacementLine,
	})

	return w
}

// isAllowed checks whether the comparison matches any allowed-errors or allowed-errors-wildcard entry.
func (w *lintNoDirectErrorComparison) isAllowed(x, y ast.Expr) bool {
	if len(w.allowedErrors) == 0 && len(w.allowedErrorsWildcard) == 0 {
		return false
	}

	errNames := w.qualifiedErrorNames(x, y)
	if len(errNames) == 0 {
		return false
	}

	for _, errName := range errNames {
		// Check exact matches
		for _, ae := range w.allowedErrors {
			if ae.Err == errName && (ae.Fun == "" || ae.Fun == w.enclosingFunc) {
				return true
			}
		}
		// Check wildcard matches
		for _, ae := range w.allowedErrorsWildcard {
			if matchWildcard(ae.Err, errName) && (ae.Fun == "" || matchWildcard(ae.Fun, w.enclosingFunc)) {
				return true
			}
		}
	}
	return false
}

// qualifiedErrorNames extracts the fully qualified name(s) of error expressions in a comparison.
func (w *lintNoDirectErrorComparison) qualifiedErrorNames(x, y ast.Expr) []string {
	var names []string
	if name := w.qualifiedExprName(x); name != "" {
		names = append(names, name)
	}
	if name := w.qualifiedExprName(y); name != "" {
		names = append(names, name)
	}
	return names
}

// qualifiedExprName returns the fully qualified name of an expression (e.g. "io.EOF", "pkg.ErrFoo").
func (w *lintNoDirectErrorComparison) qualifiedExprName(expr ast.Expr) string {
	info := w.file.Pkg.TypesInfo()
	if info == nil {
		return ""
	}

	switch e := expr.(type) {
	case *ast.Ident:
		if obj, ok := info.Uses[e]; ok {
			return qualifiedObjName(obj)
		}
	case *ast.SelectorExpr:
		if obj, ok := info.Uses[e.Sel]; ok {
			return qualifiedObjName(obj)
		}
	}
	return ""
}

// qualifiedObjName returns the fully qualified name of a types.Object (e.g. "io.EOF").
func qualifiedObjName(obj types.Object) string {
	if obj == nil || obj.Pkg() == nil {
		return ""
	}
	return obj.Pkg().Path() + "." + obj.Name()
}

// qualifiedFuncName returns the fully qualified name of a function declaration.
func (w *lintNoDirectErrorComparison) qualifiedFuncName(fn *ast.FuncDecl) string {
	pkg := w.file.Pkg.TypesPkg()
	if pkg == nil {
		return ""
	}
	pkgPath := pkg.Path()
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		recvType := w.file.Render(fn.Recv.List[0].Type)
		recvType = strings.TrimPrefix(recvType, "*")
		return pkgPath + "." + recvType + "." + fn.Name.Name
	}
	return pkgPath + "." + fn.Name.Name
}

// matchWildcard matches a pattern against a value, supporting * as a glob wildcard.
func matchWildcard(pattern, value string) bool {
	matched, _ := path.Match(pattern, value)
	return matched
}

func isNilIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "nil"
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name lowercased and without hyphens.
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}
