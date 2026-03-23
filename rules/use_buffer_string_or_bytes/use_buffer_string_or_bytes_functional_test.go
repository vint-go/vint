package use_buffer_string_or_bytes_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_buffer_string_or_bytes"
)

func TestUseBufferStringOrBytes(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_buffer_string_or_bytes", &use_buffer_string_or_bytes.UseBufferStringOrBytesRule{})
}
