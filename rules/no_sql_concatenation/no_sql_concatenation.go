package no_sql_concatenation

import (
	"go/ast"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoSqlConcatenationRule detects SQL query construction using string concatenation,
// which can lead to SQL injection vulnerabilities.
type NoSqlConcatenationRule struct{}

// Apply applies the rule to given file.
func (r *NoSqlConcatenationRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoSqlConcatenation{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules
func (r *NoSqlConcatenationRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoSqlConcatenation{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoSqlConcatenationRule) Name() string {
	return "noSqlConcatenation"
}

// Group returns the rule group.
func (*NoSqlConcatenationRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoSqlConcatenationRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// sqlMethods lists the database/sql method names whose first string argument
// is a SQL query that must not be built via concatenation.
// These include methods on *sql.DB, *sql.Tx, and *sql.Stmt.
var sqlMethods = map[string]bool{
	"Exec":            true,
	"ExecContext":     true,
	"Query":           true,
	"QueryContext":    true,
	"QueryRow":        true,
	"QueryRowContext": true,
	"Prepare":         true,
	"PrepareContext":  true,
}

type lintNoSqlConcatenation struct {
	onFailure func(lint.Failure)
}

// Visit processes function declarations to track variable assignments and
// detect SQL injection via string concatenation.
func (w *lintNoSqlConcatenation) Visit(node ast.Node) ast.Visitor {
	fd, ok := node.(*ast.FuncDecl)
	if !ok {
		return w
	}

	if fd.Body == nil {
		return nil
	}

	// Collect variables that were tainted by string concatenation.
	taintedVars := collectTaintedVars(fd.Body)

	// Walk the function body looking for SQL method calls.
	fw := &sqlCallWalker{
		onFailure:   w.onFailure,
		taintedVars: taintedVars,
	}
	ast.Walk(fw, fd.Body)

	return nil // don't recurse into the function again
}

// sqlCallWalker walks a function body to find SQL method calls with
// concatenated query strings.
type sqlCallWalker struct {
	onFailure   func(lint.Failure)
	taintedVars map[string]bool
}

func (fw *sqlCallWalker) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return fw
	}

	sel, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return fw
	}

	methodName := sel.Sel.Name
	if !sqlMethods[methodName] {
		return fw
	}

	// Determine the query argument index:
	// - Context variants have ctx as first arg, query as second
	// - Non-context variants have query as first arg
	queryArgIdx := 0
	if isContextVariant(methodName) {
		queryArgIdx = 1
	}

	if len(ce.Args) <= queryArgIdx {
		return fw
	}

	queryArg := ce.Args[queryArgIdx]
	fw.checkQueryArg(ce, queryArg, methodName)

	return fw
}

// isContextVariant returns true if the method name is a context-aware variant.
func isContextVariant(methodName string) bool {
	switch methodName {
	case "ExecContext", "QueryContext", "QueryRowContext", "PrepareContext":
		return true
	}
	return false
}

// checkQueryArg checks if the query argument contains string concatenation.
func (fw *sqlCallWalker) checkQueryArg(call *ast.CallExpr, queryArg ast.Expr, methodName string) {
	if containsConcatenation(queryArg) {
		fw.onFailure(lint.Failure{
			Confidence: 1,
			Node:       call,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "SQL query constructed using string concatenation",
		})
		return
	}

	// Check if the argument is a variable that was tainted.
	if ident, ok := queryArg.(*ast.Ident); ok {
		if fw.taintedVars[ident.Name] {
			fw.onFailure(lint.Failure{
				Confidence: 1,
				Node:       call,
				Category:   lint.FailureCategoryBadPractice,
				Failure:    "SQL query constructed using string concatenation",
			})
		}
	}
}

// containsConcatenation returns true if the expression is a binary
// expression using the + operator where at least one part is not a
// string literal (i.e., contains dynamic data).
func containsConcatenation(expr ast.Expr) bool {
	bin, ok := expr.(*ast.BinaryExpr)
	if !ok {
		return false
	}

	// Check if this is a string concatenation using +
	if bin.Op.String() != "+" {
		return false
	}

	// If both sides are string literals, it's just constant folding - safe.
	if isConstantStringExpr(bin.X) && isConstantStringExpr(bin.Y) {
		return false
	}

	return true
}

// isConstantStringExpr returns true if the expression is entirely composed of
// string literal constants (including concatenation of literals).
func isConstantStringExpr(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return true
	case *ast.BinaryExpr:
		if e.Op.String() != "+" {
			return false
		}
		return isConstantStringExpr(e.X) && isConstantStringExpr(e.Y)
	case *ast.ParenExpr:
		return isConstantStringExpr(e.X)
	default:
		return false
	}
}

// collectTaintedVars scans a block statement for variable assignments
// where the value is built via string concatenation (using = with + or +=).
func collectTaintedVars(body *ast.BlockStmt) map[string]bool {
	tainted := map[string]bool{}

	ast.Inspect(body, func(n ast.Node) bool {
		stmt, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		switch stmt.Tok.String() {
		case "+=":
			// query += something — the variable is tainted by concatenation.
			for _, lhs := range stmt.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok {
					tainted[ident.Name] = true
				}
			}
		case ":=", "=":
			// query := "..." + x or query = "..." + x
			for i, rhs := range stmt.Rhs {
				if containsConcatenation(rhs) && i < len(stmt.Lhs) {
					if ident, ok := stmt.Lhs[i].(*ast.Ident); ok {
						tainted[ident.Name] = true
					}
				}
			}
		}

		return true
	})

	return tainted
}
