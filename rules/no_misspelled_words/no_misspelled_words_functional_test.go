package no_misspelled_words_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_misspelled_words"
)

func TestNoMisspelledWords(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_misspelled_words", &no_misspelled_words.NoMisspelledWordsRule{})
}
