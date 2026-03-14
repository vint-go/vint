package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestUseErrorsNew(t *testing.T) {
	testRule(t, "use_errors_new", &rule.UseErrorsNewRule{})
}
