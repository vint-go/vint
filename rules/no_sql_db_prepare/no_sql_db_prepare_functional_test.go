package no_sql_db_prepare_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_sql_db_prepare"
)

func TestNoSqlDbPrepare(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_sql_db_prepare", &no_sql_db_prepare.NoSqlDbPrepareRule{})
}
