package no_separator_in_filepath_join_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_separator_in_filepath_join"
)

func TestNoSeparatorInFilepathJoin(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_separator_in_filepath_join", &no_separator_in_filepath_join.NoSeparatorInFilepathJoinRule{})
}
