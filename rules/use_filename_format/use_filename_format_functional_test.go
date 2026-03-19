package use_filename_format_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_filename_format"
)

func TestUseFilenameFormatDefault(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_filename_format", &use_filename_format.FilenameFormatRule{})
}

func TestUseFilenameFormatNonASCII(t *testing.T) {
	functional_test_helpers.TestRule(t, "filenamе_with_non_ascii_char", &use_filename_format.FilenameFormatRule{})
}

func TestUseFilenameFormatWithUnderscores(t *testing.T) {
	functional_test_helpers.TestRule(t, "filename_with_underscores", &use_filename_format.FilenameFormatRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"^[A-Za-z][A-Za-z0-9]*.go$"},
	})
}
