package no_stdlib_version_mismatch_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_stdlib_version_mismatch"
)

func TestNoStdlibVersionMismatch(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_stdlib_version_mismatch", &no_stdlib_version_mismatch.NoStdlibVersionMismatchRule{})
}
