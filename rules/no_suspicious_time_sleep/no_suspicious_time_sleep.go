package no_suspicious_time_sleep

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoSuspiciousTimeSleepRule detects suspiciously small untyped integer constants
// passed to time.Sleep, which almost certainly indicate a bug (sleeping for
// nanoseconds instead of the intended duration).
type NoSuspiciousTimeSleepRule struct{}

// Apply applies the rule to given file.
func (r *NoSuspiciousTimeSleepRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}
	w := &lintSuspiciousTimeSleep{onFailure: onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoSuspiciousTimeSleepRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}
	w := &lintSuspiciousTimeSleep{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoSuspiciousTimeSleepRule) Name() string {
	return "noSuspiciousTimeSleep"
}

// Group returns the rule group.
func (*NoSuspiciousTimeSleepRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoSuspiciousTimeSleepRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintSuspiciousTimeSleep struct {
	onFailure func(lint.Failure)
}

func (w *lintSuspiciousTimeSleep) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "time", "Sleep") {
		return w
	}

	if len(ce.Args) != 1 {
		return w
	}

	arg := ce.Args[0]
	lit, ok := arg.(*ast.BasicLit)
	if !ok {
		return w
	}

	if lit.Kind != token.INT {
		return w
	}

	val, err := strconv.ParseInt(lit.Value, 0, 64)
	if err != nil {
		return w
	}

	// Flag small integer literals (0 through 1000) as suspicious.
	// These are untyped constants that resolve to nanoseconds, which is
	// almost never the intended behavior.
	if val >= 0 && val <= 1000 {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryLogic,
			Failure:    fmt.Sprintf("suspiciously small untyped constant %s in time.Sleep", lit.Value),
		})
	}

	return w
}
