package no_unchecked_error

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/lint"
)

// defaultExcludedFunctions is the set of functions whose error returns are
// commonly considered non-critical, matching errcheck's defaults.
var defaultExcludedFunctions = map[string]bool{
	// bytes.Buffer write methods
	"bytes.Buffer.Write":       true,
	"bytes.Buffer.WriteByte":   true,
	"bytes.Buffer.WriteRune":   true,
	"bytes.Buffer.WriteString": true,

	// strings.Builder write methods
	"strings.Builder.Write":       true,
	"strings.Builder.WriteByte":   true,
	"strings.Builder.WriteRune":   true,
	"strings.Builder.WriteString": true,

	// fmt print functions
	"fmt.Print":    true,
	"fmt.Printf":   true,
	"fmt.Println":  true,
	"fmt.Fprint":   true,
	"fmt.Fprintf":  true,
	"fmt.Fprintln": true,

	// rand.Read
	"math/rand.Read":   true,
	"crypto/rand.Read": true,

	// hash write methods
	"hash.Hash.Write":         true,
	"hash/maphash.Hash.Write": true,

	// io.Pipe close methods
	"io.PipeReader.CloseWithError": true,
	"io.PipeWriter.CloseWithError": true,
}

// NoUncheckedErrorRule detects when error return values from function calls
// are silently ignored, either by discarding the entire return value or by
// assigning the error to the blank identifier.
type NoUncheckedErrorRule struct {
	disableDefaultExclusions bool
	excludeFunctions         map[string]bool
}

// Configure validates and applies the rule configuration.
func (r *NoUncheckedErrorRule) Configure(arguments lint.Arguments) error {
	r.excludeFunctions = map[string]bool{}

	for _, arg := range arguments {
		cfg, ok := arg.(map[string]any)
		if !ok {
			continue
		}

		if v, ok := cfg["disable-default-exclusions"]; ok {
			if b, ok := v.(bool); ok {
				r.disableDefaultExclusions = b
			}
		}

		if v, ok := cfg["exclude-functions"]; ok {
			if list, ok := v.([]any); ok {
				for _, item := range list {
					if s, ok := item.(string); ok {
						r.excludeFunctions[s] = true
					}
				}
			}
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoUncheckedErrorRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoUncheckedError{
		pkg:                      file.Pkg,
		disableDefaultExclusions: r.disableDefaultExclusions,
		excludeFunctions:         r.excludeFunctions,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules
func (r *NoUncheckedErrorRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoUncheckedError{
		pkg:                      file.Pkg,
		disableDefaultExclusions: r.disableDefaultExclusions,
		excludeFunctions:         r.excludeFunctions,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoUncheckedErrorRule) Name() string {
	return "noUncheckedError"
}

// Group returns the rule group.
func (*NoUncheckedErrorRule) Group() string {
	return "correctness"
}

func (*NoUncheckedErrorRule) RequiresTypecheck() bool {
	return true
}

type lintNoUncheckedError struct {
	pkg                      *lint.Package
	disableDefaultExclusions bool
	excludeFunctions         map[string]bool
	onFailure                func(lint.Failure)
}

func (w *lintNoUncheckedError) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.ExprStmt:
		// Check for function calls whose error return is completely discarded
		w.checkExprStmt(n)
		return nil
	case *ast.AssignStmt:
		// Check for error values assigned to blank identifier
		w.checkAssignStmt(n)
		return w
	case *ast.GoStmt:
		// go f() — the error return is discarded
		w.checkCallExpr(n.Call, n.Call)
		return nil
	case *ast.DeferStmt:
		// defer f() — check if error return is discarded
		w.checkCallExpr(n.Call, n.Call)
		return nil
	}
	return w
}

func (w *lintNoUncheckedError) checkExprStmt(stmt *ast.ExprStmt) {
	call, ok := stmt.X.(*ast.CallExpr)
	if !ok {
		return
	}

	w.checkCallExpr(call, stmt)
}

func (w *lintNoUncheckedError) checkCallExpr(call *ast.CallExpr, reportNode ast.Node) {
	funcType := w.pkg.TypeOf(call)
	if funcType == nil {
		return
	}

	if !w.returnsError(funcType) {
		return
	}

	name := w.funcName(call)
	if w.isExcluded(name) {
		return
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryErrors,
		Confidence: 1,
		Node:       reportNode,
		Failure:    fmt.Sprintf("unchecked error in call to %v", name),
	})
}

func (w *lintNoUncheckedError) checkAssignStmt(assign *ast.AssignStmt) {
	// We only care about assignments where at least one LHS is blank identifier
	for i, lhs := range assign.Lhs {
		ident, ok := lhs.(*ast.Ident)
		if !ok || ident.Name != "_" {
			continue
		}

		if w.isErrorAtIndex(assign, i) {
			// Get the function name if this is a function call
			name := w.assignFuncName(assign)
			if w.isExcluded(name) {
				return
			}

			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryErrors,
				Confidence: 1,
				Node:       assign,
				Failure:    fmt.Sprintf("unchecked error in call to %v", name),
			})
			return // Only report once per assignment statement
		}
	}
}

