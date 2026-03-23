package no_sql_tx_query_row_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_sql_tx_query_row"
)

func TestNoSqlTxQueryRow(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_sql_tx_query_row", &no_sql_tx_query_row.NoSqlTxQueryRowRule{})
}
