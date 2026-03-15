package no_sql_stmt_exec_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_sql_stmt_exec"
)

func TestNoSqlStmtExec(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_sql_stmt_exec", &no_sql_stmt_exec.NoSqlStmtExecRule{})
}
