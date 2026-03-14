package no_trojan_source_bidi_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_trojan_source_bidi"
)

func TestNoTrojanSourceBidi(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_trojan_source_bidi", &no_trojan_source_bidi.NoTrojanSourceBidiRule{})
}
