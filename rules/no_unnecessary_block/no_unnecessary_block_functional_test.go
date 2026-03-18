package no_unnecessary_block_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unnecessary_block"
)

func TestNoUnnecessaryBlock(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unnecessary_block", &no_unnecessary_block.NoUnnecessaryBlockRule{})
}
