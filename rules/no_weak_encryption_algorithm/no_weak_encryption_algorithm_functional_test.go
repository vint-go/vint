package no_weak_encryption_algorithm_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_weak_encryption_algorithm"
)

func TestNoWeakEncryptionAlgorithm(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_weak_encryption_algorithm", &no_weak_encryption_algorithm.NoWeakEncryptionAlgorithmRule{})
}
