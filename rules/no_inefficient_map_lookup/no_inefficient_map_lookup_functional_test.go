package no_inefficient_map_lookup_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_inefficient_map_lookup"
)

func TestNoInefficientMapLookup(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_inefficient_map_lookup", &no_inefficient_map_lookup.InefficientMapLookupRule{})
}
