package no_filesystem_root_serving_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_filesystem_root_serving"
)

func TestNoFilesystemRootServing(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_filesystem_root_serving", &no_filesystem_root_serving.NoFilesystemRootServingRule{})
}
