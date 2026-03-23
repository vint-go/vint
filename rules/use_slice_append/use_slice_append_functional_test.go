package use_slice_append_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_slice_append"
)

func TestUseSliceAppend(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_slice_append", &use_slice_append.UseSliceAppendRule{})
}
