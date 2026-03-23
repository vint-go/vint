package use_time_date_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_time_date"
)

func TestUseTimeDate(t *testing.T) {
	functional_test_helpers.TestRule(t, "time_date_decimal_literal", &use_time_date.TimeDateRule{})
	functional_test_helpers.TestRule(t, "time_date_nil_timezone", &use_time_date.TimeDateRule{})
	functional_test_helpers.TestRule(t, "time_date_out_of_bounds", &use_time_date.TimeDateRule{})
}
