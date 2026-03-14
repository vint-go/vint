package no_insecure_tls_config_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_insecure_tls_config"
)

func TestNoInsecureTlsConfig(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_insecure_tls_config", &no_insecure_tls_config.NoInsecureTlsConfigRule{})
}
