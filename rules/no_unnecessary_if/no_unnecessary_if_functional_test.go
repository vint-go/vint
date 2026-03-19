package no_unnecessary_if_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unnecessary_if"
)

func TestNoUnnecessaryIf(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unnecessary_if", &no_unnecessary_if.UnnecessaryIfRule{})
}
