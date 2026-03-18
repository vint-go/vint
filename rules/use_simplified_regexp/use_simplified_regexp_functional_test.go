package use_simplified_regexp_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_simplified_regexp"
)

func TestUseSimplifiedRegexp(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_simplified_regexp", &use_simplified_regexp.UseSimplifiedRegexpRule{})
}
