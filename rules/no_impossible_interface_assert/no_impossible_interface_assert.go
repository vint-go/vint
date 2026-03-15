package no_impossible_interface_assert

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoImpossibleInterfaceAssertRule flags impossible interface-to-interface type assertions.
// A type assertion from one interface to another is impossible when both interfaces
// contain methods with the same name but different signatures, since no type could
// implement both interfaces simultaneously.
type NoImpossibleInterfaceAssertRule struct{}

// Apply applies the rule to given file.
func (r *NoImpossibleInterfaceAssertRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintImpossibleIfaceAssert{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoImpossibleInterfaceAssertRule) Name() string {
	return "noImpossibleInterfaceAssert"
}

// Group returns the rule group.
func (*NoImpossibleInterfaceAssertRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoImpossibleInterfaceAssertRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoImpossibleInterfaceAssertRule) RequiresTypecheck() bool {
	return true
}

type lintImpossibleIfaceAssert struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintImpossibleIfaceAssert) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.TypeAssertExpr:
		// e.g. x.(T)
		if n.Type == nil {
			// type switch with no explicit type (x.(type)) is handled below
			return w
		}
		w.checkAssertion(n, n.X, n.Type)

	case *ast.TypeSwitchStmt:
		// e.g. switch x.(type) { case T1: ... case T2: ... }
		var expr ast.Expr
		switch a := n.Assign.(type) {
		case *ast.ExprStmt:
			assert, ok := a.X.(*ast.TypeAssertExpr)
			if !ok {
				return w
			}
			expr = assert.X
		case *ast.AssignStmt:
			assert, ok := a.Rhs[0].(*ast.TypeAssertExpr)
			if !ok {
				return w
			}
			expr = assert.X
		default:
			return w
		}

		for _, clause := range n.Body.List {
			cc, ok := clause.(*ast.CaseClause)
			if !ok {
				continue
			}
			for _, caseType := range cc.List {
				w.checkAssertion(caseType, expr, caseType)
			}
		}
	}

	return w
}

// checkAssertion checks if a type assertion from the type of src to the type of dst
// is impossible because both are interfaces with conflicting method signatures.
func (w *lintImpossibleIfaceAssert) checkAssertion(node ast.Node, src ast.Expr, dst ast.Expr) {
	srcType := w.pkg.TypeOf(src)
	if srcType == nil {
		return
	}
	dstType := w.pkg.TypeOf(dst)
	if dstType == nil {
		return
	}

	// For type expressions (like in type assertions), TypeOf returns
	// the *types.Type value. But for switch case types we need the actual type.
	// Try to get the underlying interface for both.
	srcIface := interfaceOf(srcType)
	dstIface := interfaceOf(dstType)

	if srcIface == nil || dstIface == nil {
		return
	}

	if srcIface.NumMethods() == 0 || dstIface.NumMethods() == 0 {
		// If either interface is empty (or any), the assertion is always possible.
		return
	}

	// Check for conflicting methods: same name, different signature
	srcMethodName, dstMethodName := findConflictingMethod(srcIface, dstIface)
	if srcMethodName == "" {
		return
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryLogic,
		Confidence: 1,
		Node:       node,
		Failure:    fmt.Sprintf("impossible type assertion: no type can implement both interfaces (method %s has conflicting signatures)", srcMethodName),
	})
	_ = dstMethodName
}

// interfaceOf extracts the *types.Interface from a type, if it is one.
func interfaceOf(t types.Type) *types.Interface {
	// Unwrap type if needed
	switch u := t.(type) {
	case *types.Interface:
		return u
	case *types.Named:
		if iface, ok := u.Underlying().(*types.Interface); ok {
			return iface
		}
	}
	// Also check underlying directly
	if iface, ok := t.Underlying().(*types.Interface); ok {
		return iface
	}
	return nil
}

// findConflictingMethod checks if two interfaces have a method with the same name
// but different signatures. Returns the method names if found, empty strings otherwise.
func findConflictingMethod(a, b *types.Interface) (string, string) {
	for i := 0; i < a.NumMethods(); i++ {
		methodA := a.Method(i)
		for j := 0; j < b.NumMethods(); j++ {
			methodB := b.Method(j)
			if methodA.Name() == methodB.Name() {
				// Same name — check if signatures differ
				if !types.Identical(methodA.Type(), methodB.Type()) {
					return methodA.Name(), methodB.Name()
				}
			}
		}
	}
	return "", ""
}
