package use_time_naming_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_time_naming"
)

func TestUseTimeNaming(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_time_naming", &use_time_naming.TimeNamingRule{})
}
