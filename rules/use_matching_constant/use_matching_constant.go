package use_matching_constant

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

const defaultMinOccurrences = 3

// UseMatchingConstantRule detects string literals that match the value of an
// existing named constant but are not referencing it.
type UseMatchingConstantRule struct {
	minOccurrences int
}

// Configure validates and applies the rule configuration.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *UseMatchingConstantRule) Configure(arguments lint.Arguments) error {
	r.minOccurrences = defaultMinOccurrences

	if len(arguments) < 1 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "useMatchingConstant" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		switch {
		case isRuleOption(k, "min-occurrences"):
			n, ok := lint.ToInt64(v)
			if !ok {
				return fmt.Errorf(`invalid configuration value for min-occurrences in "useMatchingConstant" rule; need integer but got %T`, v)
			}
			r.minOccurrences = int(n)
		default:
			return fmt.Errorf(`unknown option %q for "useMatchingConstant" rule`, k)
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *UseMatchingConstantRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	minOcc := r.minOccurrences
	if minOcc == 0 {
		minOcc = defaultMinOccurrences
	}

	// First pass: collect all named string constants and their values.
	// Maps a string literal value (including quotes, e.g. `"active"`) to the constant name.
	constsByValue := map[string]string{}

	for _, decl := range file.AST.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}

		for _, spec := range genDecl.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for i, val := range vs.Values {
				if i >= len(vs.Names) {
					break
				}

				lit, ok := val.(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					continue
				}

				constName := vs.Names[i].Name
				// Store the first constant found for each value.
				if _, exists := constsByValue[lit.Value]; !exists {
					constsByValue[lit.Value] = constName
				}
			}
		}
	}

	if len(constsByValue) == 0 {
		return nil
	}

	// Second pass: walk the AST to count occurrences and collect positions
	// of string literals outside constant declarations that match a known
	// constant value. Skip map keys in composite literals.
	type occurrence struct {
		node *ast.BasicLit
	}
	// occurrences maps constant value -> list of matching literal nodes
	occurrences := map[string][]occurrence{}

	counter := &lintMatchingConstantCounter{
		constsByValue: constsByValue,
		onMatch: func(lit *ast.BasicLit) {
			occurrences[lit.Value] = append(occurrences[lit.Value], occurrence{node: lit})
		},
	}
	ast.Walk(counter, file.AST)

	// Third pass: only report if occurrences >= threshold.
	var failures []lint.Failure
	for value, occs := range occurrences {
		if len(occs) < minOcc {
			continue
		}
		constName := constsByValue[value]
		for _, occ := range occs {
			failures = append(failures, lint.Failure{
				Confidence: 1,
				Category:   lint.FailureCategoryLogic,
				Node:       occ.node,
				Failure:    fmt.Sprintf("string literal %s matches constant %s, use the constant instead", occ.node.Value, constName),
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*UseMatchingConstantRule) Name() string {
	return "useMatchingConstant"
}

// Group returns the rule group.
func (*UseMatchingConstantRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*UseMatchingConstantRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// lintMatchingConstantCounter walks the AST counting string literal
// occurrences that match known constants, skipping constant declarations
// and map keys in composite literals.
type lintMatchingConstantCounter struct {
	constsByValue map[string]string
	onMatch       func(*ast.BasicLit)
}

func (w *lintMatchingConstantCounter) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	// Skip constant declarations entirely — literals inside const blocks
	// are the definitions themselves and should not be counted.
	if genDecl, ok := node.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
		return nil
	}

	// For KeyValueExpr inside composite literals, skip the Key but walk the Value.
	// This avoids counting map keys as occurrences.
	if kv, ok := node.(*ast.KeyValueExpr); ok {
		// Only walk the value side, not the key.
		ast.Walk(w, kv.Value)
		return nil
	}

	lit, ok := node.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return w
	}

	if _, matches := w.constsByValue[lit.Value]; matches {
		w.onMatch(lit)
	}

	return w
}

func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}
