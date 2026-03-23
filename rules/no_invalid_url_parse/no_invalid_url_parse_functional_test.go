package no_invalid_url_parse_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_invalid_url_parse"
)

func TestNoInvalidUrlParse(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_invalid_url_parse", &no_invalid_url_parse.NoInvalidUrlParseRule{})
}
