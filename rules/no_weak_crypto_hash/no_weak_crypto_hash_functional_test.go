package no_weak_crypto_hash_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_weak_crypto_hash"
)

func TestNoWeakCryptoHash(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_weak_crypto_hash", &no_weak_crypto_hash.NoWeakCryptoHashRule{})
}
