package no_untrappable_signal

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoUntrappableSignalRule detects attempts to trap signals that cannot be
// intercepted by a program, such as SIGKILL and SIGSTOP.
type NoUntrappableSignalRule struct{}

// Apply applies the rule to given file.
func (r *NoUntrappableSignalRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUntrappableSignal{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoUntrappableSignalRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintUntrappableSignal{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoUntrappableSignalRule) Name() string {
	return "noUntrappableSignal"
}

// Group returns the rule group.
func (*NoUntrappableSignalRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoUntrappableSignalRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintUntrappableSignal struct {
	onFailure func(lint.Failure)
}

// untrappableSignals maps signal names that cannot be caught or ignored.
var untrappableSignals = map[string]bool{
	"SIGKILL": true,
	"SIGSTOP": true,
}

func (w *lintUntrappableSignal) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check for signal.Notify or signal.Ignore calls.
	isNotify := astutils.IsPkgDotName(ce.Fun, "signal", "Notify")
	isIgnore := astutils.IsPkgDotName(ce.Fun, "signal", "Ignore")
	if !isNotify && !isIgnore {
		return w
	}

	funcName := "signal.Notify"
	if isIgnore {
		funcName = "signal.Ignore"
	}

	// For signal.Notify, the first argument is the channel; signals start at index 1.
	// For signal.Ignore, all arguments are signals starting at index 0.
	startIdx := 0
	if isNotify {
		startIdx = 1
	}

	for i := startIdx; i < len(ce.Args); i++ {
		arg := ce.Args[i]
		for sigName := range untrappableSignals {
			if astutils.IsPkgDotName(arg, "syscall", sigName) || astutils.IsPkgDotName(arg, "os", sigName) {
				w.onFailure(lint.Failure{
					Confidence: 1,
					Node:       arg,
					Category:   lint.FailureCategoryLogic,
					Failure:    sigName + " cannot be trapped by " + funcName,
				})
			}
		}
	}

	return w
}
