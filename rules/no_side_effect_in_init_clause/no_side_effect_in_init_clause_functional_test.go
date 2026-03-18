package no_side_effect_in_init_clause_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_side_effect_in_init_clause"
)

func TestNoSideEffectInInitClause(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_side_effect_in_init_clause", &no_side_effect_in_init_clause.NoSideEffectInInitClauseRule{})
}
