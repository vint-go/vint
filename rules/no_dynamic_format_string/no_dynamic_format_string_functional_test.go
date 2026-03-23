package no_dynamic_format_string_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_dynamic_format_string"
)

func TestNoDynamicFormatString(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_dynamic_format_string", &no_dynamic_format_string.NoDynamicFormatStringRule{})
}
