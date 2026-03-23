package use_idiomatic_duration_name_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_idiomatic_duration_name"
)

func TestUseIdiomaticDurationName(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_idiomatic_duration_name", &use_idiomatic_duration_name.UseIdiomaticDurationNameRule{})
}
