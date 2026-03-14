package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestModifiesParam(t *testing.T) {
	testRule(t, "modifies_param", &rule.ModifiesParamRule{})
}
