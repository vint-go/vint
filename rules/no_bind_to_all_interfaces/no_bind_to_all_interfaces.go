package no_bind_to_all_interfaces

import (
	"go/ast"
	"go/token"
	"regexp"
	"strconv"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoBindToAllInterfacesRule detects when a network listener binds to all interfaces.
type NoBindToAllInterfacesRule struct{}

// Apply applies the rule to given file.
func (r *NoBindToAllInterfacesRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintBindToAllInterfaces{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoBindToAllInterfacesRule) Name() string {
	return "noBindToAllInterfaces"
}

// Group returns the rule group.
func (*NoBindToAllInterfacesRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoBindToAllInterfacesRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

var bindAllPattern = regexp.MustCompile(`^(0\.0\.0\.0|:).*$`)

type lintBindToAllInterfaces struct {
	onFailure func(lint.Failure)
}

func (w *lintBindToAllInterfaces) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check for net.Listen or tls.Listen
	if !astutils.IsPkgDotName(ce.Fun, "net", "Listen") &&
		!astutils.IsPkgDotName(ce.Fun, "tls", "Listen") {
		return w
	}

	// net.Listen(network, address) - address is the second argument
	// tls.Listen(network, address, config) - address is the second argument
	if len(ce.Args) < 2 {
		return w
	}

	addrArg := ce.Args[1]
	lit, ok := addrArg.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return w
	}

	addr, err := strconv.Unquote(lit.Value)
	if err != nil {
		return w
	}

	if bindAllPattern.MatchString(addr) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "binding to all interfaces is a security risk, bind to a specific interface instead",
		})
	}

	return w
}
