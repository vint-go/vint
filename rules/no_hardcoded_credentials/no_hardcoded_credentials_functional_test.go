package no_hardcoded_credentials_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_hardcoded_credentials"
)

func TestNoHardcodedCredentials(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_hardcoded_credentials", &no_hardcoded_credentials.NoHardcodedCredentialsRule{})
}
