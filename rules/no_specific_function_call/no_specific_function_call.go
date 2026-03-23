package no_specific_function_call

import (
	"fmt"
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoSpecificFunctionCallRule reports calls to a specific function by name.
// This is primarily used as a demonstration and test of the Analysis API,
// but can also be useful for auditing purposes.
type NoSpecificFunctionCallRule struct {
	name string
}

const defaultName = "example"

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configure implements the [lint.ConfigurableRule] interface.
func (r *NoSpecificFunctionCallRule) Configure(arguments lint.Arguments) error {
	if len(arguments) < 1 {
		r.name = defaultName
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		// try direct string argument
		name, ok := arguments[0].(string)
		if !ok {
			return fmt.Errorf(`invalid argument to the "noSpecificFunctionCall" rule, expecting a k,v map or string, got %T`, arguments[0])
		}
		r.name = name
		return nil
	}

	r.name = defaultName
	for k, v := range argKV {
		if k == "name" {
			name, ok := v.(string)
			if !ok {
				return fmt.Errorf(`invalid configuration value for name in "noSpecificFunctionCall" rule; need string but got %T`, v)
			}
			r.name = name
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoSpecificFunctionCallRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if r.name == "" {
		return nil
	}

	var failures []lint.Failure

	w := &lintSpecificFunctionCall{
		name: r.name,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoSpecificFunctionCallRule) Name() string {
	return "noSpecificFunctionCall"
}

// Group returns the rule group.
func (*NoSpecificFunctionCallRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoSpecificFunctionCallRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintSpecificFunctionCall struct {
	name      string
	onFailure func(lint.Failure)
}

func (w *lintSpecificFunctionCall) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check if the called function name matches.
	// We match only simple identifiers (e.g. println), not qualified names (e.g. fmt.Println).
	ident, ok := ce.Fun.(*ast.Ident)
	if !ok {
		return w
	}

	if ident.Name != w.name {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    fmt.Sprintf("call of %s(...)", w.name),
	})

	return w
}
