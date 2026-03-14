package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestNestedStructs(t *testing.T) {
	testRule(t, "nested_structs", &rule.NestedStructs{}, &lint.RuleConfig{})
}
