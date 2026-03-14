package no_insecure_random_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_insecure_random"
)

func TestNoInsecureRandom(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_insecure_random", &no_insecure_random.NoInsecureRandomRule{})
}
