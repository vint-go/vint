package no_space_in_directive

import (
	"fmt"
	"strings"

	"github.com/strowk/vint/lint"
)

// NoSpaceInDirectiveRule detects Go compiler directives that contain a space
// between the comment slashes and the go: prefix.
type NoSpaceInDirectiveRule struct{}

// Apply applies the rule to given file.
func (r *NoSpaceInDirectiveRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, group := range file.AST.Comments {
		for _, comment := range group.List {
			text := comment.Text
			if strings.HasPrefix(text, "// go:") {
				directive := strings.TrimPrefix(text, "// ")
				failures = append(failures, lint.Failure{
					Category:   lint.FailureCategoryBadPractice,
					Confidence: 1,
					Node:       comment,
					Failure:    fmt.Sprintf("compiler directive contains space: %s", directive),
				})
			}
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoSpaceInDirectiveRule) Name() string {
	return "noSpaceInDirective"
}

// Group returns the rule group.
func (*NoSpaceInDirectiveRule) Group() string {
	return "correctness"
}
