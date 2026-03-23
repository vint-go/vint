package no_zip_slip_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_zip_slip"
)

func TestNoZipSlip(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_zip_slip", &no_zip_slip.NoZipSlipRule{})
}
