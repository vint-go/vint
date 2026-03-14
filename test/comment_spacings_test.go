package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestCommentSpacings(t *testing.T) {
	testRule(t, "comment_spacings", &rule.CommentSpacingsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"myOwnDirective:", "+optional"},
	},
	)
}
