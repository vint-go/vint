package no_weak_crypto_hash

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoWeakCryptoHashRule detects the usage of weak cryptographic hash functions MD5 or SHA1.
type NoWeakCryptoHashRule struct{}

// Apply applies the rule to given file.
func (r *NoWeakCryptoHashRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintWeakCryptoHash{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoWeakCryptoHashRule) Name() string {
	return "noWeakCryptoHash"
}

// Group returns the rule group.
func (*NoWeakCryptoHashRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoWeakCryptoHashRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintWeakCryptoHash struct {
	onFailure func(lint.Failure)
}

func (w *lintWeakCryptoHash) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if astutils.IsPkgDotName(ce.Fun, "md5", "New") || astutils.IsPkgDotName(ce.Fun, "md5", "Sum") {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "use of weak cryptographic hash function md5",
		})
		return w
	}

	if astutils.IsPkgDotName(ce.Fun, "sha1", "New") || astutils.IsPkgDotName(ce.Fun, "sha1", "Sum") {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "use of weak cryptographic hash function sha1",
		})
		return w
	}

	return w
}
