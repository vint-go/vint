package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestEpochNaming(t *testing.T) {
	testRule(t, "epoch_naming", &rule.EpochNamingRule{})
}
