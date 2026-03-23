package no_http_head_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_http_head"
)

func TestNoHttpHead(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_http_head", &no_http_head.NoHttpHeadRule{})
}
