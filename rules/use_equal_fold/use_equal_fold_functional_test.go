package use_equal_fold_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_equal_fold"
)

func TestUseEqualFold(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_equal_fold", &use_equal_fold.UseEqualFoldRule{})
}
