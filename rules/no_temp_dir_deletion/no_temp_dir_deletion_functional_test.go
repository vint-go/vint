package no_temp_dir_deletion_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_temp_dir_deletion"
)

func TestNoTempDirDeletion(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_temp_dir_deletion", &no_temp_dir_deletion.NoTempDirDeletionRule{})
}
