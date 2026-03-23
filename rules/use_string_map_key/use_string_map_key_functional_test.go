package use_string_map_key_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_string_map_key"
)

func TestUseStringMapKey(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_string_map_key", &use_string_map_key.UseStringMapKeyRule{})
}
