package use_compiled_regexp_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_compiled_regexp"
)

func TestUseCompiledRegexp(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_compiled_regexp", &use_compiled_regexp.UseCompiledRegexpRule{})
}
