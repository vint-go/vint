package no_banned_characters_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_banned_characters"
)

func TestNoBannedCharactersDefault(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_banned_characters_default", &no_banned_characters.NoBannedCharactersRule{})
}

// Test banned characters in a const, var and func names.
// One banned character is in the comment and should not be checked.
// One banned character from the list is not present in the fixture file.
func TestNoBannedCharacters(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_banned_characters", &no_banned_characters.NoBannedCharactersRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"Ω", "Σ", "σ", "1"},
	})
}
