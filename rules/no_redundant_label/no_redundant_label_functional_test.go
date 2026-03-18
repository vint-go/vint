package no_redundant_label_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_redundant_label"
)

func TestNoRedundantLabel(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_label", &no_redundant_label.NoRedundantLabelRule{})
}
