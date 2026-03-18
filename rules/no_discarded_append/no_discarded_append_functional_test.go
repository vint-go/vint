package no_discarded_append_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_discarded_append"
)

func TestNoDiscardedAppend(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_discarded_append", &no_discarded_append.NoDiscardedAppendRule{})
}
