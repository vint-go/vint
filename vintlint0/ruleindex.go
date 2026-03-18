package vintlint0

import "github.com/strowk/vint/lint"

// RuleIndex partitions rules by their requirements for efficient dispatch
// to separate worker pools.
//
// Currently indexes by RequiresTypecheck. The struct is designed to be
// extended with additional indexes (e.g., ByNodeType, RequiresPackage)
// in the future.
type RuleIndex struct {
	// ASTOnly contains rules that do NOT require type checking.
	// These can run immediately after AST parsing with no lock contention.
	ASTOnly []lint.Rule

	// TypeCheck contains rules that require type information.
	// These must wait for Package.TypeCheck() to complete before running.
	TypeCheck []lint.Rule
}

// NewRuleIndex builds a RuleIndex by partitioning the given rules based on
// whether they implement RequiresTypecheck() bool returning true.
func NewRuleIndex(rules []lint.Rule) *RuleIndex {
	idx := &RuleIndex{
		ASTOnly:   make([]lint.Rule, 0, len(rules)),
		TypeCheck: make([]lint.Rule, 0),
	}
	for _, r := range rules {
		if requiresTypecheck(r) {
			idx.TypeCheck = append(idx.TypeCheck, r)
		} else {
			idx.ASTOnly = append(idx.ASTOnly, r)
		}
	}
	return idx
}

// typecheckAware is the interface that rules can implement to declare they
// need type information. This is checked via type assertion rather than
// being a formal part of the Rule interface.
type typecheckAware interface {
	RequiresTypecheck() bool
}

func requiresTypecheck(r lint.Rule) bool {
	tc, ok := r.(typecheckAware)
	return ok && tc.RequiresTypecheck()
}
