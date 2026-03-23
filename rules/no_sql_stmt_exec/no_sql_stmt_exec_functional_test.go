package no_sql_stmt_exec_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_sql_stmt_exec"
)

func TestNoSqlStmtExec(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_sql_stmt_exec", &no_sql_stmt_exec.NoSqlStmtExecRule{})
}
