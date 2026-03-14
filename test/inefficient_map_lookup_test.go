package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestInefficientMapLookup(t *testing.T) {
	testRule(t, "inefficient_map_lookup", &rule.InefficientMapLookupRule{}, &lint.RuleConfig{})
}
