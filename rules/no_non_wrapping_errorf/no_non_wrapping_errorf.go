package no_non_wrapping_errorf

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoNonWrappingErrorfRule detects fmt.Errorf calls that use non-wrapping format verbs
// (%v, %s, etc.) instead of %w for error values.
type NoNonWrappingErrorfRule struct {
	errorfMulti bool
}

const defaultErrorfMulti = true

// Configure validates and applies the rule configuration.
func (r *NoNonWrappingErrorfRule) Configure(arguments lint.Arguments) error {
	r.errorfMulti = defaultErrorfMulti
	if len(arguments) == 0 {
		return nil
	}
	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noNonWrappingErrorf" rule, expecting a k,v map, got %T`, arguments[0])
	}
	for k, v := range argKV {
		switch strings.ToLower(strings.ReplaceAll(k, "_", "-")) {
		case "errorf-multi":
			b, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for errorf-multi in "noNonWrappingErrorf" rule; need bool but got %T`, v)
			}
			r.errorfMulti = b
		}
	}
	return nil
}

// Apply applies the rule to given file.
func (r *NoNonWrappingErrorfRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()
	w := &lintNonWrappingErrorf{
		file:        file,
		errorfMulti: r.errorfMulti,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoNonWrappingErrorfRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	w := &lintNonWrappingErrorf{
		file:        file,
		errorfMulti: r.errorfMulti,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoNonWrappingErrorfRule) Name() string {
	return "noNonWrappingErrorf"
}

// Group returns the rule group.
func (*NoNonWrappingErrorfRule) Group() string {
	return "correctness"
}

// RequiresTypecheck returns true because this rule uses type information.
func (*NoNonWrappingErrorfRule) RequiresTypecheck() bool { return true }

// CacheTier returns the cache tier for this rule.
func (*NoNonWrappingErrorfRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

type lintNonWrappingErrorf struct {
	file        *lint.File
	errorfMulti bool
	onFailure   func(lint.Failure)
}

func (w *lintNonWrappingErrorf) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "fmt", "Errorf") {
		return w
	}

	if len(ce.Args) < 2 {
		return w
	}

	// First argument must be a string literal (format string)
	formatLit, ok := ce.Args[0].(*ast.BasicLit)
	if !ok {
		return w
	}

	formatStr := strings.Trim(formatLit.Value, `"`)
	formatStr = strings.Trim(formatStr, "`")

	verbs := parseFormatVerbs(formatStr)
	args := ce.Args[1:]

	// Handle variadic calls: if the last argument is ... (spread), skip analysis
	if ce.Ellipsis.IsValid() {
		return w
	}

	wCount := 0
	for i, verb := range verbs {
		if i >= len(args) {
			break
		}

		if verb == "w" {
			wCount++
			continue
		}

		// Check if this argument is an error type
		arg := args[i]
		if !w.isErrorType(arg) {
			continue
		}

		// This is an error argument with a non-wrapping verb
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryErrors,
			Node:       ce,
			Confidence: 1,
			Failure:    fmt.Sprintf("non-wrapping format verb for fmt.Errorf; use %%w to format error %s", w.file.Render(arg)),
		})
	}

	// Check for multiple %w when errorf-multi is false
	if !w.errorfMulti && wCount > 1 {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryErrors,
			Node:       ce,
			Confidence: 1,
			Failure:    "multiple %%w format verbs in fmt.Errorf call",
		})
	}

	return w
}

func (w *lintNonWrappingErrorf) isErrorType(expr ast.Expr) bool {
	typ := w.file.Pkg.TypeOf(expr)
	if typ == nil {
		return false
	}
	errorIface := types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
	return types.Implements(typ, errorIface) || types.Implements(types.NewPointer(typ), errorIface)
}

// parseFormatVerbs extracts format verbs from a format string.
// It returns the verb letter(s) for each positional argument.
// For example, "failed: %v, code: %d" returns ["v", "d"].
func parseFormatVerbs(format string) []string {
	var verbs []string
	i := 0
	for i < len(format) {
		if format[i] != '%' {
			i++
			continue
		}
		i++ // skip '%'
		if i >= len(format) {
			break
		}
		// Skip literal %%
		if format[i] == '%' {
			i++
			continue
		}
		// Skip flags: #, 0, -, +, ' ', etc.
		for i < len(format) && isFlag(format[i]) {
			i++
		}
		// Skip width: [0-9]* or *
		if i < len(format) && format[i] == '*' {
			i++
			// * consumes an argument for width, add a placeholder
			verbs = append(verbs, "*")
		} else {
			for i < len(format) && format[i] >= '0' && format[i] <= '9' {
				i++
			}
		}
		// Skip precision: . [0-9]* or .*
		if i < len(format) && format[i] == '.' {
			i++
			if i < len(format) && format[i] == '*' {
				i++
				// * consumes an argument for precision, add a placeholder
				verbs = append(verbs, "*")
			} else {
				for i < len(format) && format[i] >= '0' && format[i] <= '9' {
					i++
				}
			}
		}
		// The verb character
		if i < len(format) {
			verbs = append(verbs, string(format[i]))
			i++
		}
	}
	return verbs
}

func isFlag(c byte) bool {
	return c == '#' || c == '0' || c == '-' || c == '+' || c == ' '
}
