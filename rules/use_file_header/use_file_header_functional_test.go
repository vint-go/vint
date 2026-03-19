package use_file_header_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_file_header"
)

func TestUseFileHeaderDefault(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_file_header_default", &use_file_header.FileHeaderRule{}, &lint.RuleConfig{})
}

func TestUseFileHeader(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_file_header1", &use_file_header.FileHeaderRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"foobar"},
	})

	functional_test_helpers.TestRule(t, "use_file_header2", &use_file_header.FileHeaderRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"foobar"},
	})

	functional_test_helpers.TestRule(t, "use_file_header3", &use_file_header.FileHeaderRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"foobar"},
	})

	functional_test_helpers.TestRule(t, "use_file_header4", &use_file_header.FileHeaderRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"^\\sfoobar$"},
	})

	functional_test_helpers.TestRule(t, "use_file_header5", &use_file_header.FileHeaderRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"^\\sfoo.*bar$"},
	})

	functional_test_helpers.TestRule(t, "use_file_header6", &use_file_header.FileHeaderRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"foobar"},
	})
}

func BenchmarkUseFileHeader(b *testing.B) {
	for b.Loop() {
		functional_test_helpers.TestRule(b, "use_file_header1", &use_file_header.FileHeaderRule{}, &lint.RuleConfig{
			Arguments: lint.Arguments{"foobar"},
		})
	}
}
