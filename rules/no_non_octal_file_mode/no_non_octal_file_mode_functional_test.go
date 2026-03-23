package no_non_octal_file_mode_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_non_octal_file_mode"
)

func TestNoNonOctalFileMode(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_non_octal_file_mode", &no_non_octal_file_mode.NoNonOctalFileModeRule{})
}
