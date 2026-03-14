package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestLintFilenameFormat(t *testing.T) {
	testRule(t, "filename_ok_default", &rule.FilenameFormatRule{}, &lint.RuleConfig{})

	testRule(t, "filenamе_with_non_ascii_char", &rule.FilenameFormatRule{}, &lint.RuleConfig{})

	testRule(t, "filename_with_underscores", &rule.FilenameFormatRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"^[A-Za-z][A-Za-z0-9]*.go$"}})
}
