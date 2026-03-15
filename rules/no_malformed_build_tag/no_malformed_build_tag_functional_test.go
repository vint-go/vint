package no_malformed_build_tag_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_malformed_build_tag"
)

func TestNoMalformedBuildTag(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_malformed_build_tag", &no_malformed_build_tag.NoMalformedBuildTagRule{})
}
