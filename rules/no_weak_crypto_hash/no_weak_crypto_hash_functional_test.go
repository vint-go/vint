package no_weak_crypto_hash_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_weak_crypto_hash"
)

func TestNoWeakCryptoHash(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_weak_crypto_hash", &no_weak_crypto_hash.NoWeakCryptoHashRule{})
}
