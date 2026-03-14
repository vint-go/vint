package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestConfusingResults(t *testing.T) {
	testRule(t, "confusing_results", &rule.ConfusingResultsRule{})
}
