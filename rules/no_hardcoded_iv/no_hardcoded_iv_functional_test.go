package no_hardcoded_iv_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_hardcoded_iv"
)

func TestNoHardcodedIv(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_hardcoded_iv", &no_hardcoded_iv.NoHardcodedIvRule{})
}
