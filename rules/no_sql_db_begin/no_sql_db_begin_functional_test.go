package no_sql_db_begin_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_sql_db_begin"
)

func TestNoSqlDbBegin(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_sql_db_begin", &no_sql_db_begin.NoSqlDbBeginRule{})
}
