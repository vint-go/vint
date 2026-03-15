package no_sql_db_ping_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_sql_db_ping"
)

func TestNoSqlDbPing(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_sql_db_ping", &no_sql_db_ping.NoSqlDbPingRule{})
}
