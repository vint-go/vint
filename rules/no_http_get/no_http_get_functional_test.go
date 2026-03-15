package no_http_get_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_http_get"
)

func TestNoHttpGet(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_http_get", &no_http_get.NoHttpGetRule{})
}
