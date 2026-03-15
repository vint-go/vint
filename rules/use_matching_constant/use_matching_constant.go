package use_matching_constant

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseMatchingConstantRule detects string literals that match the value of an
// existing named constant but are not referencing it.
type UseMatchingConstantRule struct{}

// Apply applies the rule to given file.
func (r *UseMatchingConstantRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
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

	// Second pass: walk the AST to find string literals outside of constant
	// declarations that match a known constant value.
	var failures []lint.Failure

	w := &lintMatchingConstant{
		constsByValue: constsByValue,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

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

type lintMatchingConstant struct {
	constsByValue map[string]string
	onFailure     func(lint.Failure)
	inConstDecl   bool
}

func (w *lintMatchingConstant) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	// Skip constant declarations entirely — literals inside const blocks
	// are the definitions themselves and should not be flagged.
	if genDecl, ok := node.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
		return nil
	}

	lit, ok := node.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return w
	}

	constName, matches := w.constsByValue[lit.Value]
	if !matches {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Category:   lint.FailureCategoryLogic,
		Node:       lit,
		Failure:    fmt.Sprintf("string literal %s matches constant %s, use the constant instead", lit.Value, constName),
	})

	return w
}
