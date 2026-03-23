package no_malformed_build_tag_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_malformed_build_tag"
)

func TestNoMalformedBuildTag(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_malformed_build_tag", &no_malformed_build_tag.NoMalformedBuildTagRule{})
}
