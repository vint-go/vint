package use_regexp_must_compile_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_regexp_must_compile"
)

func TestUseRegexpMustCompile(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_regexp_must_compile", &use_regexp_must_compile.UseRegexpMustCompileRule{})
}
