package no_unused_nolint_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unused_nolint"
)

func TestNoUnusedNolint(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unused_nolint", &no_unused_nolint.NoUnusedNolintRule{})
}
