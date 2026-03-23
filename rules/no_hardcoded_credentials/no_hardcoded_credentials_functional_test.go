package no_hardcoded_credentials_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_hardcoded_credentials"
)

func TestNoHardcodedCredentials(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_hardcoded_credentials", &no_hardcoded_credentials.NoHardcodedCredentialsRule{})
}

func TestNoHardcodedCredentialsCustomConfig(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_hardcoded_credentials_custom", &no_hardcoded_credentials.NoHardcodedCredentialsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{
			"pattern":          "(?i)(my_custom_secret|my_custom_token)",
			"entropyThreshold": float64(4.0),
		}},
	})
}
