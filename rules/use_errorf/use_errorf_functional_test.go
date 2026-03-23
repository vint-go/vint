package use_errorf_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_errorf"
)

func TestUseErrorf(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_errorf", &use_errorf.ErrorfRule{})
}
