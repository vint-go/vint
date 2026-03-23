package use_standard_codegen_comment

import (
	"go/ast"
	"regexp"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// standardPattern matches the correct Go codegen comment format:
// "// Code generated .* DO NOT EDIT."
var standardPattern = regexp.MustCompile(`^// Code generated .* DO NOT EDIT\.$`)

// heuristicPatterns are case-insensitive patterns that suggest a comment is
// about code generation but may not match the standard format.
var heuristicPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\bgenerated\b.*\bdo\s+not\s+edit\b`),
	regexp.MustCompile(`(?i)\bdo\s+not\s+edit\b.*\bgenerated\b`),
	regexp.MustCompile(`(?i)\bauto[\-\s]?generated\b`),
	regexp.MustCompile(`(?i)\bcode[\-\s]?generated\b`),
	regexp.MustCompile(`(?i)\bgenerated\s+(by|from|code|file)\b`),
	regexp.MustCompile(`(?i)\bthis file was .*generated\b`),
	regexp.MustCompile(`(?i)\bgenerated.*(do not|don'?t)\s+(edit|modify)\b`),
}

// UseStandardCodegenCommentRule detects malformed code-generated file comments.
type UseStandardCodegenCommentRule struct{}

// Apply applies the rule to given file.
func (r *UseStandardCodegenCommentRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, group := range file.AST.Comments {
		for _, comment := range group.List {
			text := comment.Text

			// Only consider line comments (// style)
			if !strings.HasPrefix(text, "//") {
				continue
			}

			// Skip if it already matches the standard format
			if standardPattern.MatchString(text) {
				continue
			}

			// Check if the comment looks like a codegen marker
			if looksLikeCodegenComment(text) {
				failures = append(failures, lint.Failure{
					Category:   lint.FailureCategoryStyle,
					Confidence: 1,
					Node:       comment,
					Failure:    "comment looks like a code-generation marker but does not follow the standard Go format",
				})
			}
		}
	}

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseStandardCodegenCommentRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	comment, ok := node.(*ast.Comment)
	if !ok {
		return nil
	}

	text := comment.Text
	if !strings.HasPrefix(text, "//") {
		return nil
	}

	if standardPattern.MatchString(text) {
		return nil
	}

	if looksLikeCodegenComment(text) {
		return []lint.Failure{{
			Category:   lint.FailureCategoryStyle,
			Confidence: 1,
			Node:       comment,
			Failure:    "comment looks like a code-generation marker but does not follow the standard Go format",
		}}
	}

	return nil
}

// looksLikeCodegenComment checks whether a comment text heuristically looks
// like a code-generation marker.
func looksLikeCodegenComment(text string) bool {
	for _, pat := range heuristicPatterns {
		if pat.MatchString(text) {
			return true
		}
	}
	return false
}

// Name returns the rule name.
func (*UseStandardCodegenCommentRule) Name() string {
	return "useStandardCodegenComment"
}

// Group returns the rule group.
func (*UseStandardCodegenCommentRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseStandardCodegenCommentRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
