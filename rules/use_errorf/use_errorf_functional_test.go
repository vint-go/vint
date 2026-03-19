package use_errorf_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_errorf"
)

func TestUseErrorf(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_errorf", &use_errorf.ErrorfRule{})
}
