package use_join_host_port

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseJoinHostPortRule checks for fmt.Sprintf calls that construct host:port
// addresses using patterns like "%s:%d" or "%s:%s". These patterns do not
// work correctly with IPv6 addresses. Use net.JoinHostPort instead.
type UseJoinHostPortRule struct{}

// Apply applies the rule to given file.
func (r *UseJoinHostPortRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintJoinHostPort{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules
func (r *UseJoinHostPortRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintJoinHostPort{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseJoinHostPortRule) Name() string {
	return "useJoinHostPort"
}

// Group returns the rule group.
func (*UseJoinHostPortRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*UseJoinHostPortRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// hostPortFormats contains format strings that construct host:port addresses
// in an IPv6-incompatible way.
var hostPortFormats = map[string]bool{
	`"%s:%d"`: true,
	`"%s:%s"`: true,
}

type lintJoinHostPort struct {
	onFailure func(lint.Failure)
}

func (w *lintJoinHostPort) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "fmt", "Sprintf") {
		return w
	}

	// fmt.Sprintf needs at least a format string argument
	if len(ce.Args) < 1 {
		return w
	}

	// Check if the first argument is a string literal matching a host:port pattern
	formatArg, ok := ce.Args[0].(*ast.BasicLit)
	if !ok || formatArg.Kind != token.STRING {
		return w
	}

	if hostPortFormats[formatArg.Value] {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "use net.JoinHostPort instead of fmt.Sprintf for host:port construction to support IPv6",
		})
	}

	return w
}
