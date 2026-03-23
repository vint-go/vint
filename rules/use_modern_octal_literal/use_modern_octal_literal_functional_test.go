package use_modern_octal_literal_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_modern_octal_literal"
)

func TestUseModernOctalLiteral(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_modern_octal_literal", &use_modern_octal_literal.UseModernOctalLiteralRule{})
}
