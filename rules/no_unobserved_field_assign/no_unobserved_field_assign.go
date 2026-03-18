package no_unobserved_field_assign

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoUnobservedFieldAssignRule detects field assignments on value receivers
// that will never be observed because the changes are made to a copy.
type NoUnobservedFieldAssignRule struct{}

// Apply applies the rule to the given file.
func (r *NoUnobservedFieldAssignRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()
	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Recv == nil {
			continue // skip non-methods
		}

		receiver := funcDecl.Recv.List[0]
		if r.mustSkip(receiver, file.Pkg) {
			continue
		}

		if len(receiver.Names) == 0 {
			continue
		}

		receiverName := receiver.Names[0].Name
		if receiverName == "_" {
			continue
		}

		fieldAssignments := r.getFieldAssignments(receiverName, funcDecl.Body)
		if len(fieldAssignments) == 0 {
			continue
		}

		// If the method returns the receiver (or a field of it), the
		// modification is intentional and observable by the caller.
		if r.methodReturnsReceiver(receiverName, funcDecl.Body) {
			continue
		}

		for _, assignment := range fieldAssignments {
			failures = append(failures, lint.Failure{
				Node:       assignment,
				Confidence: 1,
				Category:   lint.FailureCategoryLogic,
				Failure:    "field assignment to value receiver will not be observed; did you mean to use a pointer receiver?",
			})
		}
	}

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoUnobservedFieldAssignRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok || funcDecl.Recv == nil {
		return nil
	}

	file.Pkg.TypeCheck()

	receiver := funcDecl.Recv.List[0]
	if r.mustSkip(receiver, file.Pkg) {
		return nil
	}

	if len(receiver.Names) == 0 {
		return nil
	}

	receiverName := receiver.Names[0].Name
	if receiverName == "_" {
		return nil
	}

	fieldAssignments := r.getFieldAssignments(receiverName, funcDecl.Body)
	if len(fieldAssignments) == 0 {
		return nil
	}

	if r.methodReturnsReceiver(receiverName, funcDecl.Body) {
		return nil
	}

	var failures []lint.Failure
	for _, assignment := range fieldAssignments {
		failures = append(failures, lint.Failure{
			Node:       assignment,
			Confidence: 1,
			Category:   lint.FailureCategoryLogic,
			Failure:    "field assignment to value receiver will not be observed; did you mean to use a pointer receiver?",
		})
	}
	return failures
}

// Name returns the rule name.
func (*NoUnobservedFieldAssignRule) Name() string {
	return "noUnobservedFieldAssign"
}

// Group returns the rule group.
func (*NoUnobservedFieldAssignRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnobservedFieldAssignRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule uses type information.
func (*NoUnobservedFieldAssignRule) RequiresTypecheck() bool {
	return true
}

// mustSkip returns true if the method should be skipped (pointer receiver,
// anonymous receiver, or receiver is a map/slice type).
func (r *NoUnobservedFieldAssignRule) mustSkip(receiver *ast.Field, pkg *lint.Package) bool {
	if _, ok := receiver.Type.(*ast.StarExpr); ok {
		return true // pointer receiver: mutations are observed
	}

	if len(receiver.Names) < 1 {
		return true // anonymous receiver
	}

	if r.skipType(receiver.Type, pkg) {
		return true // map or slice: mutations are observed via reference
	}

	return false
}

// skipType returns true if the receiver type is a map or slice (reference types
// where mutations through a value receiver are still observed).
func (*NoUnobservedFieldAssignRule) skipType(t ast.Expr, pkg *lint.Package) bool {
	rt := pkg.TypeOf(t)
	if rt == nil {
		return false
	}

	rt = rt.Underlying()
	rtName := rt.String()

	return strings.HasPrefix(rtName, "[]") || strings.HasPrefix(rtName, "map[")
}

// getFieldAssignments returns all AST nodes where a field of the value receiver
// is assigned (e.g. c.count = ... or c.count++).
func (r *NoUnobservedFieldAssignRule) getFieldAssignments(receiverName string, body *ast.BlockStmt) []ast.Node {
	if body == nil {
		return nil
	}

	finder := func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.IncDecStmt:
			// c.field++ or c.field--
			return r.isSelectorOnReceiver(node.X, receiverName)
		case *ast.AssignStmt:
			for _, lhs := range node.Lhs {
				if r.isSelectorOnReceiver(lhs, receiverName) {
					return true
				}
			}
		}
		return false
	}

	return astutils.PickNodes(body, finder)
}

// isSelectorOnReceiver checks whether the given expression is a selector
// expression on the named receiver (e.g. c.field).
func (*NoUnobservedFieldAssignRule) isSelectorOnReceiver(expr ast.Expr, receiverName string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == receiverName
}

// methodReturnsReceiver checks whether the method returns the receiver
// (or a field/address of it), which would make the mutation observable.
func (r *NoUnobservedFieldAssignRule) methodReturnsReceiver(receiverName string, body *ast.BlockStmt) bool {
	if body == nil {
		return false
	}

	finder := func(n ast.Node) bool {
		ret, ok := n.(*ast.ReturnStmt)
		if !ok {
			return false
		}
		for _, result := range ret.Results {
			switch e := result.(type) {
			case *ast.Ident:
				if e.Name == receiverName {
					return true
				}
			case *ast.UnaryExpr:
				if e.Op == token.AND {
					if ident, ok := e.X.(*ast.Ident); ok && ident.Name == receiverName {
						return true
					}
				}
			case *ast.SelectorExpr:
				if ident, ok := e.X.(*ast.Ident); ok && ident.Name == receiverName {
					return true
				}
			}
		}
		return false
	}

	return astutils.SeekNode[ast.Node](body, finder) != nil
}
