package no_sql_format_string_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_sql_format_string"
)

func TestNoSqlFormatString(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_sql_format_string", &no_sql_format_string.NoSqlFormatStringRule{})
}
