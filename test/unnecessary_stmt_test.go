package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestUnnecessaryStmt(t *testing.T) {
	testRule(t, "unnecessary_stmt", &rule.UnnecessaryStmtRule{})
}
