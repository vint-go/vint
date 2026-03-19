package no_unexported_naming_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unexported_naming"
)

func TestNoUnexportedNaming(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unexported_naming", &no_unexported_naming.UnexportedNamingRule{})
}
