package no_separator_in_filepath_join

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoSeparatorInFilepathJoinRule detects string literal arguments to filepath.Join
// that contain path separators (forward slashes or backslashes).
type NoSeparatorInFilepathJoinRule struct{}

// Apply applies the rule to given file.
func (r *NoSeparatorInFilepathJoinRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintSeparatorInFilepathJoin{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoSeparatorInFilepathJoinRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintSeparatorInFilepathJoin{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoSeparatorInFilepathJoinRule) Name() string {
	return "noSeparatorInFilepathJoin"
}

// Group returns the rule group.
func (*NoSeparatorInFilepathJoinRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoSeparatorInFilepathJoinRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintSeparatorInFilepathJoin struct {
	onFailure func(lint.Failure)
}

func (w *lintSeparatorInFilepathJoin) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "filepath", "Join") {
		return w
	}

	for _, arg := range ce.Args {
		lit, ok := arg.(*ast.BasicLit)
		if !ok {
			continue
		}
		if lit.Kind != token.STRING {
			continue
		}

		// Extract the string value (remove quotes)
		val := lit.Value
		if len(val) >= 2 {
			val = val[1 : len(val)-1]
		}

		if strings.ContainsAny(val, "/\\") {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       ce,
				Category:   lint.FailureCategoryStyle,
				Failure:    fmt.Sprintf("filepath.Join argument %s contains a path separator", lit.Value),
			})
			break
		}
	}

	return w
}
