package no_ignored_query_modification_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_ignored_query_modification"
)

func TestNoIgnoredQueryModification(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_ignored_query_modification", &no_ignored_query_modification.NoIgnoredQueryModificationRule{})
}
