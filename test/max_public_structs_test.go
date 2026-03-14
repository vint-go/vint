package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestMaxPublicStructs(t *testing.T) {
	testRule(t, "max_public_structs", &rule.MaxPublicStructsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(1)},
	})
}

func TestMaxPublicStructsDefaultConfig(t *testing.T) {
	testRule(t, "max_public_structs_ok", &rule.MaxPublicStructsRule{}, &lint.RuleConfig{})
}
