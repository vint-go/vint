package no_conflicting_http_mux_patterns_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_conflicting_http_mux_patterns"
)

func TestNoConflictingHttpMuxPatterns(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_conflicting_http_mux_patterns", &no_conflicting_http_mux_patterns.NoConflictingHttpMuxPatternsRule{})
}
