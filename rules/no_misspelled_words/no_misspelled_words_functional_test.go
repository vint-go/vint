package no_misspelled_words_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_misspelled_words"
)

func TestNoMisspelledWords(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_misspelled_words", &no_misspelled_words.NoMisspelledWordsRule{})
}
