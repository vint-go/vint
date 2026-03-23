package no_unused_write_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unused_write"
)

func TestNoUnusedWrite(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unused_write", &no_unused_write.NoUnusedWriteRule{})
}
