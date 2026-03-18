package use_filepath_join_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_filepath_join"
)

func TestUseFilepathJoin(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_filepath_join", &use_filepath_join.UseFilepathJoinRule{})
}
