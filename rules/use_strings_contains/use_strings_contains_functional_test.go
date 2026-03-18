package use_strings_contains_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_strings_contains"
)

func TestUseStringsContains(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_strings_contains", &use_strings_contains.UseStringsContainsRule{})
}
