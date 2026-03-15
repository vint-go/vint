package no_zip_slip_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_zip_slip"
)

func TestNoZipSlip(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_zip_slip", &no_zip_slip.NoZipSlipRule{})
}
