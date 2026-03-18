package no_select_break_confusion_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_select_break_confusion"
)

func TestNoSelectBreakConfusion(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_select_break_confusion", &no_select_break_confusion.NoSelectBreakConfusionRule{})
}
