package no_unused_write_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unused_write"
)

func TestNoUnusedWrite(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unused_write", &no_unused_write.NoUnusedWriteRule{})
}
