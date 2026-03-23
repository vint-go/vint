package no_sql_tx_exec_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_sql_tx_exec"
)

func TestNoSqlTxExec(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_sql_tx_exec", &no_sql_tx_exec.NoSqlTxExecRule{})
}
