package no_sql_injection_taint_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_sql_injection_taint"
)

func TestNoSqlInjectionTaint(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_sql_injection_taint", &no_sql_injection_taint.NoSqlInjectionTaintRule{})
}
