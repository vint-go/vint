package no_http_response_misuse

import (
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoHttpResponseMisuseRule checks for mistakes using HTTP responses.
// A common mistake is to defer the closing of the response body before
// checking the error returned by http.Client.Do or similar methods.
// If the error is non-nil, the response may be nil, causing a nil pointer
// dereference in the deferred Body.Close() call.
type NoHttpResponseMisuseRule struct{}

// Apply applies the rule to given file.
func (r *NoHttpResponseMisuseRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintHttpResponseMisuse{
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
		w.checkBlockStmt(funcDecl.Body)
	}

	return failures
}

// Name returns the rule name.
func (*NoHttpResponseMisuseRule) Name() string {
	return "noHttpResponseMisuse"
}

// Group returns the rule group.
func (*NoHttpResponseMisuseRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoHttpResponseMisuseRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoHttpResponseMisuseRule) RequiresTypecheck() bool {
	return true
}

type lintHttpResponseMisuse struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

// checkBlockStmt examines consecutive statements in a block to find
// the pattern: assign resp+err, then defer resp.Body.Close() before
// checking the error.
func (w *lintHttpResponseMisuse) checkBlockStmt(block *ast.BlockStmt) {
	for i, stmt := range block.List {
		// Recurse into nested blocks (if, for, switch, etc.)
		w.recurseIntoStmt(stmt)

		// Look for assignment statements that produce (*http.Response, error)
		assign, ok := stmt.(*ast.AssignStmt)
		if !ok {
			continue
		}

		// Need at least 2 LHS values: resp, err := ...
		if len(assign.Lhs) < 2 {
			continue
		}

		// Find the response variable name
		respVarName := w.findHTTPResponseVar(assign)
		if respVarName == "" {
			continue
		}

		// Check if any subsequent statement before an error check is a
		// defer resp.Body.Close()
		w.checkStatementsAfterAssign(block.List[i+1:], respVarName)
	}
}

// findHTTPResponseVar checks if an assignment produces an *http.Response
// and returns the variable name of that response, or "" if not found.
func (w *lintHttpResponseMisuse) findHTTPResponseVar(assign *ast.AssignStmt) string {
	for _, lhs := range assign.Lhs {
		ident, ok := lhs.(*ast.Ident)
		if !ok || ident.Name == "_" {
			continue
		}
		t := w.pkg.TypeOf(ident)
		if t != nil && isHTTPResponseType(t) {
			return ident.Name
		}
	}
	return ""
}

// checkStatementsAfterAssign scans statements after the HTTP call assignment
// looking for a defer resp.Body.Close() that appears before any error check.
func (w *lintHttpResponseMisuse) checkStatementsAfterAssign(stmts []ast.Stmt, respVarName string) {
	for _, stmt := range stmts {
		// If we hit an if statement (likely error check), stop scanning
		if isErrorCheck(stmt) {
			return
		}

		// Check if this is a defer resp.Body.Close()
		deferStmt, ok := stmt.(*ast.DeferStmt)
		if !ok {
			continue
		}

		if isBodyCloseCall(deferStmt.Call, respVarName) {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 1,
				Node:       deferStmt,
				Failure:    "deferring Body.Close before checking error from HTTP call may cause nil pointer dereference",
			})
		}
	}
}

// recurseIntoStmt recurses into nested block statements to check them as well.
func (w *lintHttpResponseMisuse) recurseIntoStmt(stmt ast.Stmt) {
	switch s := stmt.(type) {
	case *ast.BlockStmt:
		w.checkBlockStmt(s)
	case *ast.IfStmt:
		w.checkBlockStmt(s.Body)
		if s.Else != nil {
			w.recurseIntoStmt(s.Else)
		}
	case *ast.ForStmt:
		w.checkBlockStmt(s.Body)
	case *ast.RangeStmt:
		w.checkBlockStmt(s.Body)
	case *ast.SwitchStmt:
		w.checkBlockStmt(s.Body)
	case *ast.TypeSwitchStmt:
		w.checkBlockStmt(s.Body)
	case *ast.SelectStmt:
		w.checkBlockStmt(s.Body)
	case *ast.CaseClause:
		fakeBlock := &ast.BlockStmt{List: s.Body}
		w.checkBlockStmt(fakeBlock)
	case *ast.CommClause:
		fakeBlock := &ast.BlockStmt{List: s.Body}
		w.checkBlockStmt(fakeBlock)
	}
}

// isHTTPResponseType checks if the given type is *net/http.Response.
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

// isBodyCloseCall checks if a call expression is resp.Body.Close()
// where resp matches the given variable name.
func isBodyCloseCall(call *ast.CallExpr, varName string) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Close" {
		return false
	}
	innerSel, ok := sel.X.(*ast.SelectorExpr)
	if !ok || innerSel.Sel.Name != "Body" {
		return false
	}
	ident, ok := innerSel.X.(*ast.Ident)
	if !ok || ident.Name != varName {
		return false
	}
	return true
}

// isErrorCheck checks if a statement is an if statement that likely checks
// an error (e.g., "if err != nil").
func isErrorCheck(stmt ast.Stmt) bool {
	ifStmt, ok := stmt.(*ast.IfStmt)
	if !ok {
		return false
	}
	binExpr, ok := ifStmt.Cond.(*ast.BinaryExpr)
	if !ok {
		return false
	}
	// Check if either side references "err" or any ident ending in err/Err
	return containsErrorIdent(binExpr.X) || containsErrorIdent(binExpr.Y)
}

// containsErrorIdent checks if an expression is an identifier that looks
// like an error variable (named "err").
func containsErrorIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == "err"
}
