package use_package_comment_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_package_comment"
)

func TestUsePackageComment(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_package_comment", &use_package_comment.UsePackageCommentRule{})
}
