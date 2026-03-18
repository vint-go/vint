package use_idiomatic_naming_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_idiomatic_naming"
)

func TestUseIdiomaticNaming(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_idiomatic_naming", &use_idiomatic_naming.UseIdiomaticNamingRule{})
}
