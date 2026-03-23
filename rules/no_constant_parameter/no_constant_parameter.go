package no_constant_parameter

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
	"sync"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/lint"
)

// Compile-time interface checks.
var (
	_ lint.Rule             = (*NoConstantParameterRule)(nil)
	_ lint.AggregatingRule  = (*NoConstantParameterRule)(nil)
	_ lint.ConfigurableRule = (*NoConstantParameterRule)(nil)
	_ lint.Grouped          = (*NoConstantParameterRule)(nil)
)

// NoConstantParameterRule reports function parameters that always receive the
// same constant value at every call site across all files in a package.
// By default, only unexported (private) functions are checked. Set
// check-exported to true to also analyze exported functions.
type NoConstantParameterRule struct {
	checkExported bool
	mu            sync.Mutex
	packages      map[string]*pkgData // pkgDir → collected data
}

// pkgData holds collected function declarations and call sites for a package.
type pkgData struct {
	funcs map[string]*collectedFunc     // funcName → declaration info
	calls map[string][]collectedCall    // funcName → call sites from all files
}

// collectedFunc stores a function declaration's parameter metadata.
type collectedFunc struct {
	paramCount int
	params     []collectedParam
}

// collectedParam stores a single parameter's name and source position.
type collectedParam struct {
	name     string
	startPos token.Position
	endPos   token.Position
}

// collectedCall stores one call site's argument data.
type collectedCall struct {
	argCount    int
	hasEllipsis bool
	args        []collectedArg
}

// collectedArg stores whether an argument is constant and its rendered form.
type collectedArg struct {
	isConstant bool
	rendered   string
}

// Configure validates and applies the rule configuration.
func (r *NoConstantParameterRule) Configure(arguments lint.Arguments) error {
	if len(arguments) < 1 {
		r.checkExported = false
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noConstantParameter" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		if isRuleOption(k, "check-exported") {
			val, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for check-exported in "noConstantParameter" rule; need bool but got %T`, v)
			}
			r.checkExported = val
		}
	}

	return nil
}

// Apply returns nil — this rule produces results via Collect/Finalize.
func (r *NoConstantParameterRule) Apply(_ *lint.File, _ lint.Arguments) []lint.Failure {
	return nil
}

// Collect gathers function declarations and call sites from a single file.
// Safe for concurrent calls.
func (r *NoConstantParameterRule) Collect(file *lint.File, _ lint.Arguments) {
	// Phase 1: Walk AST without lock to gather local data.
	localFuncs := map[string]*collectedFunc{}
	localCalls := map[string][]collectedCall{}

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}
		// Skip methods (receiver != nil) since matching call sites is more
		// complex and unreliable without type info.
		if funcDecl.Recv != nil {
			continue
		}
		// Skip exported functions unless check-exported is enabled.
		if !r.checkExported && funcDecl.Name != nil && ast.IsExported(funcDecl.Name.Name) {
			continue
		}
		params := collectParams(funcDecl, file)
		if len(params) == 0 {
			continue
		}
		localFuncs[funcDecl.Name.Name] = &collectedFunc{
			paramCount: len(params),
			params:     params,
		}
	}

	// Collect call sites from the entire file.
	ast.Inspect(file.AST, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		ident, ok := call.Fun.(*ast.Ident)
		if !ok {
			return true
		}

		args := make([]collectedArg, len(call.Args))
		for i, arg := range call.Args {
			if isConstantExpr(arg) {
				args[i] = collectedArg{
					isConstant: true,
					rendered:   astutils.GoFmt(arg),
				}
			}
		}
		localCalls[ident.Name] = append(localCalls[ident.Name], collectedCall{
			argCount:    len(call.Args),
			hasEllipsis: call.Ellipsis.IsValid(),
			args:        args,
		})

		return true
	})

	if len(localFuncs) == 0 && len(localCalls) == 0 {
		return
	}

	// Phase 2: Merge into shared state under lock.
	pkgDir := filepath.Dir(file.Name)
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.packages == nil {
		r.packages = map[string]*pkgData{}
	}
	pkg := r.packages[pkgDir]
	if pkg == nil {
		pkg = &pkgData{
			funcs: map[string]*collectedFunc{},
			calls: map[string][]collectedCall{},
		}
		r.packages[pkgDir] = pkg
	}

	for name, fi := range localFuncs {
		pkg.funcs[name] = fi
	}
	for name, calls := range localCalls {
		pkg.calls[name] = append(pkg.calls[name], calls...)
	}
}

// Finalize analyzes collected data across all files and returns failures.
func (r *NoConstantParameterRule) Finalize() []lint.Failure {
	var failures []lint.Failure

	for _, pkg := range r.packages {
		for funcName, fi := range pkg.funcs {
			calls := pkg.calls[funcName]
			if len(calls) == 0 {
				continue
			}

			// Aggregate per-param constancy across all call sites.
			type paramAgg struct {
				constValues map[string]bool
				callCount   int
			}
			agg := make([]paramAgg, fi.paramCount)
			for i := range agg {
				agg[i].constValues = map[string]bool{}
			}

			for _, call := range calls {
				if call.argCount != fi.paramCount {
					continue
				}
				if call.hasEllipsis {
					continue
				}
				for i, arg := range call.args {
					if i >= fi.paramCount {
						break
					}
					agg[i].callCount++
					if arg.isConstant {
						agg[i].constValues[arg.rendered] = true
					} else {
						agg[i].constValues[""] = true
						agg[i].constValues["\x00non-const"] = true
					}
				}
			}

			for i, pi := range fi.params {
				if i >= len(agg) {
					break
				}
				info := agg[i]
				if info.callCount == 0 || len(info.constValues) != 1 {
					continue
				}
				var constVal string
				for v := range info.constValues {
					constVal = v
				}
				if constVal == "" || constVal == "\x00non-const" {
					continue
				}

				failures = append(failures, lint.Failure{
					Confidence: 1,
					Category:   lint.FailureCategoryBadPractice,
					Failure:    fmt.Sprintf("%s always receives %s", pi.name, constVal),
					Position: lint.FailurePosition{
						Start: pi.startPos,
						End:   pi.endPos,
					},
				})
			}
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoConstantParameterRule) Name() string {
	return "noConstantParameter"
}

// Group returns the rule group.
func (*NoConstantParameterRule) Group() string {
	return "suspicious"
}

// collectParams extracts a flat list of named parameters from a function,
// capturing their source positions for later failure reporting.
func collectParams(fn *ast.FuncDecl, file *lint.File) []collectedParam {
	if fn.Type.Params == nil {
		return nil
	}
	var result []collectedParam
	for _, field := range fn.Type.Params.List {
		if len(field.Names) == 0 {
			continue
		}
		for _, name := range field.Names {
			if name.Name == "_" {
				continue
			}
			result = append(result, collectedParam{
				name:     name.Name,
				startPos: file.ToPosition(name.Pos()),
				endPos:   file.ToPosition(name.End()),
			})
		}
	}
	return result
}

// isConstantExpr returns true if the expression is a compile-time constant
// literal (basic lit, unary of basic lit, or a composite of those).
func isConstantExpr(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return true
	case *ast.Ident:
		return e.Name == "true" || e.Name == "false" || e.Name == "nil"
	case *ast.UnaryExpr:
		return isConstantExpr(e.X)
	case *ast.ParenExpr:
		return isConstantExpr(e.X)
	default:
		return false
	}
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name lowercased and without hyphens.
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}
