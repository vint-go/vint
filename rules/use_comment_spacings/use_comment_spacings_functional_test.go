package use_comment_spacings_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_comment_spacings"
)

func TestUseCommentSpacings(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_comment_spacings", &use_comment_spacings.CommentSpacingsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"myOwnDirective:", "+optional"},
	})
}
