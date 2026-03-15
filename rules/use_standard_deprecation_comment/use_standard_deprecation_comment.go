package use_standard_deprecation_comment

import (
	"go/ast"
	"regexp"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// standardPattern matches the correct Go deprecation comment format:
// "// Deprecated: <explanation>"
// It must start with exactly "// Deprecated: " (with proper casing, colon, and space).
var standardPattern = regexp.MustCompile(`^// Deprecated: \S`)

// heuristicPatterns detect comments that look like deprecation notices
// but do not follow the standard "// Deprecated: ..." format.
var heuristicPatterns = []*regexp.Regexp{
	// Wrong casing: DEPRECATED:, deprecated:, Deprecated: at start (handled separately)
	regexp.MustCompile(`(?i)^//\s*deprecated\s*:`),
	// Using comma instead of colon: "Deprecated, ..."
	regexp.MustCompile(`(?i)^//\s*deprecated\s*,`),
	// Alternative inline patterns: "this function is deprecated", "this method is deprecated", etc.
	regexp.MustCompile(`(?i)^//\s*(this\s+\w+\s+is\s+)?deprecated[\.\,\s]`),
	// "deprecated. use ..." pattern
	regexp.MustCompile(`(?i)\bdeprecated\.\s*use\b`),
	// Typos of "Deprecated" near the start of the comment
	regexp.MustCompile(`(?i)^//\s*(deprected|depricated|depracated|depreacted|depercated|depriciated|depreciated|deprcated|depprecated)\s*[:\.,\s]`),
}

// UseStandardDeprecationCommentRule detects malformed deprecation doc-comments.
type UseStandardDeprecationCommentRule struct{}

// Apply applies the rule to given file.
func (r *UseStandardDeprecationCommentRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, group := range file.AST.Comments {
		for _, comment := range group.List {
			if f := checkComment(comment); f != nil {
				failures = append(failures, *f)
			}
		}
	}

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseStandardDeprecationCommentRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	comment, ok := node.(*ast.Comment)
	if !ok {
		return nil
	}

	if f := checkComment(comment); f != nil {
		return []lint.Failure{*f}
	}

	return nil
}

// checkComment checks a single comment for malformed deprecation notices.
func checkComment(comment *ast.Comment) *lint.Failure {
	text := comment.Text

	// Only consider line comments (// style)
	if !strings.HasPrefix(text, "//") {
		return nil
	}

	// Skip if it already matches the standard format
	if standardPattern.MatchString(text) {
		return nil
	}

	// Check if the comment looks like a deprecation notice
	if looksLikeDeprecationComment(text) {
		return &lint.Failure{
			Category:   lint.FailureCategoryStyle,
			Confidence: 1,
			Node:       comment,
			Failure:    `comment looks like a deprecation notice but does not follow the standard "// Deprecated: <explanation>" format`,
		}
	}

	return nil
}

// looksLikeDeprecationComment checks whether a comment text heuristically
// looks like a deprecation notice.
func looksLikeDeprecationComment(text string) bool {
	for _, pat := range heuristicPatterns {
		if pat.MatchString(text) {
			return true
		}
	}
	return false
}

// Name returns the rule name.
func (*UseStandardDeprecationCommentRule) Name() string {
	return "useStandardDeprecationComment"
}

// Group returns the rule group.
func (*UseStandardDeprecationCommentRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseStandardDeprecationCommentRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
