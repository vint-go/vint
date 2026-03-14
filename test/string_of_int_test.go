package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestStringOfInt(t *testing.T) {
	testRule(t, "string_of_int", &rule.StringOfIntRule{})
}
