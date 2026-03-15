package no_sql_tx_query_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_sql_tx_query"
)

func TestNoSqlTxQuery(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_sql_tx_query", &no_sql_tx_query.NoSqlTxQueryRule{})
}
