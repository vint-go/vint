package no_builtin_shadow_decl

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// predeclaredIdentifiers contains Go's predeclared identifiers: built-in functions,
// types, constants, and variables that should not be shadowed via top-level declarations.
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

// NoBuiltinShadowDeclRule detects top-level declarations that shadow predeclared identifiers.
type NoBuiltinShadowDeclRule struct{}

// Apply applies the rule to given file.
func (r *NoBuiltinShadowDeclRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			// Only check top-level functions, not methods (methods have a receiver).
			if d.Recv != nil {
				continue
			}
			if predeclaredIdentifiers[d.Name.Name] {
				failures = append(failures, lint.Failure{
					Confidence: 1,
					Node:       d.Name,
					Category:   lint.FailureCategoryLogic,
					Failure:    fmt.Sprintf("top-level declaration shadows predeclared identifier %s", d.Name.Name),
				})
			}

		case *ast.GenDecl:
			switch d.Tok {
			case token.VAR, token.CONST:
				for _, spec := range d.Specs {
					vs, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					for _, name := range vs.Names {
						if predeclaredIdentifiers[name.Name] {
							failures = append(failures, lint.Failure{
								Confidence: 1,
								Node:       name,
								Category:   lint.FailureCategoryLogic,
								Failure:    fmt.Sprintf("top-level declaration shadows predeclared identifier %s", name.Name),
							})
						}
					}
				}
			case token.TYPE:
				for _, spec := range d.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					if predeclaredIdentifiers[ts.Name.Name] {
						failures = append(failures, lint.Failure{
							Confidence: 1,
							Node:       ts.Name,
							Category:   lint.FailureCategoryLogic,
							Failure:    fmt.Sprintf("top-level declaration shadows predeclared identifier %s", ts.Name.Name),
						})
					}
				}
			}
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoBuiltinShadowDeclRule) Name() string {
	return "noBuiltinShadowDecl"
}

// Group returns the rule group.
func (*NoBuiltinShadowDeclRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoBuiltinShadowDeclRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
