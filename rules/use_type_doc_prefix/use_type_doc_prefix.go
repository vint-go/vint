package use_type_doc_prefix

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseTypeDocPrefixRule checks that documentation comments for exported types
// start with the type's name, following Go convention.
type UseTypeDocPrefixRule struct{}

// Apply applies the rule to the given file.
func (r *UseTypeDocPrefixRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// Only check exported types.
			if typeSpec.Name == nil || !typeSpec.Name.IsExported() {
				continue
			}

			// Doc comment can be on the TypeSpec (grouped declaration)
			// or on the GenDecl (standalone declaration).
			doc := typeSpec.Doc
			if doc == nil && len(genDecl.Specs) == 1 {
				doc = genDecl.Doc
			}

			// Skip types without doc comments.
			if doc == nil || len(doc.List) == 0 {
				continue
			}

			name := typeSpec.Name.Name
			docText := strings.TrimSpace(doc.Text())

			if docText == "" {
				continue
			}

			// Check if doc comment starts with the type name.
			if !strings.HasPrefix(docText, name) {
				node := ast.Node(genDecl)
				if typeSpec.Doc != nil {
					node = typeSpec
				}
				failures = append(failures, lint.Failure{
					Category:   lint.FailureCategoryStyle,
					Confidence: 1,
					Node:       node,
					Failure:    fmt.Sprintf("comment on exported type %s should be of the form \"%s ...\"", name, name),
				})
			}
		}
	}

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseTypeDocPrefixRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	genDecl, ok := node.(*ast.GenDecl)
	if !ok || genDecl.Tok != token.TYPE {
		return nil
	}

	var failures []lint.Failure

	for _, spec := range genDecl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		if typeSpec.Name == nil || !typeSpec.Name.IsExported() {
			continue
		}

		doc := typeSpec.Doc
		if doc == nil && len(genDecl.Specs) == 1 {
			doc = genDecl.Doc
		}

		if doc == nil || len(doc.List) == 0 {
			continue
		}

		name := typeSpec.Name.Name
		docText := strings.TrimSpace(doc.Text())

		if docText == "" {
			continue
		}

		if !strings.HasPrefix(docText, name) {
			failureNode := ast.Node(genDecl)
			if typeSpec.Doc != nil {
				failureNode = typeSpec
			}
			failures = append(failures, lint.Failure{
				Category:   lint.FailureCategoryStyle,
				Confidence: 1,
				Node:       failureNode,
				Failure:    fmt.Sprintf("comment on exported type %s should be of the form \"%s ...\"", name, name),
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*UseTypeDocPrefixRule) Name() string {
	return "useTypeDocPrefix"
}

// Group returns the rule group.
func (*UseTypeDocPrefixRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseTypeDocPrefixRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
