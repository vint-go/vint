package use_simplified_regexp_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_simplified_regexp"
)

func TestUseSimplifiedRegexp(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_simplified_regexp", &use_simplified_regexp.UseSimplifiedRegexpRule{})
}
