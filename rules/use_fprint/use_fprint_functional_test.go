package use_fprint_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_fprint"
)

func TestUseFprint(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_fprint", &use_fprint.UseFprintRule{})
}
