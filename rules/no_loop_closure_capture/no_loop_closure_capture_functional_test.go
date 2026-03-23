package no_loop_closure_capture_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_loop_closure_capture"
)

func TestNoLoopClosureCapture(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_loop_closure_capture", &no_loop_closure_capture.NoLoopClosureCaptureRule{})
}
