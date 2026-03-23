package use_package_comments_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_package_comments"
)

func TestUsePackageComments(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_package_comments", &use_package_comments.PackageCommentsRule{})
}
