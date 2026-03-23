package no_self_referencing_finalizer_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_self_referencing_finalizer"
)

func TestNoSelfReferencingFinalizer(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_self_referencing_finalizer", &no_self_referencing_finalizer.NoSelfReferencingFinalizerRule{})
}
