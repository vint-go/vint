package use_convenience_func_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_convenience_func"
)

func TestUseConvenienceFunc(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_convenience_func", &use_convenience_func.UseConvenienceFuncRule{})
}
