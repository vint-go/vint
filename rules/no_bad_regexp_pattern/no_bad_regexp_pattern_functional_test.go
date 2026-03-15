package no_bad_regexp_pattern_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_bad_regexp_pattern"
)

func TestNoBadRegexpPattern(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_bad_regexp_pattern", &no_bad_regexp_pattern.NoBadRegexpPatternRule{})
}
