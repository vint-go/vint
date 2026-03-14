package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestTimeDate(t *testing.T) {
	testRule(t, "time_date_decimal_literal", &rule.TimeDateRule{})
	testRule(t, "time_date_nil_timezone", &rule.TimeDateRule{})
	testRule(t, "time_date_out_of_bounds", &rule.TimeDateRule{})
}
