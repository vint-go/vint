package no_huge_param_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_huge_param"
)

func TestNoHugeParam(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_huge_param", &no_huge_param.NoHugeParamRule{})
}
