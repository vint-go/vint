package no_temp_dir_deletion_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_temp_dir_deletion"
)

func TestNoTempDirDeletion(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_temp_dir_deletion", &no_temp_dir_deletion.NoTempDirDeletionRule{})
}
