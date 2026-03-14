package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestRedundantBuildTagRule(t *testing.T) {
	testRule(t, "redundant_build_tag", &rule.RedundantBuildTagRule{}, &lint.RuleConfig{})
}

func TestRedundantBuildTagRuleNoFailure(t *testing.T) {
	testRule(t, "redundant_build_tag_no_failure", &rule.RedundantBuildTagRule{}, &lint.RuleConfig{})
}

func TestRedundantBuildTagRuleGo116(t *testing.T) {
	testRule(t, "redundant_build_tag_go116", &rule.RedundantBuildTagRule{}, &lint.RuleConfig{})
}
