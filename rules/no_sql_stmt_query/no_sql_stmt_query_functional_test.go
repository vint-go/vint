package no_sql_stmt_query_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_sql_stmt_query"
)

func TestNoSqlStmtQuery(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_sql_stmt_query", &no_sql_stmt_query.NoSqlStmtQueryRule{})
}
