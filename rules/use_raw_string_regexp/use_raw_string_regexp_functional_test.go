package use_raw_string_regexp_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_raw_string_regexp"
)

func TestUseRawStringRegexp(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_raw_string_regexp", &use_raw_string_regexp.UseRawStringRegexpRule{})
}
