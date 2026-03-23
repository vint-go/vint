package no_unused_parameter

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoUnusedParameterRule reports function parameters that are completely unused
// within the function body. By default, only unexported (private) functions are
// checked. Set check-exported to true to also analyze exported functions.
type NoUnusedParameterRule struct {
	checkExported bool
}

// Configure validates and applies the rule configuration.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *NoUnusedParameterRule) Configure(arguments lint.Arguments) error {
	if len(arguments) < 1 {
		r.checkExported = false
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noUnusedParameter" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		if isRuleOption(k, "check-exported") {
			val, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for check-exported in "noUnusedParameter" rule; need bool but got %T`, v)
			}
			r.checkExported = val
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoUnusedParameterRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	w := &lintUnusedParam{
		checkExported: r.checkExported,
		onFailure:     onFailure,
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoUnusedParameterRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintUnusedParam{
		checkExported: r.checkExported,
		onFailure:     onFailure,
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoUnusedParameterRule) Name() string {
	return "noUnusedParameter"
}

// Group returns the rule group.
func (*NoUnusedParameterRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnusedParameterRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintUnusedParam struct {
	checkExported bool
	onFailure     func(lint.Failure)
}

func (w *lintUnusedParam) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		if n.Body == nil {
			return w // skip function prototypes (e.g. CGo, assembly)
		}

		// Skip exported functions unless check-exported is enabled
		if !w.checkExported && n.Name != nil && ast.IsExported(n.Name.Name) {
			return w
		}

		w.checkFunc(n.Type.Params, n.Body)
	case *ast.FuncLit:
		w.checkFunc(n.Type.Params, n.Body)
	}

	return w
}

func (w *lintUnusedParam) checkFunc(params *ast.FieldList, body *ast.BlockStmt) {
	if params == nil || len(params.List) == 0 || body == nil {
		return
	}

	// Collect all named, non-blank parameters
	//nolint:staticcheck // TODO: ast.Object is deprecated
	namedParams := map[*ast.Object]bool{}
	for _, field := range params.List {
		for _, name := range field.Names {
			if name.Name == "_" {
				continue
			}
			namedParams[name.Obj] = true // true means unused
		}
	}

	if len(namedParams) == 0 {
		return
	}

	// Walk the function body looking for references to the parameters
	ast.Inspect(body, func(n ast.Node) bool {
		ident, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		if _, isParam := namedParams[ident.Obj]; isParam {
			namedParams[ident.Obj] = false // mark as used
		}
		return true
	})

	// Report any parameters that remain unused
	for _, field := range params.List {
		for _, name := range field.Names {
			if name.Name == "_" {
				continue
			}
			if namedParams[name.Obj] {
				w.onFailure(lint.Failure{
					Confidence: 1,
					Node:       name,
					Category:   lint.FailureCategoryBadPractice,
					Failure:    fmt.Sprintf("parameter '%s' seems to be unused, consider removing or renaming it as _", name.Name),
				})
			}
		}
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
