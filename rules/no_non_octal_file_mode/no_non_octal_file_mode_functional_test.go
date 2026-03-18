package no_non_octal_file_mode_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_non_octal_file_mode"
)

func TestNoNonOctalFileMode(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_non_octal_file_mode", &no_non_octal_file_mode.NoNonOctalFileModeRule{})
}
