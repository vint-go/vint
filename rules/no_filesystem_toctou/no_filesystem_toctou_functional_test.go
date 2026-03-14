package no_filesystem_toctou_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_filesystem_toctou"
)

func TestNoFilesystemToctou(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_filesystem_toctou", &no_filesystem_toctou.NoFilesystemToctouRule{})
}
