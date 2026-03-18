package no_writer_buffer_modification_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_writer_buffer_modification"
)

func TestNoWriterBufferModification(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_writer_buffer_modification", &no_writer_buffer_modification.NoWriterBufferModificationRule{})
}
