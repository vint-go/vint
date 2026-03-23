package no_insecure_host_key_callback

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoInsecureHostKeyCallbackRule flags calls to ssh.InsecureIgnoreHostKey() which
// disables SSH host key verification, making connections vulnerable to MITM attacks.
type NoInsecureHostKeyCallbackRule struct{}

// Apply applies the rule to given file.
func (r *NoInsecureHostKeyCallbackRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintInsecureHostKey{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoInsecureHostKeyCallbackRule) Name() string {
	return "noInsecureHostKeyCallback"
}

// Group returns the rule group.
func (*NoInsecureHostKeyCallbackRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoInsecureHostKeyCallbackRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintInsecureHostKey struct {
	onFailure func(lint.Failure)
}

func (w *lintInsecureHostKey) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w // nothing to do, the node is not a function call
	}

	if !astutils.IsPkgDotName(ce.Fun, "ssh", "InsecureIgnoreHostKey") {
		return w // not a call to ssh.InsecureIgnoreHostKey
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "use of ssh.InsecureIgnoreHostKey disables host key verification and is vulnerable to MITM attacks",
	})

	return w
}
