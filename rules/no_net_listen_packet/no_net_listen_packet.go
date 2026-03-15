package no_net_listen_packet

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoNetListenPacketRule disallows calling net.ListenPacket without a context.
type NoNetListenPacketRule struct{}

// Apply applies the rule to given file.
func (r *NoNetListenPacketRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoNetListenPacket{onFailure: func(f lint.Failure) {
		failures = append(failures, f)
	}}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoNetListenPacketRule) Name() string {
	return "noNetListenPacket"
}

// Group returns the rule group.
func (*NoNetListenPacketRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoNetListenPacketRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoNetListenPacket struct {
	onFailure func(lint.Failure)
}

func (w *lintNoNetListenPacket) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(ce.Fun, "net", "ListenPacket") {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "net.ListenPacket does not accept a context; use (*net.ListenConfig).ListenPacket instead",
	})

	return w
}
