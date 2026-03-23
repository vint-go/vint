package no_http_get_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_http_get"
)

func TestNoHttpGet(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_http_get", &no_http_get.NoHttpGetRule{})
}
