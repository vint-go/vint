package use_consistent_receiver_name

import (
	"fmt"
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/internal/typeparams"
	"github.com/vint-go/vint/lint"
)

// UseConsistentReceiverNameRule checks that all methods on the same type use the
// same receiver name consistently.
type UseConsistentReceiverNameRule struct{}

// Apply applies the rule to the given file.
func (r *UseConsistentReceiverNameRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	typeReceiver := map[string]string{}
	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
			continue
		}

		names := fn.Recv.List[0].Names
		if len(names) < 1 {
			continue
		}
		name := names[0].Name

		recv := typeparams.ReceiverType(fn)
		if prev, ok := typeReceiver[recv]; ok && prev != name {
			failures = append(failures, lint.Failure{
				Node:       fn.Recv.List[0],
				Confidence: 1,
				Category:   lint.FailureCategoryStyle,
				Failure:    fmt.Sprintf("receiver name %s should be consistent with previous receiver name %s for %s", name, prev, recv),
			})
			continue
		}

		typeReceiver[recv] = name
	}

	return failures
}

// Name returns the rule name.
func (*UseConsistentReceiverNameRule) Name() string {
	return "useConsistentReceiverName"
}

// Group returns the rule group.
func (*UseConsistentReceiverNameRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseConsistentReceiverNameRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
