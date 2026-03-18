package use_named_result_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_named_result"
)

func TestUseNamedResult(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_named_result", &use_named_result.UseNamedResultRule{})
}
