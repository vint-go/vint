package no_ignored_query_modification_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_ignored_query_modification"
)

func TestNoIgnoredQueryModification(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_ignored_query_modification", &no_ignored_query_modification.NoIgnoredQueryModificationRule{})
}
