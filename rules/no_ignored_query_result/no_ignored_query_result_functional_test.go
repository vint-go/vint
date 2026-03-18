package no_ignored_query_result_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_ignored_query_result"
)

func TestNoIgnoredQueryResult(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_ignored_query_result", &no_ignored_query_result.NoIgnoredQueryResultRule{})
}
