package no_unnecessary_stmt_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unnecessary_stmt"
)

func TestNoUnnecessaryStmt(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unnecessary_stmt", &no_unnecessary_stmt.UnnecessaryStmtRule{})
}
