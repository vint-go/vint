package use_var_const_doc_prefix

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseVarConstDocPrefixRule checks that documentation comments for exported variables
// and constants start with the variable's or constant's name, following Go convention.
type UseVarConstDocPrefixRule struct{}

// Apply applies the rule to the given file.
func (r *UseVarConstDocPrefixRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.VAR && genDecl.Tok != token.CONST {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for _, name := range valueSpec.Names {
				if name == nil || !name.IsExported() {
					continue
				}

				// Doc comment can be on the ValueSpec (grouped declaration)
				// or on the GenDecl (standalone declaration).
				doc := valueSpec.Doc
				if doc == nil && len(genDecl.Specs) == 1 {
					doc = genDecl.Doc
				}

				// Skip variables/constants without doc comments.
				if doc == nil || len(doc.List) == 0 {
					continue
				}

				varName := name.Name
				docText := strings.TrimSpace(doc.Text())

				if docText == "" {
					continue
				}

				// Check if doc comment starts with the variable/constant name.
				if !strings.HasPrefix(docText, varName) {
					kind := "variable"
					if genDecl.Tok == token.CONST {
						kind = "constant"
					}
					node := ast.Node(genDecl)
					if valueSpec.Doc != nil {
						node = valueSpec
					}
					failures = append(failures, lint.Failure{
						Category:   lint.FailureCategoryStyle,
						Confidence: 1,
						Node:       node,
						Failure:    fmt.Sprintf("comment on exported %s %s should be of the form \"%s ...\"", kind, varName, varName),
					})
				}
			}
		}
	}

	return failures
}

// Name returns the rule name.
func (*UseVarConstDocPrefixRule) Name() string {
	return "useVarConstDocPrefix"
}

// Group returns the rule group.
func (*UseVarConstDocPrefixRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseVarConstDocPrefixRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
