package use_standard_deprecation_comment_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_standard_deprecation_comment"
)

func TestUseStandardDeprecationComment(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_standard_deprecation_comment", &use_standard_deprecation_comment.UseStandardDeprecationCommentRule{})
}
