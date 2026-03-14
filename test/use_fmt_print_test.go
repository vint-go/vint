package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestUseFmtPrint(t *testing.T) {
	testRule(t, "use_fmt_print", &rule.UseFmtPrintRule{})
	testRule(t, "use_fmt_print_with_redefinition", &rule.UseFmtPrintRule{})
}
