package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestCommentsDensity(t *testing.T) {
	testRule(t, "comments_density_1", &rule.CommentsDensityRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(60)},
	})

	testRule(t, "comments_density_2", &rule.CommentsDensityRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(90)},
	})

	testRule(t, "comments_density_3", &rule.CommentsDensityRule{}, &lint.RuleConfig{})
}
