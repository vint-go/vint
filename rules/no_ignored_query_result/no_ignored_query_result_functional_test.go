package no_ignored_query_result_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_ignored_query_result"
)

func TestNoIgnoredQueryResult(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_ignored_query_result", &no_ignored_query_result.NoIgnoredQueryResultRule{})
}
