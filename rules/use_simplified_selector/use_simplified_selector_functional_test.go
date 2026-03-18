package use_simplified_selector_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_simplified_selector"
)

func TestUseSimplifiedSelector(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_simplified_selector", &use_simplified_selector.UseSimplifiedSelectorRule{})
}
