package no_atoi_overflow_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_atoi_overflow"
)

func TestNoAtoiOverflow(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_atoi_overflow", &no_atoi_overflow.NoAtoiOverflowRule{})
}
