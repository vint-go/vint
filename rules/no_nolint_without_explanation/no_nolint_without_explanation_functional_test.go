package no_nolint_without_explanation_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_nolint_without_explanation"
)

func TestNoNolintWithoutExplanation(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_nolint_without_explanation", &no_nolint_without_explanation.NoNolintWithoutExplanationRule{})
}
