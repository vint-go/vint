package no_http_head_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_http_head"
)

func TestNoHttpHead(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_http_head", &no_http_head.NoHttpHeadRule{})
}
