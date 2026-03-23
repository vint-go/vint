package no_duplicate_cutset_chars_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_duplicate_cutset_chars"
)

func TestNoDuplicateCutsetChars(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_duplicate_cutset_chars", &no_duplicate_cutset_chars.NoDuplicateCutsetCharsRule{})
}
