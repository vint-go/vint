package no_sql_tx_stmt_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_sql_tx_stmt"
)

func TestNoSqlTxStmt(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_sql_tx_stmt", &no_sql_tx_stmt.NoSqlTxStmtRule{})
}
