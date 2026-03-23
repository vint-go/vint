package no_doc_comment_stub

import (
	"go/ast"
	"strings"
	"unicode"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoDocCommentStubRule detects documentation comment stubs on exported
// declarations where the comment is a placeholder rather than real documentation.
type NoDocCommentStubRule struct{}

// stubPatterns contains lowercase patterns that indicate a stub doc comment.
// The check strips the leading "// Name " prefix and compares the remainder.
var stubPatterns = []string{
	"...",
	".",
	"xxx",
	"whatever",
	"todo",
	"nolint",
	"fixme",
	"hack",
	"placeholder",
	"-",
	"_",
}

// Apply applies the rule to given file.
func (r *NoDocCommentStubRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Name == nil || !d.Name.IsExported() {
				continue
			}
			if d.Doc == nil {
				continue
			}
			if isStubComment(d.Doc, d.Name.Name) {
				failures = append(failures, lint.Failure{
					Category:   lint.FailureCategoryStyle,
					Confidence: 1,
					Node:       d,
					Failure:    "doc comment for " + d.Name.Name + " appears to be a stub",
				})
			}
		case *ast.GenDecl:
			// For single-spec GenDecls (type, var, const) the Doc is on the GenDecl.
			// For grouped declarations, Doc may be on individual specs.
			if d.Doc != nil && len(d.Specs) == 1 {
				if name := specExportedName(d.Specs[0]); name != "" {
					if isStubComment(d.Doc, name) {
						failures = append(failures, lint.Failure{
							Category:   lint.FailureCategoryStyle,
							Confidence: 1,
							Node:       d,
							Failure:    "doc comment for " + name + " appears to be a stub",
						})
					}
				}
			}
			// Also check individual specs inside grouped decls.
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					if s.Doc != nil && s.Name != nil && s.Name.IsExported() {
						if isStubComment(s.Doc, s.Name.Name) {
							failures = append(failures, lint.Failure{
								Category:   lint.FailureCategoryStyle,
								Confidence: 1,
								Node:       s,
								Failure:    "doc comment for " + s.Name.Name + " appears to be a stub",
							})
						}
					}
				case *ast.ValueSpec:
					if s.Doc != nil && len(s.Names) > 0 && s.Names[0].IsExported() {
						if isStubComment(s.Doc, s.Names[0].Name) {
							failures = append(failures, lint.Failure{
								Category:   lint.FailureCategoryStyle,
								Confidence: 1,
								Node:       s,
								Failure:    "doc comment for " + s.Names[0].Name + " appears to be a stub",
							})
						}
					}
				}
			}
		}
	}

	return failures
}

// isStubComment returns true if the comment group looks like a documentation
// stub for an exported symbol with the given name.
func isStubComment(doc *ast.CommentGroup, name string) bool {
	if doc == nil || len(doc.List) == 0 {
		return false
	}

	// Combine all comment lines.
	var lines []string
	for _, c := range doc.List {
		text := c.Text
		if strings.HasPrefix(text, "//") {
			text = strings.TrimPrefix(text, "//")
		} else if strings.HasPrefix(text, "/*") {
			text = strings.TrimPrefix(text, "/*")
			text = strings.TrimSuffix(text, "*/")
		}
		text = strings.TrimSpace(text)
		if text != "" {
			lines = append(lines, text)
		}
	}

	if len(lines) == 0 {
		return false
	}

	// The conventional doc comment starts with the symbol name.
	// We consider the full text after joining all lines.
	fullText := strings.Join(lines, " ")

	// Strip leading name if present (e.g. "ProcessData ..." -> "...").
	body := fullText
	if strings.HasPrefix(fullText, name+" ") {
		body = strings.TrimSpace(fullText[len(name)+1:])
	} else if strings.EqualFold(fullText, name) {
		// Just the name with no description at all.
		return true
	}

	lower := strings.ToLower(body)
	lower = strings.TrimSpace(lower)

	// Check patterns before stripping punctuation.
	for _, pattern := range stubPatterns {
		if lower == pattern {
			return true
		}
	}

	// Remove trailing punctuation for a second round of comparison.
	stripped := strings.TrimRightFunc(lower, func(r rune) bool {
		return unicode.IsPunct(r) && r != '.'
	})

	for _, pattern := range stubPatterns {
		if stripped == pattern {
			return true
		}
	}

	// Check for repeated dots or ellipsis.
	trimmed := strings.TrimRight(lower, ".")
	if trimmed == "" && len(lower) > 0 {
		return true // all dots
	}

	return false
}

// specExportedName returns the exported name from a spec, or "" if not exported.
func specExportedName(spec ast.Spec) string {
	switch s := spec.(type) {
	case *ast.TypeSpec:
		if s.Name != nil && s.Name.IsExported() {
			return s.Name.Name
		}
	case *ast.ValueSpec:
		if len(s.Names) > 0 && s.Names[0].IsExported() {
			return s.Names[0].Name
		}
	}
	return ""
}

// Name returns the rule name.
func (*NoDocCommentStubRule) Name() string {
	return "noDocCommentStub"
}

// Group returns the rule group.
func (*NoDocCommentStubRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoDocCommentStubRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
