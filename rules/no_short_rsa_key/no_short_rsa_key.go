package no_short_rsa_key

import (
	"fmt"
	"go/ast"
	"strconv"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoShortRsaKeyRule detects calls to rsa.GenerateKey with key sizes less than 2048 bits.
type NoShortRsaKeyRule struct{}

const minRSAKeyBits = 2048

// Apply applies the rule to given file.
func (*NoShortRsaKeyRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintShortRsaKey{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoShortRsaKeyRule) Name() string {
	return "noShortRsaKey"
}

// Group returns the rule group.
func (*NoShortRsaKeyRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoShortRsaKeyRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintShortRsaKey struct {
	onFailure func(lint.Failure)
}

func (w *lintShortRsaKey) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "rsa", "GenerateKey") {
		return w
	}

	// rsa.GenerateKey(random, bits) — check the second argument
	if len(ce.Args) < 2 {
		return w
	}

	bitsArg := ce.Args[1]
	lit, ok := bitsArg.(*ast.BasicLit)
	if !ok {
		return w
	}

	bits, err := strconv.Atoi(lit.Value)
	if err != nil {
		return w
	}

	if bits < minRSAKeyBits {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    fmt.Sprintf("RSA key length %d is too short, minimum 2048 bits recommended", bits),
		})
	}

	return w
}
