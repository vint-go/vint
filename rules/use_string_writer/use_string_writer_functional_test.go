package use_string_writer_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_string_writer"
)

func TestUseStringWriter(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_string_writer", &use_string_writer.UseStringWriterRule{})
}
