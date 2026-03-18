package use_write_byte_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_write_byte"
)

func TestUseWriteByte(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_write_byte", &use_write_byte.UseWriteByteRule{})
}
