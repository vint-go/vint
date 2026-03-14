package no_exposed_pprof_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_exposed_pprof"
)

func TestNoExposedPprof(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_exposed_pprof", &no_exposed_pprof.NoExposedPprofRule{})
}
