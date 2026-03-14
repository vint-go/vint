package no_unclosed_bodies_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unclosed_bodies"
)

func TestNoUnclosedBodies(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unclosed_bodies", &no_unclosed_bodies.NoUnclosedBodiesRule{})
}
