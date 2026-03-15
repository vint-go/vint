package no_conflicting_http_mux_patterns_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_conflicting_http_mux_patterns"
)

func TestNoConflictingHttpMuxPatterns(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_conflicting_http_mux_patterns", &no_conflicting_http_mux_patterns.NoConflictingHttpMuxPatternsRule{})
}
