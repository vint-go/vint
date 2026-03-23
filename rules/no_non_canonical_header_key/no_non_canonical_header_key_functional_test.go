package no_non_canonical_header_key_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_non_canonical_header_key"
)

func TestNoNonCanonicalHeaderKey(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_non_canonical_header_key", &no_non_canonical_header_key.NoNonCanonicalHeaderKeyRule{})
}
