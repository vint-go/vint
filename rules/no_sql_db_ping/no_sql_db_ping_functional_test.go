package no_sql_db_ping_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_sql_db_ping"
)

func TestNoSqlDbPing(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_sql_db_ping", &no_sql_db_ping.NoSqlDbPingRule{})
}
