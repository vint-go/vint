package no_unclosed_bodies

import (
	"go/ast"
	"go/types"
	"strings"

	"github.com/strowk/vint/lint"
)

// NoUnclosedBodiesRule warns when HTTP response bodies are not closed.
type NoUnclosedBodiesRule struct{}

// Apply applies the rule to given file.
func (r *NoUnclosedBodiesRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoUnclosedBodies{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}
		w.checkFunctionBody(funcDecl.Body)
	}

	return failures
}

// Name returns the rule name.
func (*NoUnclosedBodiesRule) Name() string {
	return "noUnclosedBodies"
}

type lintNoUnclosedBodies struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintNoUnclosedBodies) checkFunctionBody(body *ast.BlockStmt) {
	type respVar struct {
		name string
		node ast.Node
	}
	var responseVars []respVar

	ast.Inspect(body, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		// Check for blank identifier assignments: when RHS is a single call
		// returning *http.Response and the corresponding LHS is _, flag
		// immediately since the body can never be closed.
		if len(assign.Rhs) == 1 {
			if call, ok := assign.Rhs[0].(*ast.CallExpr); ok {
				if w.callReturnsHTTPResponseAt(call, assign.Lhs) {
					return true
				}
			}
		}

		for _, lhs := range assign.Lhs {
			ident, ok := lhs.(*ast.Ident)
			if !ok || ident.Name == "_" {
				continue
			}
			t := w.pkg.TypeOf(ident)
			if t != nil && isHTTPResponseType(t) {
				responseVars = append(responseVars, respVar{
					name: ident.Name,
					node: assign,
				})
			}
		}
		return true
	})

	for _, rv := range responseVars {
		if !hasBodyClose(body, rv.name) {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryBadPractice,
				Confidence: 1,
				Node:       rv.node,
				Failure:    "response body must be closed",
			})
		}
	}

	// Detect calls returning *http.Response where the return value is
	// entirely discarded (expression statements, go statements, defer statements).
	ast.Inspect(body, func(n ast.Node) bool {
		var call *ast.CallExpr
		switch stmt := n.(type) {
		case *ast.ExprStmt:
			c, ok := stmt.X.(*ast.CallExpr)
			if !ok {
				return true
			}
			call = c
		case *ast.GoStmt:
			call = stmt.Call
		case *ast.DeferStmt:
			call = stmt.Call
		default:
			return true
		}

		if w.callReturnsHTTPResponse(call) {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryBadPractice,
				Confidence: 1,
				Node:       call,
				Failure:    "response body must be closed",
			})
		}
		return true
	})
}

// callReturnsHTTPResponseAt checks whether a function call returns *http.Response
// at a position where the corresponding LHS variable is a blank identifier.
// If so, it emits a failure and returns true.
func (w *lintNoUnclosedBodies) callReturnsHTTPResponseAt(call *ast.CallExpr, lhs []ast.Expr) bool {
	t := w.pkg.TypeOf(call.Fun)
	if t == nil {
		return false
	}
	sig, ok := t.(*types.Signature)
	if !ok {
		return false
	}
	results := sig.Results()
	for i := 0; i < results.Len() && i < len(lhs); i++ {
		if !isHTTPResponseType(results.At(i).Type()) {
			continue
		}
		ident, ok := lhs[i].(*ast.Ident)
		if ok && ident.Name == "_" {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryBadPractice,
				Confidence: 1,
				Node:       call,
				Failure:    "response body must be closed",
			})
			return true
		}
	}
	return false
}

// callReturnsHTTPResponse checks whether a function call returns *http.Response
// as one of its return values.
func (w *lintNoUnclosedBodies) callReturnsHTTPResponse(call *ast.CallExpr) bool {
	t := w.pkg.TypeOf(call.Fun)
	if t == nil {
		return false
	}
	sig, ok := t.(*types.Signature)
	if !ok {
		return false
	}
	results := sig.Results()
	for i := 0; i < results.Len(); i++ {
		if isHTTPResponseType(results.At(i).Type()) {
			return true
		}
	}
	return false
}

func isHTTPResponseType(t types.Type) bool {
	ptr, ok := t.(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := ptr.Elem().(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	return obj.Name() == "Response" && obj.Pkg() != nil && obj.Pkg().Path() == "net/http"
}

func hasBodyClose(body ast.Node, varName string) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if found {
			return false
		}
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "Close" {
			return true
		}
		innerSel, ok := sel.X.(*ast.SelectorExpr)
		if !ok || innerSel.Sel.Name != "Body" {
			return true
		}
		ident, ok := innerSel.X.(*ast.Ident)
		if !ok || ident.Name != varName {
			return true
		}
		found = true
		return false
	})
	return found
}

func (*NoUnclosedBodiesRule) Group() string {
	return "correctness"
}

func (*NoUnclosedBodiesRule) RequiresTypecheck() bool {
	return true
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name lowercased and without hyphens.
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}
