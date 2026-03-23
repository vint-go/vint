package no_duplicate_cutset_chars

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoDuplicateCutsetCharsRule detects duplicate characters in cutset arguments
// passed to strings.TrimLeft or strings.TrimRight.
type NoDuplicateCutsetCharsRule struct{}

// Apply applies the rule to given file.
func (r *NoDuplicateCutsetCharsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintDuplicateCutset{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoDuplicateCutsetCharsRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintDuplicateCutset{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoDuplicateCutsetCharsRule) Name() string {
	return "noDuplicateCutsetChars"
}

// Group returns the rule group.
func (*NoDuplicateCutsetCharsRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoDuplicateCutsetCharsRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintDuplicateCutset struct {
	onFailure func(lint.Failure)
}

func (w *lintDuplicateCutset) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	isTrimLeft := astutils.IsPkgDotName(ce.Fun, "strings", "TrimLeft")
	isTrimRight := astutils.IsPkgDotName(ce.Fun, "strings", "TrimRight")

	if !isTrimLeft && !isTrimRight {
		return w
	}

	if len(ce.Args) != 2 {
		return w
	}

	cutsetArg, ok := ce.Args[1].(*ast.BasicLit)
	if !ok {
		return w
	}

	if cutsetArg.Kind != token.STRING {
		return w
	}

	// Extract the string value (remove surrounding quotes)
	val := cutsetArg.Value
	if len(val) >= 2 {
		val = val[1 : len(val)-1]
	}

	if hasDuplicateChars(val) {
		funcName := "strings.TrimLeft"
		suggested := "strings.TrimPrefix"
		if isTrimRight {
			funcName = "strings.TrimRight"
			suggested = "strings.TrimSuffix"
		}
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryLogic,
			Failure:    fmt.Sprintf("cutset %s has duplicate characters, did you mean to use %s instead of %s?", cutsetArg.Value, suggested, funcName),
		})
	}

	return w
}

// hasDuplicateChars returns true if the string contains any duplicate characters.
func hasDuplicateChars(s string) bool {
	seen := make(map[rune]bool, len(s))
	for _, ch := range s {
		if seen[ch] {
			return true
		}
		seen[ch] = true
	}
	return false
}
