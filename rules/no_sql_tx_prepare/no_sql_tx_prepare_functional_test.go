package no_sql_tx_prepare_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_sql_tx_prepare"
)

func TestNoSqlTxPrepare(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_sql_tx_prepare", &no_sql_tx_prepare.NoSqlTxPrepareRule{})
}
