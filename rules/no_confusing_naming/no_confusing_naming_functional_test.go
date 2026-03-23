package no_confusing_naming_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_confusing_naming"
)

func TestNoConfusingNaming(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_confusing_naming", &no_confusing_naming.ConfusingNamingRule{})
}
