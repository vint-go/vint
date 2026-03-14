package no_sql_format_string

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoSqlFormatStringRule detects SQL query construction using format string functions,
// which can lead to SQL injection vulnerabilities.
type NoSqlFormatStringRule struct{}

// Apply applies the rule to given file.
func (r *NoSqlFormatStringRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoSqlFormatString{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoSqlFormatStringRule) Name() string {
	return "noSqlFormatString"
}

// Group returns the rule group.
func (*NoSqlFormatStringRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoSqlFormatStringRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// sqlMethods lists the database/sql method names whose first string argument
// is a SQL query that must not be built via format string functions.
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

// fmtFormatFuncs lists fmt package functions that produce formatted strings.
var fmtFormatFuncs = map[string]bool{
	"Sprintf":  true,
	"Sprint":   true,
	"Sprintln": true,
	"Fprintf":  true,
	"Fprint":   true,
	"Fprintln": true,
}

type lintNoSqlFormatString struct {
	onFailure func(lint.Failure)
}

// Visit processes function declarations to track variable assignments and
// detect SQL injection via format string functions.
func (w *lintNoSqlFormatString) Visit(node ast.Node) ast.Visitor {
	fd, ok := node.(*ast.FuncDecl)
	if !ok {
		return w
	}

	if fd.Body == nil {
		return nil
	}

	// Collect variables that were tainted by fmt format functions with non-constant args.
	taintedVars := collectFmtTaintedVars(fd.Body)

	// Walk the function body looking for SQL method calls.
	fw := &sqlFmtCallWalker{
		onFailure:   w.onFailure,
		taintedVars: taintedVars,
	}
	ast.Walk(fw, fd.Body)

	return nil // don't recurse into the function again
}

// sqlFmtCallWalker walks a function body to find SQL method calls with
// format-string-constructed query strings.
type sqlFmtCallWalker struct {
	onFailure   func(lint.Failure)
	taintedVars map[string]bool
}

func (fw *sqlFmtCallWalker) Visit(node ast.Node) ast.Visitor {
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
	fw.checkQueryArg(ce, queryArg)

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

// checkQueryArg checks if the query argument was constructed via a fmt format function.
func (fw *sqlFmtCallWalker) checkQueryArg(call *ast.CallExpr, queryArg ast.Expr) {
	// Check if the argument is directly a fmt.Sprintf(...) call with non-constant args.
	if isFmtFormatCallWithNonConstantArgs(queryArg) {
		fw.onFailure(lint.Failure{
			Confidence: 1,
			Node:       call,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "SQL query constructed using format string function",
		})
		return
	}

	// Check if the argument is a variable that was tainted by fmt formatting.
	if ident, ok := queryArg.(*ast.Ident); ok {
		if fw.taintedVars[ident.Name] {
			fw.onFailure(lint.Failure{
				Confidence: 1,
				Node:       call,
				Category:   lint.FailureCategoryBadPractice,
				Failure:    "SQL query constructed using format string function",
			})
		}
	}
}

// isFmtFormatCallWithNonConstantArgs returns true if the expression is a call
// to a fmt format function (e.g., fmt.Sprintf) that has at least one
// non-constant argument beyond the format string.
func isFmtFormatCallWithNonConstantArgs(expr ast.Expr) bool {
	ce, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	if !isFmtFormatCall(ce) {
		return false
	}

	return hasNonConstantFormatArgs(ce)
}

// isFmtFormatCall returns true if the call expression calls a fmt format function.
func isFmtFormatCall(ce *ast.CallExpr) bool {
	for funcName := range fmtFormatFuncs {
		if astutils.IsPkgDotName(ce.Fun, "fmt", funcName) {
			return true
		}
	}
	return false
}

// hasNonConstantFormatArgs checks if a fmt call has non-constant arguments
// beyond the format string. For Sprintf/Fprintf/etc. the first arg is the
// format string, and subsequent args are values. For Sprint/Fprint/Sprintln
// all args are values.
func hasNonConstantFormatArgs(ce *ast.CallExpr) bool {
	if len(ce.Args) == 0 {
		return false
	}

	sel, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	funcName := sel.Sel.Name

	// Determine which args are "value" args (not the format string itself)
	var valueArgs []ast.Expr
	switch funcName {
	case "Sprintf", "Fprintf", "Fprintln":
		// Sprintf(format, args...) — skip format string
		// Fprintf(writer, format, args...) — skip writer and format string
		if funcName == "Fprintf" || funcName == "Fprintln" {
			if len(ce.Args) > 2 {
				valueArgs = ce.Args[2:]
			} else if funcName == "Fprintln" && len(ce.Args) > 1 {
				valueArgs = ce.Args[1:]
			}
			// For Fprintf with just format + writer, no non-constant args
		} else {
			if len(ce.Args) > 1 {
				valueArgs = ce.Args[1:]
			}
		}
	case "Sprint", "Sprintln":
		// Sprint(args...) / Sprintln(args...) — all args are values
		valueArgs = ce.Args
	case "Fprint":
		// Fprint(writer, args...) — skip writer
		if len(ce.Args) > 1 {
			valueArgs = ce.Args[1:]
		}
	default:
		return false
	}

	// Check if any value arg is non-constant
	for _, arg := range valueArgs {
		if !isConstantExpr(arg) {
			return true
		}
	}

	return false
}

// isConstantExpr returns true if the expression is a compile-time constant
// (basic literal or concatenation of literals).
func isConstantExpr(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return true
	case *ast.BinaryExpr:
		return isConstantExpr(e.X) && isConstantExpr(e.Y)
	case *ast.ParenExpr:
		return isConstantExpr(e.X)
	default:
		return false
	}
}

// collectFmtTaintedVars scans a block statement for variable assignments
// where the value is built via a fmt format function with non-constant args.
func collectFmtTaintedVars(body *ast.BlockStmt) map[string]bool {
	tainted := map[string]bool{}

	ast.Inspect(body, func(n ast.Node) bool {
		stmt, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		// Check := and = assignments
		if stmt.Tok.String() == ":=" || stmt.Tok.String() == "=" {
			for i, rhs := range stmt.Rhs {
				if isFmtFormatCallWithNonConstantArgs(rhs) && i < len(stmt.Lhs) {
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
