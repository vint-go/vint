package no_invalid_regexp_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_invalid_regexp"
)

func TestNoInvalidRegexp(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_invalid_regexp", &no_invalid_regexp.NoInvalidRegexpRule{})
}
