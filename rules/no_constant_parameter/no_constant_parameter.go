package no_constant_parameter

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoConstantParameterRule reports function parameters that always receive the
// same constant value at every call site in the file. By default, only
// unexported (private) functions are checked. Set check-exported to true to
// also analyze exported functions.
type NoConstantParameterRule struct {
	checkExported bool
}

// Configure validates and applies the rule configuration.
//
// Configuration implements the [lint.ConfigurableRule] interface.
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

// paramInfo tracks a single named parameter of a function declaration.
type paramInfo struct {
	field     *ast.Field
	nameIdent *ast.Ident
	paramName string
	index     int // positional index among all individual parameter names
}

// funcInfo holds the function declaration and its parameter metadata.
type funcInfo struct {
	decl   *ast.FuncDecl
	params []paramInfo
}

// Apply applies the rule to given file.
func (r *NoConstantParameterRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	// Phase 1: Collect all function declarations to analyze.
	funcs := map[string]*funcInfo{}
	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if funcDecl.Body == nil {
			continue // skip prototypes
		}
		// Skip methods (receiver != nil) since matching call sites is more
		// complex and unreliable without type info.
		if funcDecl.Recv != nil {
			continue
		}
		// Skip exported functions unless check-exported is enabled
		if !r.checkExported && funcDecl.Name != nil && ast.IsExported(funcDecl.Name.Name) {
			continue
		}
		params := collectParams(funcDecl)
		if len(params) == 0 {
			continue
		}
		funcs[funcDecl.Name.Name] = &funcInfo{
			decl:   funcDecl,
			params: params,
		}
	}

	if len(funcs) == 0 {
		return nil
	}

	// Phase 2: Walk the entire file and collect call-site argument values.
	// For each function parameter, we track:
	// - the set of constant string representations seen
	// - the total number of calls
	type paramCallInfo struct {
		constValues map[string]bool // set of rendered constant values
		callCount   int             // total number of call sites
	}
	callData := map[string][]paramCallInfo{} // funcName -> per-param call info

	for name, fi := range funcs {
		pci := make([]paramCallInfo, len(fi.params))
		for i := range pci {
			pci[i].constValues = map[string]bool{}
		}
		callData[name] = pci
	}

	ast.Inspect(file.AST, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		ident, ok := call.Fun.(*ast.Ident)
		if !ok {
			return true
		}
		fi, exists := funcs[ident.Name]
		if !exists {
			return true
		}
		pci := callData[ident.Name]

		// Match arguments to parameter positions.
		// If the call has variadic args or different arg count, skip.
		if len(call.Args) != len(fi.params) {
			return true
		}
		if call.Ellipsis.IsValid() {
			return true // skip calls with ... expansion
		}

		for i, arg := range call.Args {
			if i >= len(pci) {
				break
			}
			pci[i].callCount++
			if isConstantExpr(arg) {
				rendered := astutils.GoFmt(arg)
				pci[i].constValues[rendered] = true
			} else {
				// Non-constant argument: mark with a sentinel so we know
				// the parameter is not always constant.
				pci[i].constValues[""] = true
				// Add a second sentinel to ensure len > 1 or use a flag
				pci[i].constValues["\x00non-const"] = true
			}
		}

		return true
	})

	// Phase 3: Report parameters that always receive the same constant.
	for name, fi := range funcs {
		pci := callData[name]
		for i, pi := range fi.params {
			if i >= len(pci) {
				break
			}
			info := pci[i]
			// Must have at least 1 call site, exactly 1 unique constant value,
			// and no non-constant arguments.
			if info.callCount == 0 {
				continue
			}
			if len(info.constValues) != 1 {
				continue
			}
			// Get the single constant value
			var constVal string
			for v := range info.constValues {
				constVal = v
			}
			// Skip if it's a non-constant sentinel
			if constVal == "" || constVal == "\x00non-const" {
				continue
			}

			failures = append(failures, lint.Failure{
				Confidence: 1,
				Node:       pi.nameIdent,
				Category:   lint.FailureCategoryBadPractice,
				Failure:    fmt.Sprintf("%s always receives %s", pi.paramName, constVal),
			})
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

// CacheTier returns the cache tier for this rule.
func (*NoConstantParameterRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// collectParams extracts a flat list of named parameters from a function.
func collectParams(fn *ast.FuncDecl) []paramInfo {
	if fn.Type.Params == nil {
		return nil
	}
	var result []paramInfo
	idx := 0
	for _, field := range fn.Type.Params.List {
		if len(field.Names) == 0 {
			idx++
			continue
		}
		for _, name := range field.Names {
			if name.Name == "_" {
				idx++
				continue
			}
			result = append(result, paramInfo{
				field:     field,
				nameIdent: name,
				paramName: name.Name,
				index:     idx,
			})
			idx++
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
		// true, false, nil are constant identifiers
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
