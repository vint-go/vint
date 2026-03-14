package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestConfusingNaming(t *testing.T) {
	testRule(t, "confusing_naming1", &rule.ConfusingNamingRule{})
}
