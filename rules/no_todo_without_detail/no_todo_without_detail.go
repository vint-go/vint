package no_todo_without_detail

import (
	"fmt"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoTodoWithoutDetailRule detects TODO, FIX, FIXME, and BUG comments that
// lack additional context such as an assignee or explanatory details.
type NoTodoWithoutDetailRule struct{}

// todoKeywords are the comment prefixes that should have detail.
// Keywords are ordered longest-first so that "FIXME" is checked before "FIX".
var todoKeywords = []string{"TODO", "FIXME", "FIX", "BUG"}

// Apply applies the rule to given file.
func (r *NoTodoWithoutDetailRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, group := range file.AST.Comments {
		for _, comment := range group.List {
			text := comment.Text

			// Strip the comment prefix to get the body.
			var body string
			if strings.HasPrefix(text, "//") {
				body = strings.TrimPrefix(text, "//")
			} else if strings.HasPrefix(text, "/*") {
				body = strings.TrimPrefix(text, "/*")
				body = strings.TrimSuffix(body, "*/")
			} else {
				continue
			}

			// Trim leading whitespace from the body.
			trimmed := strings.TrimLeft(body, " \t")

			for _, keyword := range todoKeywords {
				if !strings.HasPrefix(trimmed, keyword) {
					continue
				}

				// Check what follows the keyword.
				rest := trimmed[len(keyword):]

				// If there's nothing after the keyword, or only whitespace,
				// this is a bare TODO/FIX/FIXME/BUG without detail.
				if rest == "" || strings.TrimSpace(rest) == "" {
					failures = append(failures, lint.Failure{
						Category:   lint.FailureCategoryContent,
						Confidence: 1,
						Node:       comment,
						Failure:    fmt.Sprintf("%s comment without detail or assignee", keyword),
					})
					break
				}

				// If the keyword is immediately followed by a letter or digit,
				// it's part of a longer word (e.g., "TODOS", "BUGFIX") - not a match.
				firstChar := rest[0]
				if (firstChar >= 'a' && firstChar <= 'z') ||
					(firstChar >= 'A' && firstChar <= 'Z') ||
					(firstChar >= '0' && firstChar <= '9') {
					break
				}

				// The keyword is followed by a non-alphanumeric character.
				// Check if there's actual detail after the keyword and its delimiter.
				// Acceptable patterns: "TODO(assignee): detail", "TODO: detail", "TODO - detail"
				afterDelim := strings.TrimLeft(rest, " \t:(-")
				if afterDelim == "" || strings.TrimSpace(afterDelim) == "" {
					failures = append(failures, lint.Failure{
						Category:   lint.FailureCategoryContent,
						Confidence: 1,
						Node:       comment,
						Failure:    fmt.Sprintf("%s comment without detail or assignee", keyword),
					})
				}
				break
			}
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoTodoWithoutDetailRule) Name() string {
	return "noTodoWithoutDetail"
}

// Group returns the rule group.
func (*NoTodoWithoutDetailRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoTodoWithoutDetailRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
