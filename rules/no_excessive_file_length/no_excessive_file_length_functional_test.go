package no_excessive_file_length_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_excessive_file_length"
)

func TestNoExcessiveFileLength(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_excessive_file_length_disabled", &no_excessive_file_length.NoExcessiveFileLengthRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{},
	})
	functional_test_helpers.TestRule(t, "no_excessive_file_length_disabled", &no_excessive_file_length.NoExcessiveFileLengthRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{"max": int64(0)}},
	})
	functional_test_helpers.TestRule(t, "no_excessive_file_length_disabled", &no_excessive_file_length.NoExcessiveFileLengthRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{"skipComments": true, "skipBlankLines": true}},
	})
	functional_test_helpers.TestRule(t, "no_excessive_file_length_9", &no_excessive_file_length.NoExcessiveFileLengthRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{"max": int64(9)}},
	})
	functional_test_helpers.TestRule(t, "no_excessive_file_length_7_skip_comments", &no_excessive_file_length.NoExcessiveFileLengthRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{"max": int64(7), "skipComments": true}},
	})
	functional_test_helpers.TestRule(t, "no_excessive_file_length_6_skip_blank", &no_excessive_file_length.NoExcessiveFileLengthRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{"max": int64(6), "skipBlankLines": true}},
	})
	functional_test_helpers.TestRule(t, "no_excessive_file_length_4_skip_comments_skip_blank", &no_excessive_file_length.NoExcessiveFileLengthRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{"max": int64(4), "skipComments": true, "skipBlankLines": true}},
	})
	functional_test_helpers.TestRule(t, "no_excessive_file_length_4_skip_comments_skip_blank", &no_excessive_file_length.NoExcessiveFileLengthRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{"max": int64(4), "skip-comments": true, "skip-blank-lines": true}},
	})
}
