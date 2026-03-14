package no_predictable_temp_file_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_predictable_temp_file"
)

func TestNoPredictableTempFile(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_predictable_temp_file", &no_predictable_temp_file.NoPredictableTempFileRule{})
}
