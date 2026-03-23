package no_frame_pointer_clobber_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_frame_pointer_clobber"
)

func TestNoFramePointerClobber(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_frame_pointer_clobber", &no_frame_pointer_clobber.NoFramePointerClobberRule{})
}
