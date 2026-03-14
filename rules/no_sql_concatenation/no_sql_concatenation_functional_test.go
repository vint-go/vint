package no_sql_concatenation_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_sql_concatenation"
)

func TestNoSqlConcatenation(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_sql_concatenation", &no_sql_concatenation.NoSqlConcatenationRule{})
}
