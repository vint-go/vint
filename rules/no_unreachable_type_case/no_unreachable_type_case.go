package no_unreachable_type_case

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoUnreachableTypeCaseRule detects erroneous case order inside type switch
// statements where a concrete type case appears after an interface type case
// that would match it, making the concrete case unreachable.
type NoUnreachableTypeCaseRule struct{}

// Apply applies the rule to given file.
func (r *NoUnreachableTypeCaseRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintUnreachableTypeCase{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoUnreachableTypeCaseRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintUnreachableTypeCase{
		pkg:       file.Pkg,
		onFailure: onFailure,
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoUnreachableTypeCaseRule) Name() string {
	return "noUnreachableTypeCase"
}

// Group returns the rule group.
func (*NoUnreachableTypeCaseRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnreachableTypeCaseRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck indicates this rule needs type information.
func (*NoUnreachableTypeCaseRule) RequiresTypecheck() bool {
	return true
}

type lintUnreachableTypeCase struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintUnreachableTypeCase) Visit(node ast.Node) ast.Visitor {
	tsStmt, ok := node.(*ast.TypeSwitchStmt)
	if !ok {
		return w
	}

	// Collect the case clauses in order
	if tsStmt.Body == nil {
		return w
	}

	// For each case clause, resolve the types and check for unreachable cases.
	// A case is unreachable if a previous case has an interface type that the
	// current case's type implements.
	type resolvedCase struct {
		typ     types.Type
		caseExpr ast.Expr
	}

	var seenCases []resolvedCase

	for _, stmt := range tsStmt.Body.List {
		cc, ok := stmt.(*ast.CaseClause)
		if !ok {
			continue
		}

		for _, caseExpr := range cc.List {
			caseType := w.pkg.TypeOf(caseExpr)
			if caseType == nil {
				continue
			}

			// Check if any previously seen interface case would match this type
			for _, prev := range seenCases {
				prevIface := extractInterface(prev.typ)
				if prevIface == nil {
					continue
				}

				// Check if the current case type implements the previous interface
				if types.Implements(caseType, prevIface) {
					w.onFailure(lint.Failure{
						Category:   lint.FailureCategoryLogic,
						Confidence: 1,
						Node:       caseExpr,
						Failure:    fmt.Sprintf("unreachable type case: %s is already matched by earlier interface case %s", types.TypeString(caseType, nil), types.TypeString(prev.typ, nil)),
					})
					break
				}

				// Also check pointer-to-type implements interface
				ptrType := types.NewPointer(caseType)
				if types.Implements(ptrType, prevIface) {
					// Only flag if the case type itself is a pointer or the interface
					// is satisfied by value (already handled above)
					// This handles cases where caseType is not a pointer but *caseType
					// implements the interface - but the switch will match the non-pointer
					// so this isn't unreachable. Skip this case.
				}
			}

			seenCases = append(seenCases, resolvedCase{typ: caseType, caseExpr: caseExpr})
		}
	}

	return w
}

// extractInterface extracts the *types.Interface from a type, if it is one.
func extractInterface(t types.Type) *types.Interface {
	switch u := t.(type) {
	case *types.Interface:
		return u
	case *types.Named:
		if iface, ok := u.Underlying().(*types.Interface); ok {
			return iface
		}
	}
	if iface, ok := t.Underlying().(*types.Interface); ok {
		return iface
	}
	return nil
}
