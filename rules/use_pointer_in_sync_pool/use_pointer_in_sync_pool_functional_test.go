package use_pointer_in_sync_pool_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_pointer_in_sync_pool"
)

func TestUsePointerInSyncPool(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_pointer_in_sync_pool", &use_pointer_in_sync_pool.UsePointerInSyncPoolRule{})
}
