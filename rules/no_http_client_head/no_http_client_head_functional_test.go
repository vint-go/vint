package no_http_client_head_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_http_client_head"
)

func TestNoHttpClientHead(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_http_client_head", &no_http_client_head.NoHttpClientHeadRule{})
}
