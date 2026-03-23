package no_blank_import

import (
	"go/ast"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoBlankImportRule lints blank imports.
type NoBlankImportRule struct{}

// Name returns the rule name.
func (*NoBlankImportRule) Name() string {
	return "noBlankImport"
}

// Group returns the rule group.
func (*NoBlankImportRule) Group() string {
	return "style"
}

// Apply applies the rule to given file.
func (r *NoBlankImportRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if file.Pkg.IsMain() || file.IsTest() {
		return nil
	}

	const (
		message         = "a blank import should be only in a main or test package, or have a comment justifying it"
		embedImportPath = `"embed"`
	)

	var failures []lint.Failure

	// The first element of each contiguous group of blank imports should have
	// an explanatory comment of some kind.
	for i, imp := range file.AST.Imports {
		pos := file.ToPosition(imp.Pos())

		if !isBlank(imp.Name) {
			continue // Ignore non-blank imports.
		}

		isNotFirstElement := i > 0
		if isNotFirstElement {
			prev := file.AST.Imports[i-1]
			prevPos := file.ToPosition(prev.Pos())

			isSubsequentBlancInAGroup := prevPos.Line+1 == pos.Line && prev.Path.Value != embedImportPath && isBlank(prev.Name)
			if isSubsequentBlancInAGroup {
				continue
			}
		}

		if imp.Path.Value == embedImportPath && r.fileHasValidEmbedComment(file.AST) {
			continue
		}

		// This is the first blank import of a group.
		if imp.Doc == nil && imp.Comment == nil {
			failures = append(failures, lint.Failure{Failure: message, Category: lint.FailureCategoryImports, Node: imp, Confidence: 1})
		}
	}

	return failures
}

func (*NoBlankImportRule) fileHasValidEmbedComment(fileAst *ast.File) bool {
	for _, commentGroup := range fileAst.Comments {
		for _, comment := range commentGroup.List {
			if strings.HasPrefix(comment.Text, "//go:embed ") {
				return true
			}
		}
	}

	return false
}

// isBlank returns whether id is the blank identifier "_".
// If id == nil, the answer is false.
func isBlank(id *ast.Ident) bool { return id != nil && id.Name == "_" }

// CacheTier returns the cache tier for this rule.
func (*NoBlankImportRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
