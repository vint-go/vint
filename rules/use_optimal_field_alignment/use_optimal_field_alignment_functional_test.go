package use_optimal_field_alignment_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_optimal_field_alignment"
)

func TestUseOptimalFieldAlignment(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_optimal_field_alignment", &use_optimal_field_alignment.UseOptimalFieldAlignmentRule{})
}
