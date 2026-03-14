package no_atoi_overflow_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_atoi_overflow"
)

func TestNoAtoiOverflow(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_atoi_overflow", &no_atoi_overflow.NoAtoiOverflowRule{})
}
