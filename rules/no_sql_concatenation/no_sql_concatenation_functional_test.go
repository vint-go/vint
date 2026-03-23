package no_sql_concatenation_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_sql_concatenation"
)

func TestNoSqlConcatenation(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_sql_concatenation", &no_sql_concatenation.NoSqlConcatenationRule{})
}
