package no_banned_characters

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoBannedCharactersRule checks if identifiers contain banned characters.
type NoBannedCharactersRule struct {
	bannedCharList []string
}

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configure implements the [lint.ConfigurableRule] interface.
func (r *NoBannedCharactersRule) Configure(arguments lint.Arguments) error {
	r.bannedCharList = nil
	for _, arg := range arguments {
		charStr, ok := arg.(string)
		if !ok {
			return fmt.Errorf("invalid argument for the %s rule: expecting a string, got %T", r.Name(), arg)
		}
		r.bannedCharList = append(r.bannedCharList, charStr)
	}
	return nil
}

// Apply applies the rule to given file.
func (r *NoBannedCharactersRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if len(r.bannedCharList) == 0 {
		return failures
	}

	w := &lintBannedCharsRule{
		bannedChars: r.bannedCharList,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (*NoBannedCharactersRule) Name() string {
	return "noBannedCharacters"
}

// Group returns the rule group.
func (*NoBannedCharactersRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoBannedCharactersRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintBannedCharsRule struct {
	bannedChars []string
	onFailure   func(lint.Failure)
}

// Visit checks for each node if an identifier contains banned characters.
func (w *lintBannedCharsRule) Visit(node ast.Node) ast.Visitor {
	n, ok := node.(*ast.Ident)
	if !ok {
		return w
	}
	for _, c := range w.bannedChars {
		if strings.Contains(n.Name, c) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Failure:    fmt.Sprintf("banned character found: %s", c),
				Node:       n,
				Category:   lint.FailureCategoryNaming,
			})
		}
	}

	return w
}
