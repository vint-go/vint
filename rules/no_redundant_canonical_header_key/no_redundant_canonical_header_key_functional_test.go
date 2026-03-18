package no_redundant_canonical_header_key_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_redundant_canonical_header_key"
)

func TestNoRedundantCanonicalHeaderKey(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_canonical_header_key", &no_redundant_canonical_header_key.NoRedundantCanonicalHeaderKeyRule{})
}
