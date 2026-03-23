package no_unsafe_redirect_policy_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unsafe_redirect_policy"
)

func TestNoUnsafeRedirectPolicy(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unsafe_redirect_policy", &no_unsafe_redirect_policy.NoUnsafeRedirectPolicyRule{})
}
