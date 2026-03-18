package use_copy_for_slide_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_copy_for_slide"
)

func TestUseCopyForSlide(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_copy_for_slide", &use_copy_for_slide.UseCopyForSlideRule{})
}
