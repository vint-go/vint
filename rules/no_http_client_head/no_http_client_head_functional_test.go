package no_http_client_head_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_http_client_head"
)

func TestNoHttpClientHead(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_http_client_head", &no_http_client_head.NoHttpClientHeadRule{})
}
