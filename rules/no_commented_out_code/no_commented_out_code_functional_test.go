package no_commented_out_code_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_commented_out_code"
)

func TestNoCommentedOutCode(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_commented_out_code", &no_commented_out_code.NoCommentedOutCodeRule{})
}
