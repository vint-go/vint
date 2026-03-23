package use_comment_spacing_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_comment_spacing"
)

func TestUseCommentSpacing(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_comment_spacing", &use_comment_spacing.UseCommentSpacingRule{})
}
