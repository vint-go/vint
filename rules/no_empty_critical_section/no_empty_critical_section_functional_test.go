package no_empty_critical_section_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_empty_critical_section"
)

func TestNoEmptyCriticalSection(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_empty_critical_section", &no_empty_critical_section.NoEmptyCriticalSectionRule{})
}
