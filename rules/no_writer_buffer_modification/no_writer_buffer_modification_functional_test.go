package no_writer_buffer_modification_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_writer_buffer_modification"
)

func TestNoWriterBufferModification(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_writer_buffer_modification", &no_writer_buffer_modification.NoWriterBufferModificationRule{})
}
