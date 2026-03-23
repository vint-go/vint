package use_strings_contains_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_strings_contains"
)

func TestUseStringsContains(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_strings_contains", &use_strings_contains.UseStringsContainsRule{})
}
