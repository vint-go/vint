package no_unexported_naming_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unexported_naming"
)

func TestNoUnexportedNaming(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unexported_naming", &no_unexported_naming.UnexportedNamingRule{})
}
