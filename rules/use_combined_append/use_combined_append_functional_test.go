package use_combined_append_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_combined_append"
)

func TestUseCombinedAppend(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_combined_append", &use_combined_append.UseCombinedAppendRule{})
}
