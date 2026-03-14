package no_sql_injection_taint_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_sql_injection_taint"
)

func TestNoSqlInjectionTaint(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_sql_injection_taint", &no_sql_injection_taint.NoSqlInjectionTaintRule{})
}
