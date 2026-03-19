package use_comment_spacings_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_comment_spacings"
)

func TestUseCommentSpacings(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_comment_spacings", &use_comment_spacings.CommentSpacingsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"myOwnDirective:", "+optional"},
	})
}
