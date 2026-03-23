package no_excessive_public_structs_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_excessive_public_structs"
)

func TestNoExcessivePublicStructs(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_excessive_public_structs", &no_excessive_public_structs.NoExcessivePublicStructsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(1)},
	})
}

func TestNoExcessivePublicStructsDefaultConfig(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_excessive_public_structs_ok", &no_excessive_public_structs.NoExcessivePublicStructsRule{}, &lint.RuleConfig{})
}
