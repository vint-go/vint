package no_commented_out_import

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoCommentedOutImportRule detects commented-out imports inside import blocks.
type NoCommentedOutImportRule struct{}

// Apply applies the rule to given file.
func (r *NoCommentedOutImportRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	// Find all import GenDecl blocks
	for _, decl := range file.AST.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			continue
		}

		// Only check parenthesized import blocks (import (...))
		if !genDecl.Lparen.IsValid() {
			continue
		}

		// Check all comment groups in the file for ones that fall inside this import block
		for _, cg := range file.AST.Comments {
			for _, comment := range cg.List {
				// Comment must be inside the import block (between Lparen and Rparen)
				if comment.Pos() <= genDecl.Lparen || comment.Pos() >= genDecl.Rparen {
					continue
				}

				if isCommentedOutImport(comment.Text) {
					failures = append(failures, lint.Failure{
						Confidence: 1,
						Failure:    fmt.Sprintf("commented-out import: remove or uncomment %s", extractImportPath(comment.Text)),
						Node:       comment,
						Category:   lint.FailureCategoryImports,
					})
				}
			}
		}
	}

	return failures
}

// isCommentedOutImport checks if a comment line looks like a commented-out import.
// It handles patterns like:
//
//	// "fmt"
//	// "os/exec"
//	// alias "pkg/path"
//	/* "net/http" */
func isCommentedOutImport(text string) bool {
	var content string
	if strings.HasPrefix(text, "//") {
		content = strings.TrimPrefix(text, "//")
	} else if strings.HasPrefix(text, "/*") && strings.HasSuffix(text, "*/") {
		content = strings.TrimPrefix(text, "/*")
		content = strings.TrimSuffix(content, "*/")
	} else {
		return false
	}

	content = strings.TrimSpace(content)

	if content == "" {
		return false
	}

	// Check for direct quoted import path: "some/path"
	if isQuotedImportPath(content) {
		return true
	}

	// Check for aliased import: alias "some/path" or . "some/path" or _ "some/path"
	parts := strings.Fields(content)
	if len(parts) == 2 && isValidImportAlias(parts[0]) && isQuotedImportPath(parts[1]) {
		return true
	}

	return false
}

// isQuotedImportPath checks if s looks like a quoted Go import path.
func isQuotedImportPath(s string) bool {
	if len(s) < 2 {
		return false
	}
	if s[0] != '"' || s[len(s)-1] != '"' {
		return false
	}
	inner := s[1 : len(s)-1]
	if inner == "" {
		return false
	}
	// Basic validation: import paths contain letters, digits, /, ., -, _
	for _, ch := range inner {
		if !isImportPathChar(ch) {
			return false
		}
	}
	return true
}

// isImportPathChar returns true if the rune is valid in a Go import path.
func isImportPathChar(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') ||
		ch == '/' || ch == '.' || ch == '-' || ch == '_' || ch == '~' || ch == '+'
}

// isValidImportAlias checks if the string is a valid Go import alias identifier.
func isValidImportAlias(s string) bool {
	if s == "." || s == "_" {
		return true
	}
	// Must be a valid Go identifier
	for i, ch := range s {
		if i == 0 {
			if !isLetter(ch) {
				return false
			}
		} else {
			if !isLetter(ch) && !isDigit(ch) {
				return false
			}
		}
	}
	return len(s) > 0
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

// extractImportPath extracts the import path string from a commented-out import line.
func extractImportPath(text string) string {
	var content string
	if strings.HasPrefix(text, "//") {
		content = strings.TrimPrefix(text, "//")
	} else if strings.HasPrefix(text, "/*") && strings.HasSuffix(text, "*/") {
		content = strings.TrimPrefix(text, "/*")
		content = strings.TrimSuffix(content, "*/")
	}
	content = strings.TrimSpace(content)

	// If aliased, return the full thing
	parts := strings.Fields(content)
	if len(parts) == 2 && isValidImportAlias(parts[0]) && isQuotedImportPath(parts[1]) {
		return content
	}

	return content
}

// Name returns the rule name.
func (*NoCommentedOutImportRule) Name() string {
	return "noCommentedOutImport"
}

// Group returns the rule group.
func (*NoCommentedOutImportRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoCommentedOutImportRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
