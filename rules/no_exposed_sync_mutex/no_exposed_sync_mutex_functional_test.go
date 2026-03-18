package no_exposed_sync_mutex_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_exposed_sync_mutex"
)

func TestNoExposedSyncMutex(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_exposed_sync_mutex", &no_exposed_sync_mutex.NoExposedSyncMutexRule{})
}
