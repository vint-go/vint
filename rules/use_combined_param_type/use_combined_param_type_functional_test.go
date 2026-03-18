package use_combined_param_type_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_combined_param_type"
)

func TestUseCombinedParamType(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_combined_param_type", &use_combined_param_type.UseCombinedParamTypeRule{})
}
