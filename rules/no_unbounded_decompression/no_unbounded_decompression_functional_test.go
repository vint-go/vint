package no_unbounded_decompression_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unbounded_decompression"
)

func TestNoUnboundedDecompression(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unbounded_decompression", &no_unbounded_decompression.NoUnboundedDecompressionRule{})
}
