package use_type_assert_result_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_type_assert_result"
)

func TestUseTypeAssertResult(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_type_assert_result", &use_type_assert_result.UseTypeAssertResultRule{})
}
