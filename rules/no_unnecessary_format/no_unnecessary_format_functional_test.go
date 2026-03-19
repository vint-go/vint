package no_unnecessary_format_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unnecessary_format"
)

func TestNoUnnecessaryFormat(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unnecessary_format", &no_unnecessary_format.UnnecessaryFormatRule{})
}
