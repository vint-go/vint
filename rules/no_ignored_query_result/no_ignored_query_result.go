package no_ignored_query_result

import (
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoIgnoredQueryResultRule detects Query()/QueryContext() calls where the
// result rows are assigned to the blank identifier (_). When the rows are
// never read or closed this can leak database connections. If no rows are
// needed, Exec()/ExecContext() should be used instead.
type NoIgnoredQueryResultRule struct{}

// Apply applies the rule to given file.
func (r *NoIgnoredQueryResultRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoIgnoredQueryResult{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoIgnoredQueryResultRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintNoIgnoredQueryResult{
		pkg:       file.Pkg,
		onFailure: onFailure,
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoIgnoredQueryResultRule) Name() string {
	return "noIgnoredQueryResult"
}

// Group returns the rule group.
func (*NoIgnoredQueryResultRule) Group() string {
	return "correctness"
}

// RequiresTypecheck returns true because this rule needs type info.
func (*NoIgnoredQueryResultRule) RequiresTypecheck() bool {
	return true
}

// CacheTier returns the cache tier for this rule.
func (*NoIgnoredQueryResultRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

type lintNoIgnoredQueryResult struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

// queryMethodNames is the set of method names that return rows which must not be ignored.
var queryMethodNames = map[string]bool{
	"Query":        true,
	"QueryContext": true,
}

func (w *lintNoIgnoredQueryResult) Visit(node ast.Node) ast.Visitor {
	assign, ok := node.(*ast.AssignStmt)
	if !ok {
		return w
	}

	// We need at least 2 LHS values (rows, err) and the first must be blank.
	if len(assign.Lhs) < 2 {
		return w
	}

	firstIdent, ok := assign.Lhs[0].(*ast.Ident)
	if !ok || firstIdent.Name != "_" {
		return w
	}

	// RHS must be a single call expression.
	if len(assign.Rhs) != 1 {
		return w
	}

	call, ok := assign.Rhs[0].(*ast.CallExpr)
	if !ok {
		return w
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	methodName := sel.Sel.Name
	if !queryMethodNames[methodName] {
		return w
	}

	// Check if the receiver type has the method from database/sql types.
	recvType := w.pkg.TypeOf(sel.X)
	if recvType == nil {
		return w
	}

	if !isDBQueryType(recvType, methodName) {
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryBadPractice,
		Confidence: 1,
		Node:       assign,
		Failure:    "query result rows are ignored; use Exec() instead of " + methodName + "() when rows are not needed",
	})

	return w
}

// isDBQueryType checks whether the given type has the named Query method
// from database/sql (i.e. the receiver is *sql.DB, *sql.Tx, or *sql.Stmt,
// or any type whose method set includes the matching method from database/sql).
func isDBQueryType(t types.Type, methodName string) bool {
	// Look up the method on the type's method set.
	mset := types.NewMethodSet(t)
	sel := mset.Lookup(nil, methodName)
	if sel == nil {
		return false
	}

	fn, ok := sel.Obj().(*types.Func)
	if !ok {
		return false
	}

	// Check that the method belongs to a type from database/sql.
	sig := fn.Signature()
	recv := sig.Recv()
	if recv == nil {
		return false
	}

	recvType := recv.Type()
	// Unwrap pointer.
	if ptr, ok := recvType.(*types.Pointer); ok {
		recvType = ptr.Elem()
	}

	named, ok := recvType.(*types.Named)
	if !ok {
		return false
	}

	obj := named.Obj()
	if obj.Pkg() == nil {
		return false
	}

	return obj.Pkg().Path() == "database/sql"
}
