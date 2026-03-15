package no_weak_encryption_algorithm

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoWeakEncryptionAlgorithmRule detects usage of DES or RC4 encryption algorithms.
type NoWeakEncryptionAlgorithmRule struct{}

// Apply applies the rule to given file.
func (r *NoWeakEncryptionAlgorithmRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintWeakEncryption{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoWeakEncryptionAlgorithmRule) Name() string {
	return "noWeakEncryptionAlgorithm"
}

// Group returns the rule group.
func (*NoWeakEncryptionAlgorithmRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoWeakEncryptionAlgorithmRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintWeakEncryption struct {
	onFailure func(lint.Failure)
}

func (w *lintWeakEncryption) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if astutils.IsPkgDotName(ce.Fun, "des", "NewCipher") {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "use of weak encryption algorithm DES: use AES instead",
		})
		return w
	}

	if astutils.IsPkgDotName(ce.Fun, "des", "NewTripleDESCipher") {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "use of weak encryption algorithm 3DES: use AES instead",
		})
		return w
	}

	if astutils.IsPkgDotName(ce.Fun, "rc4", "NewCipher") {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "use of weak encryption algorithm RC4: use AES instead",
		})
		return w
	}

	return w
}
