package no_sql_db_ping

import (
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoSqlDbPingRule disallows calling (*database/sql.DB).Ping without a context.
// The method (*sql.DB).Ping does not accept a context.Context parameter,
// which prevents the caller from controlling cancellation, deadlines, and tracing.
// Use (*sql.DB).PingContext instead, which accepts a context.Context as its first argument.
type NoSqlDbPingRule struct{}

// Apply applies the rule to given file.
func (r *NoSqlDbPingRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoSqlDbPing{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoSqlDbPingRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	file.Pkg.TypeCheck()

	w := &lintNoSqlDbPing{
		pkg:       file.Pkg,
		onFailure: onFailure,
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoSqlDbPingRule) Name() string {
	return "noSqlDbPing"
}

// Group returns the rule group.
func (*NoSqlDbPingRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoSqlDbPingRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoSqlDbPingRule) RequiresTypecheck() bool {
	return true
}

type lintNoSqlDbPing struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintNoSqlDbPing) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check if this is a selector expression call: <expr>.Ping(...)
	sel, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	if sel.Sel.Name != "Ping" {
		return w
	}

	// Check if the receiver type is *database/sql.DB
	t := w.pkg.TypeOf(sel.X)
	if t == nil {
		return w
	}

	if !isSqlDBType(t) {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "(*sql.DB).Ping does not accept a context; use (*sql.DB).PingContext instead",
	})

	return w
}

// isSqlDBType checks if the given type is *database/sql.DB.
func isSqlDBType(t types.Type) bool {
	ptr, ok := t.(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := ptr.Elem().(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	return obj.Name() == "DB" && obj.Pkg() != nil && obj.Pkg().Path() == "database/sql"
}
