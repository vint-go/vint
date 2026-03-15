package use_comment_spacing

import (
	"regexp"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// directiveOrSpecialRE matches comments that are exempted from the space-after-// rule.
// This includes:
//   - Compiler directives like //go:generate, //line, //extern, //export
//   - Keyed directives matching //key:value patterns (e.g., //nolint:foo)
//   - Custom markers starting with +, -, #, or !
//   - Linting directives like //nolint, //noinspection
//   - IDE folding regions like //region, //editor-fold
var directiveOrSpecialRE = regexp.MustCompile(
	`^//` +
		`(?:` +
		`line |extern |export |` + // compiler directives with space
		`[a-z0-9]+:[a-z0-9]|` + // keyed directives like //go:generate, //nolint:foo
		`[+\-#!]|` + // custom markers
		`nolint\b|` + // nolint directive
		`noinspection\b|` + // noinspection directive
		`region\b|` + // IDE folding region start
		`endregion\b|` + // IDE folding region end
		`editor-fold` + // editor-fold directive
		`)`,
)

// UseCommentSpacingRule checks that single-line comments have a space after //.
type UseCommentSpacingRule struct{}

// Apply applies the rule to given file.
func (r *UseCommentSpacingRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, cg := range file.AST.Comments {
		for _, comment := range cg.List {
			commentLine := comment.Text

			// Skip short comments (just "//" is fine)
			if len(commentLine) < 3 {
				continue
			}

			// Skip block comments (/* ... */)
			if strings.HasPrefix(commentLine, "/*") {
				continue
			}

			// Check if the third character (after //) is a space or tab -- that's fine
			if commentLine[2] == ' ' || commentLine[2] == '\t' {
				continue
			}

			// Check if this is an exempted directive or special comment
			if directiveOrSpecialRE.MatchString(commentLine) {
				continue
			}

			failures = append(failures, lint.Failure{
				Node:       comment,
				Confidence: 1,
				Category:   lint.FailureCategoryStyle,
				Failure:    "no space between comment delimiter and comment text",
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*UseCommentSpacingRule) Name() string {
	return "useCommentSpacing"
}

// Group returns the rule group.
func (*UseCommentSpacingRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseCommentSpacingRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
