package use_named_result

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseNamedResultRule detects unnamed results that may benefit from names.
// It identifies function return values lacking names that could benefit
// from being named to improve code clarity.
type UseNamedResultRule struct {
	checkExported bool
}

// Configure validates and applies the rule configuration.
func (r *UseNamedResultRule) Configure(arguments lint.Arguments) error {
	r.checkExported = false

	if len(arguments) == 0 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "useNamedResult" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		if normalizeOption(k) == "checkexported" {
			b, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for checkExported in "useNamedResult" rule; need bool but got %T`, v)
			}
			r.checkExported = b
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *UseNamedResultRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNamedResult{
		checkExported: r.checkExported,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseNamedResultRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNamedResult{
		checkExported: r.checkExported,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*UseNamedResultRule) Name() string {
	return "useNamedResult"
}

// Group returns the rule group.
func (*UseNamedResultRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseNamedResultRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNamedResult struct {
	checkExported bool
	onFailure     func(lint.Failure)
}

func (w *lintNamedResult) Visit(node ast.Node) ast.Visitor {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return w
	}

	// Skip exported functions if checkExported is false.
	if !w.checkExported && funcDecl.Name.IsExported() {
		return w
	}

	results := funcDecl.Type.Results
	if results == nil {
		return w
	}

	// Collect individual return types, expanding grouped fields.
	// Only consider unnamed results.
	types := collectResultTypes(results)
	if types == nil {
		// All results are already named.
		return w
	}

	n := len(types)
	if n < 2 {
		return w
	}

	if n == 2 {
		// For 2-return functions: warn when both returns share the same type
		// (ambiguous which is which), unless the second is error or bool
		// (conventionally obvious).
		if types[0] != types[1] {
			// Different types: no warning (types disambiguate).
			return w
		}
		// Same type: warn unless second is error or bool.
		if types[1] == "error" || types[1] == "bool" {
			return w
		}
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryStyle,
			Confidence: 1,
			Node:       funcDecl.Type.Results,
			Failure:    fmt.Sprintf("unnamed results of type '%s' may benefit from named results", strings.Join(types, "', '")),
		})
		return w
	}

	// For >2 returns: warn when duplicate types appear, excluding trailing error or bool.
	typesToCheck := types
	// Exclude trailing error or bool.
	for len(typesToCheck) > 0 {
		last := typesToCheck[len(typesToCheck)-1]
		if last == "error" || last == "bool" {
			typesToCheck = typesToCheck[:len(typesToCheck)-1]
		} else {
			break
		}
	}

	if hasDuplicates(typesToCheck) {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryStyle,
			Confidence: 1,
			Node:       funcDecl.Type.Results,
			Failure:    fmt.Sprintf("unnamed results of type '%s' may benefit from named results", strings.Join(types, "', '")),
		})
	}

	return w
}

// collectResultTypes returns the list of type strings for function results.
// Returns nil if any result is already named.
func collectResultTypes(results *ast.FieldList) []string {
	var types []string
	for _, field := range results.List {
		if len(field.Names) > 0 {
			// Already named; skip this function.
			return nil
		}
		typeStr := typeToString(field.Type)
		types = append(types, typeStr)
	}
	return types
}

// typeToString returns a string representation of an AST type expression.
func typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return typeToString(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + typeToString(t.X)
	case *ast.ArrayType:
		if t.Len == nil {
			return "[]" + typeToString(t.Elt)
		}
		return "[...]" + typeToString(t.Elt)
	case *ast.MapType:
		return "map[" + typeToString(t.Key) + "]" + typeToString(t.Value)
	case *ast.ChanType:
		switch t.Dir {
		case ast.SEND:
			return "chan<- " + typeToString(t.Value)
		case ast.RECV:
			return "<-chan " + typeToString(t.Value)
		default:
			return "chan " + typeToString(t.Value)
		}
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.FuncType:
		return "func(...)"
	case *ast.Ellipsis:
		return "..." + typeToString(t.Elt)
	default:
		return fmt.Sprintf("%T", expr)
	}
}

// hasDuplicates returns true if there are any duplicate strings in the slice.
func hasDuplicates(strs []string) bool {
	seen := make(map[string]bool, len(strs))
	for _, s := range strs {
		if seen[s] {
			return true
		}
		seen[s] = true
	}
	return false
}

// normalizeOption normalizes a configuration option name by removing hyphens,
// underscores, and lowering case.
func normalizeOption(name string) string {
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "_", "")
	return strings.ToLower(name)
}
