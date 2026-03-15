package no_duplicate_constants

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoDuplicateConstantsRule detects named constants that have the same value
// as another constant, indicating potential duplication or a missed opportunity
// for consolidation.
type NoDuplicateConstantsRule struct{}

// Apply applies the rule to given file.
func (r *NoDuplicateConstantsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	// constantEntry tracks a named constant's name and the node where it was declared.
	type constantEntry struct {
		name string
		node ast.Node
	}

	// Map from literal value string to the first constant that declared it.
	seen := map[string]constantEntry{}

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

			// Skip iota-based constants (no explicit values or iota expressions).
			if len(vs.Values) == 0 {
				continue
			}

			for i, val := range vs.Values {
				if i >= len(vs.Names) {
					break
				}

				constName := vs.Names[i].Name

				// Only consider basic literal values (strings, ints, floats).
				lit, ok := val.(*ast.BasicLit)
				if !ok {
					continue
				}

				key := litKey(lit)

				if prev, exists := seen[key]; exists {
					failures = append(failures, lint.Failure{
						Confidence: 1,
						Category:   lint.FailureCategoryLogic,
						Failure:    fmt.Sprintf("duplicate constant value %s: %s has the same value as %s", lit.Value, constName, prev.name),
						Node:       vs,
					})
				} else {
					seen[key] = constantEntry{name: constName, node: vs}
				}
			}
		}
	}

	return failures
}

// litKey returns a canonical key for a basic literal, combining its kind and value
// so that e.g. the string "30" and the integer 30 are not conflated.
func litKey(lit *ast.BasicLit) string {
	return lit.Kind.String() + ":" + lit.Value
}

// Name returns the rule name.
func (*NoDuplicateConstantsRule) Name() string {
	return "noDuplicateConstants"
}

// Group returns the rule group.
func (*NoDuplicateConstantsRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoDuplicateConstantsRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
