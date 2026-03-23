package no_exposed_sync_mutex_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_exposed_sync_mutex"
)

func TestNoExposedSyncMutex(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_exposed_sync_mutex", &no_exposed_sync_mutex.NoExposedSyncMutexRule{})
}
