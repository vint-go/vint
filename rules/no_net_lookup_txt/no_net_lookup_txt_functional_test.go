package no_net_lookup_txt_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_net_lookup_txt"
)

func TestNoNetLookupTxt(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_net_lookup_txt", &no_net_lookup_txt.NoNetLookupTxtRule{})
}
