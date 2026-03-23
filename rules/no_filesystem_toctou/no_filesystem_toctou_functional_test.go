package no_filesystem_toctou_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_filesystem_toctou"
)

func TestNoFilesystemToctou(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_filesystem_toctou", &no_filesystem_toctou.NoFilesystemToctouRule{})
}
