package use_http_no_body_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_http_no_body"
)

func TestUseHttpNoBody(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_http_no_body", &use_http_no_body.UseHttpNoBodyRule{})
}
