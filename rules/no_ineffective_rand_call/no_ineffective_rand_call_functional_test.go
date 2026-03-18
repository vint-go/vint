package no_ineffective_rand_call_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_ineffective_rand_call"
)

func TestNoIneffectiveRandCall(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_ineffective_rand_call", &no_ineffective_rand_call.NoIneffectiveRandCallRule{})
}
