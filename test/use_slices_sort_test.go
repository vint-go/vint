package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestUseSlicesSort(t *testing.T) {
	testRule(t, "use_slices_sort", &rule.UseSlicesSort{})
	testRule(t, "go1.21/use_slices_sort", &rule.UseSlicesSort{})
}
