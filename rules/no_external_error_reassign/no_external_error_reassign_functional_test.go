package no_external_error_reassign_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_external_error_reassign"
)

func TestNoExternalErrorReassign(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_external_error_reassign", &no_external_error_reassign.NoExternalErrorReassignRule{})
}
