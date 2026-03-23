package no_unbounded_decompression_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unbounded_decompression"
)

func TestNoUnboundedDecompression(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unbounded_decompression", &no_unbounded_decompression.NoUnboundedDecompressionRule{})
}
