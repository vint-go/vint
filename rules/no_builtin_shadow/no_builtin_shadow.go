package no_builtin_shadow

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// predeclaredIdentifiers contains Go's predeclared identifiers: built-in functions,
// types, constants, and variables that should not be shadowed via assignments.
var predeclaredIdentifiers = map[string]bool{
	// Built-in functions
	"append":  true,
	"cap":     true,
	"clear":   true,
	"close":   true,
	"complex": true,
	"copy":    true,
	"delete":  true,
	"imag":    true,
	"len":     true,
	"make":    true,
	"max":     true,
	"min":     true,
	"new":     true,
	"panic":   true,
	"print":   true,
	"println": true,
	"real":    true,
	"recover": true,

	// Built-in types
	"any":        true,
	"bool":       true,
	"byte":       true,
	"comparable": true,
	"complex64":  true,
	"complex128": true,
	"error":      true,
	"float32":    true,
	"float64":    true,
	"int":        true,
	"int8":       true,
	"int16":      true,
	"int32":      true,
	"int64":      true,
	"rune":       true,
	"string":     true,
	"uint":       true,
	"uint8":      true,
	"uint16":     true,
	"uint32":     true,
	"uint64":     true,
	"uintptr":    true,

	// Built-in constants
	"true":  true,
	"false": true,
	"iota":  true,

	// Zero value
	"nil": true,
}

// NoBuiltinShadowRule detects when predeclared identifiers are shadowed in assignments.
type NoBuiltinShadowRule struct{}

// Apply applies the rule to given file.
func (r *NoBuiltinShadowRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintBuiltinShadow{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoBuiltinShadowRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintBuiltinShadow{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoBuiltinShadowRule) Name() string {
	return "noBuiltinShadow"
}

// Group returns the rule group.
func (*NoBuiltinShadowRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoBuiltinShadowRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintBuiltinShadow struct {
	onFailure func(lint.Failure)
}

func (w *lintBuiltinShadow) Visit(node ast.Node) ast.Visitor {
	assign, ok := node.(*ast.AssignStmt)
	if !ok {
		return w
	}

	// Only check short variable declarations (:=) and plain assignments (=)
	if assign.Tok != token.DEFINE && assign.Tok != token.ASSIGN {
		return w
	}

	for _, lhs := range assign.Lhs {
		ident, ok := lhs.(*ast.Ident)
		if !ok {
			continue
		}

		if predeclaredIdentifiers[ident.Name] {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       ident,
				Category:   lint.FailureCategoryLogic,
				Failure:    fmt.Sprintf("assignment shadows predeclared identifier %s", ident.Name),
			})
		}
	}

	return w
}
