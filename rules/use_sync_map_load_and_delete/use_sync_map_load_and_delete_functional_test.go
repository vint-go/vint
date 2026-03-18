package use_sync_map_load_and_delete_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_sync_map_load_and_delete"
)

func TestUseSyncMapLoadAndDelete(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_sync_map_load_and_delete", &use_sync_map_load_and_delete.UseSyncMapLoadAndDeleteRule{})
}
