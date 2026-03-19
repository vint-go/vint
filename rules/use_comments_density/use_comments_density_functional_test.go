package use_comments_density_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_comments_density"
)

func TestUseCommentsDensity(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_comments_density_1", &use_comments_density.CommentsDensityRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(60)},
	})

	functional_test_helpers.TestRule(t, "use_comments_density_2", &use_comments_density.CommentsDensityRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(90)},
	})

	functional_test_helpers.TestRule(t, "use_comments_density_3", &use_comments_density.CommentsDensityRule{}, &lint.RuleConfig{})
}
