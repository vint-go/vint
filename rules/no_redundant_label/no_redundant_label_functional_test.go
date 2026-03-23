package no_redundant_label_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_redundant_label"
)

func TestNoRedundantLabel(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_label", &no_redundant_label.NoRedundantLabelRule{})
}
