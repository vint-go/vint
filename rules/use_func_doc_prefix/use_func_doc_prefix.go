package use_func_doc_prefix

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseFuncDocPrefixRule checks that documentation comments for exported functions
// start with the function's name, following Go convention.
type UseFuncDocPrefixRule struct{}

// Apply applies the rule to the given file.
func (r *UseFuncDocPrefixRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// Only check exported functions.
		if funcDecl.Name == nil || !funcDecl.Name.IsExported() {
			continue
		}

		// Skip functions without doc comments.
		if funcDecl.Doc == nil || len(funcDecl.Doc.List) == 0 {
			continue
		}

		name := funcDecl.Name.Name
		docText := funcDecl.Doc.Text()
		docText = strings.TrimSpace(docText)

		if docText == "" {
			continue
		}

		// Check if doc comment starts with the function name.
		if !strings.HasPrefix(docText, name) {
			failures = append(failures, lint.Failure{
				Category:   lint.FailureCategoryStyle,
				Confidence: 1,
				Node:       funcDecl,
				Failure:    fmt.Sprintf("comment on exported function %s should be of the form \"%s ...\"", name, name),
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*UseFuncDocPrefixRule) Name() string {
	return "useFuncDocPrefix"
}

// Group returns the rule group.
func (*UseFuncDocPrefixRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseFuncDocPrefixRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
