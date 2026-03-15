package no_sql_stmt_query_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_sql_stmt_query"
)

func TestNoSqlStmtQuery(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_sql_stmt_query", &no_sql_stmt_query.NoSqlStmtQueryRule{})
}
