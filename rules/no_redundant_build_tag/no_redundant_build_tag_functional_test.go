package no_redundant_build_tag_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_redundant_build_tag"
)

func TestNoRedundantBuildTag(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_build_tag", &no_redundant_build_tag.RedundantBuildTagRule{})
}

func TestNoRedundantBuildTagNoFailure(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_build_tag_no_failure", &no_redundant_build_tag.RedundantBuildTagRule{})
}

func TestNoRedundantBuildTagGo116(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_build_tag_go116", &no_redundant_build_tag.RedundantBuildTagRule{})
}
