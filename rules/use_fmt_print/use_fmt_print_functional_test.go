package use_fmt_print_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_fmt_print"
)

func TestUseFmtPrint(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_fmt_print", &use_fmt_print.UseFmtPrintRule{})
}

func TestUseFmtPrintWithRedefinition(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_fmt_print_with_redefinition", &use_fmt_print.UseFmtPrintRule{})
}
