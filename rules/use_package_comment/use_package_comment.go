package use_package_comment

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UsePackageCommentRule checks that packages have a correctly formatted
// package comment beginning with "Package <name> ...".
type UsePackageCommentRule struct{}

// Apply applies the rule to given file.
func (r *UsePackageCommentRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if file.IsTest() {
		return nil
	}

	fileAST := file.AST
	pkgName := fileAST.Name.Name

	// Skip the main package -- it has no godoc requirements.
	if pkgName == "main" {
		return nil
	}

	// If there is no doc comment at all, report it.
	if isEmptyDoc(fileAST.Doc) {
		return nil // other files in the package might have the comment; skip silently per-file
	}

	docText := fileAST.Doc.Text()

	// Skip directive comments (e.g. //go:build).
	if isDirectiveComment(docText) {
		return nil
	}

	prefix := "Package " + pkgName + " "
	if !strings.HasPrefix(docText, prefix) {
		return []lint.Failure{{
			Category:   lint.FailureCategoryStyle,
			Node:       fileAST.Name,
			Confidence: 1,
			Failure:    fmt.Sprintf(`package comment should be of the form "Package %s ..."`, pkgName),
		}}
	}

	return nil
}

// Name returns the rule name.
func (*UsePackageCommentRule) Name() string {
	return "usePackageComment"
}

// Group returns the rule group.
func (*UsePackageCommentRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UsePackageCommentRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

func isEmptyDoc(cg *ast.CommentGroup) bool {
	return cg == nil || cg.Text() == ""
}

// isDirectiveComment reports whether the comment text starts with a Go directive.
func isDirectiveComment(text string) bool {
	return strings.HasPrefix(text, "//go:") || strings.HasPrefix(text, "// +build")
}
