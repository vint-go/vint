package no_forbidden_call_in_wg_go_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_forbidden_call_in_wg_go"
)

func TestNoForbiddenCallInWgGo(t *testing.T) {
	functional_test_helpers.TestRule(t, "go1.25/no_forbidden_call_in_wg_go", &no_forbidden_call_in_wg_go.ForbiddenCallInWgGoRule{})
}
