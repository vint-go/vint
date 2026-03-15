package no_sql_db_query_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_sql_db_query"
)

func TestNoSqlDbQuery(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_sql_db_query", &no_sql_db_query.NoSqlDbQueryRule{})
}
