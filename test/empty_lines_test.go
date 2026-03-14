package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestEmptyLines(t *testing.T) {
	testRule(t, "empty_lines", &rule.EmptyLinesRule{})
}
