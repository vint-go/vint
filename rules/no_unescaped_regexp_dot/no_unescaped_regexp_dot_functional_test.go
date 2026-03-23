package no_unescaped_regexp_dot_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unescaped_regexp_dot"
)

func TestNoUnescapedRegexpDot(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unescaped_regexp_dot", &no_unescaped_regexp_dot.NoUnescapedRegexpDotRule{})
}
