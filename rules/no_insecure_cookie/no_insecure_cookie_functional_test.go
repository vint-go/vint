package no_insecure_cookie_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_insecure_cookie"
)

func TestNoInsecureCookie(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_insecure_cookie", &no_insecure_cookie.NoInsecureCookieRule{})
}
