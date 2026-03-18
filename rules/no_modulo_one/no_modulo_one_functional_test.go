package no_modulo_one_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_modulo_one"
)

func TestNoModuloOne(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_modulo_one", &no_modulo_one.NoModuloOneRule{})
}
