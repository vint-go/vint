package no_discarded_append_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_discarded_append"
)

func TestNoDiscardedAppend(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_discarded_append", &no_discarded_append.NoDiscardedAppendRule{})
}
