package use_standard_codegen_comment_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_standard_codegen_comment"
)

func TestUseStandardCodegenComment(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_standard_codegen_comment", &use_standard_codegen_comment.UseStandardCodegenCommentRule{})
}
