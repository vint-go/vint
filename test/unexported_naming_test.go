package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestUnexportedNaming(t *testing.T) {
	testRule(t, "unexported_naming", &rule.UnexportedNamingRule{})
}
