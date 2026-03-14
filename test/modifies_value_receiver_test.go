package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestModifiesValRec(t *testing.T) {
	testRule(t, "modifies_value_receiver", &rule.ModifiesValRecRule{})
}