func (w *lintNoUncheckedError) isErrorAtIndex(assign *ast.AssignStmt, index int) bool {
	// For tuple assignments like `a, _ := f()`
	if len(assign.Rhs) == 1 && len(assign.Lhs) > 1 {
		return w.isTupleErrorAtIndex(assign.Rhs[0], index)
	}

	// For single assignments like `_ = f.Close()`
	if index < len(assign.Rhs) {
		t := w.pkg.TypeOf(assign.Rhs[index])
		return t != nil && w.isTypeError(t)
	}

	return false
}

func (w *lintNoUncheckedError) isTupleErrorAtIndex(expr ast.Expr, index int) bool {
	t := w.pkg.TypeOf(expr)
	if t == nil {
		return false
	}

	tuple, ok := t.(*types.Tuple)
	if !ok {
		return false
	}

	if index >= tuple.Len() {
		return false
	}

	return w.isTypeError(tuple.At(index).Type())
}

func (*lintNoUncheckedError) isTypeError(t types.Type) bool {
	named, ok := t.(*types.Named)
	if ok {
		return named.Obj().Id() == "_.error"
	}
	return false
}

func (w *lintNoUncheckedError) returnsError(funcType types.Type) bool {
	switch t := funcType.(type) {
	case *types.Named:
		return w.isTypeError(t)
	default:
		retTypes, ok := funcType.Underlying().(*types.Tuple)
		if !ok {
			return false
		}
		for v := range retTypes.Variables() {
			nt, ok := v.Type().(*types.Named)
			if ok && w.isTypeError(nt) {
				return true
			}
		}
		return false
	}
}

func (w *lintNoUncheckedError) funcName(call *ast.CallExpr) string {
	fn, ok := w.getFunc(call)
	if !ok {
		return astutils.GoFmt(call.Fun)
	}

	name := fn.FullName()
	name = strings.ReplaceAll(name, "(", "")
	name = strings.ReplaceAll(name, ")", "")
	name = strings.ReplaceAll(name, "*", "")

	return name
}

func (w *lintNoUncheckedError) assignFuncName(assign *ast.AssignStmt) string {
	if len(assign.Rhs) == 1 {
		if call, ok := assign.Rhs[0].(*ast.CallExpr); ok {
			return w.funcName(call)
		}
	}
	// For multi-value RHS, try each one
	for _, rhs := range assign.Rhs {
		if call, ok := rhs.(*ast.CallExpr); ok {
			return w.funcName(call)
		}
	}
	return "<unknown>"
}

func (w *lintNoUncheckedError) isExcluded(funcName string) bool {
	// Check user-configured exclusions
	if w.excludeFunctions != nil && w.excludeFunctions[funcName] {
		return true
	}

	// Check default exclusions unless disabled
	if !w.disableDefaultExclusions && defaultExcludedFunctions[funcName] {
		return true
	}

	return false
}

func (w *lintNoUncheckedError) getFunc(call *ast.CallExpr) (*types.Func, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		// Try ident (local function call)
		if ident, ok := call.Fun.(*ast.Ident); ok {
			if obj := w.pkg.TypesInfo().ObjectOf(ident); obj != nil {
				if fn, ok := obj.(*types.Func); ok {
					return fn, true
				}
			}
		}
		return nil, false
	}

	if w.pkg.TypesInfo() == nil {
		return nil, false
	}

	fn, ok := w.pkg.TypesInfo().ObjectOf(sel.Sel).(*types.Func)
	if !ok {
		return nil, false
	}

	return fn, true
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name lowercased and without hyphens.
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}
