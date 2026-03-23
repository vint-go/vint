package use_infinite_for_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_infinite_for"
)

func TestUseInfiniteFor(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_infinite_for", &use_infinite_for.UseInfiniteForRule{})
}
