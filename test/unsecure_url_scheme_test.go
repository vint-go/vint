package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestUnsecureURLScheme(t *testing.T) {
	testRule(t, "unsecure_url_scheme", &rule.UnsecureURLSchemeRule{}, &lint.RuleConfig{})
}
