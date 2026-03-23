package use_epoch_naming_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_epoch_naming"
)

func TestUseEpochNaming(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_epoch_naming", &use_epoch_naming.EpochNamingRule{})
}
