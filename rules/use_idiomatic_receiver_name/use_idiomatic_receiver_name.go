package use_idiomatic_receiver_name

import (
	"fmt"
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseIdiomaticReceiverNameRule checks that method receiver names are idiomatic:
// not generic names like "this" or "self", and not underscores.
type UseIdiomaticReceiverNameRule struct{}

// Apply applies the rule to the given file.
func (r *UseIdiomaticReceiverNameRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
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

		if name == "_" {
			failures = append(failures, lint.Failure{
				Node:       fn.Recv.List[0],
				Confidence: 1,
				Category:   lint.FailureCategoryStyle,
				Failure:    "receiver name should not be an underscore, omit the name if it is unused",
			})
			continue
		}

		if name == "this" || name == "self" {
			failures = append(failures, lint.Failure{
				Node:       fn.Recv.List[0],
				Confidence: 1,
				Category:   lint.FailureCategoryStyle,
				Failure:    fmt.Sprintf("receiver name should be a reflection of its identity; don't use generic names such as %q", name),
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*UseIdiomaticReceiverNameRule) Name() string {
	return "useIdiomaticReceiverName"
}

// Group returns the rule group.
func (*UseIdiomaticReceiverNameRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseIdiomaticReceiverNameRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
