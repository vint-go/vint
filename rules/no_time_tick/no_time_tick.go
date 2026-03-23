package no_time_tick

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoTimeTickRule detects calls to time.Tick in functions other than main() or init(),
// which leak the underlying ticker since it can never be stopped or garbage collected.
type NoTimeTickRule struct{}

// Apply applies the rule to given file.
func (r *NoTimeTickRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintTimeTick{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoTimeTickRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintTimeTick{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoTimeTickRule) Name() string {
	return "noTimeTick"
}

// Group returns the rule group.
func (*NoTimeTickRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoTimeTickRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintTimeTick struct {
	onFailure func(lint.Failure)
}

func (w *lintTimeTick) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		// main() and init() are long-lived, so time.Tick is acceptable there
		if n.Name != nil && (n.Name.Name == "main" || n.Name.Name == "init") {
			return nil // skip main and init entirely
		}
		// For other named functions, walk their bodies looking for time.Tick
		return w

	case *ast.CallExpr:
		if astutils.IsPkgDotName(n.Fun, "time", "Tick") {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       n,
				Category:   lint.FailureCategoryBadPractice,
				Failure:    "calling time.Tick leaks the underlying ticker; use time.NewTicker and call Stop() when done",
			})
		}
		return w
	}

	return w
}
