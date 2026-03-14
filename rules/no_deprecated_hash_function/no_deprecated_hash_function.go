package no_deprecated_hash_function

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoDeprecatedHashFunctionRule detects usage of deprecated hash functions MD4 and RIPEMD160.
type NoDeprecatedHashFunctionRule struct{}

// deprecatedHash holds metadata about a deprecated hash function.
type deprecatedHash struct {
	pkg     string // local package name (e.g. "md4")
	name    string // function name (e.g. "New")
	message string // failure message
}

var deprecatedHashes = []deprecatedHash{
	{pkg: "md4", name: "New", message: "use of deprecated hash function md4.New: MD4 is cryptographically broken"},
	{pkg: "ripemd160", name: "New", message: "use of deprecated hash function ripemd160.New: RIPEMD160 has insufficient collision resistance"},
}

// Apply applies the rule to given file.
func (*NoDeprecatedHashFunctionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintDeprecatedHash{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoDeprecatedHashFunctionRule) Name() string {
	return "noDeprecatedHashFunction"
}

// Group returns the rule group.
func (*NoDeprecatedHashFunctionRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoDeprecatedHashFunctionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintDeprecatedHash struct {
	onFailure func(lint.Failure)
}

func (w *lintDeprecatedHash) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	for _, dh := range deprecatedHashes {
		if astutils.IsPkgDotName(ce.Fun, dh.pkg, dh.name) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       ce,
				Category:   lint.FailureCategoryBadPractice,
				Failure:    dh.message,
			})
			return w
		}
	}

	return w
}
