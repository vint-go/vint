package use_comment_spacing_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_comment_spacing"
)

func TestUseCommentSpacing(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_comment_spacing", &use_comment_spacing.UseCommentSpacingRule{})
}
