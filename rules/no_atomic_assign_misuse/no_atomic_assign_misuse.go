package no_atomic_assign_misuse

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoAtomicAssignMisuseRule detects misuse of sync/atomic operations where the
// result is assigned back to the variable passed as argument, introducing a race.
type NoAtomicAssignMisuseRule struct{}

// Apply applies the rule to given file.
func (r *NoAtomicAssignMisuseRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoAtomicAssignMisuse{
		pkgTypesInfo: file.Pkg.TypesInfo(),
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoAtomicAssignMisuseRule) Name() string {
	return "noAtomicAssignMisuse"
}

// Group returns the rule group.
func (*NoAtomicAssignMisuseRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoAtomicAssignMisuseRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

type lintNoAtomicAssignMisuse struct {
	pkgTypesInfo *types.Info
	onFailure    func(lint.Failure)
}

// atomicAddFunctions is the set of sync/atomic Add functions that return the new value.
var atomicAddFunctions = map[string]bool{
	"AddInt32":   true,
	"AddInt64":   true,
	"AddUint32":  true,
	"AddUint64":  true,
	"AddUintptr": true,
}

func (w *lintNoAtomicAssignMisuse) Visit(node ast.Node) ast.Visitor {
	n, ok := node.(*ast.AssignStmt)
	if !ok {
		return w
	}

	if len(n.Lhs) != len(n.Rhs) {
		return w
	}

	// Skip short variable declarations (:=) with a single pair — that is always fine.
	if len(n.Lhs) == 1 && n.Tok == token.DEFINE {
		return w
	}

	for i, right := range n.Rhs {
		call, ok := right.(*ast.CallExpr)
		if !ok {
			continue
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		// Verify it is from sync/atomic using type info when available.
		pkgIdent, _ := sel.X.(*ast.Ident)
		if w.pkgTypesInfo != nil {
			pkgName, ok := w.pkgTypesInfo.Uses[pkgIdent].(*types.PkgName)
			if !ok || pkgName.Imported().Path() != "sync/atomic" {
				continue
			}
		}

		if !atomicAddFunctions[sel.Sel.Name] {
			continue
		}

		left := n.Lhs[i]
		if len(call.Args) != 2 {
			continue
		}
		arg := call.Args[0]
		broken := false

		if uarg, ok := arg.(*ast.UnaryExpr); ok && uarg.Op == token.AND {
			broken = astutils.GoFmt(left) == astutils.GoFmt(uarg.X)
		} else if star, ok := left.(*ast.StarExpr); ok {
			broken = astutils.GoFmt(star.X) == astutils.GoFmt(arg)
		}

		if broken {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 1,
				Node:       n,
				Failure:    "direct assignment to atomic value",
			})
		}
	}

	return w
}
