package no_sql_tx_prepare_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_sql_tx_prepare"
)

func TestNoSqlTxPrepare(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_sql_tx_prepare", &no_sql_tx_prepare.NoSqlTxPrepareRule{})
}
