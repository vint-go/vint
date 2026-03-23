package no_unnecessary_deref_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unnecessary_deref"
)

func TestNoUnnecessaryDeref(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unnecessary_deref", &no_unnecessary_deref.NoUnnecessaryDerefRule{})
}
