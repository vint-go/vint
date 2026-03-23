package no_malformed_struct_tag_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_malformed_struct_tag"
)

func TestNoMalformedStructTag(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_malformed_struct_tag", &no_malformed_struct_tag.NoMalformedStructTagRule{})
}
