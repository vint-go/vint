package no_excessive_statements

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"

	"github.com/vint-go/vint/lint"
)

// NoExcessiveStatementsRule checks that functions do not exceed a maximum number of statements.
type NoExcessiveStatementsRule struct {
	maxStatements int
}

const defaultMaxStatements = 40

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *NoExcessiveStatementsRule) Configure(arguments lint.Arguments) error {
	if len(arguments) < 1 {
		r.maxStatements = defaultMaxStatements
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		// try direct int64 argument
		stmts, ok := arguments[0].(int64)
		if !ok {
			return fmt.Errorf(`invalid argument to the "noExcessiveStatements" rule, expecting a k,v map or int64, got %T`, arguments[0])
		}
		r.maxStatements = int(stmts)
		return nil
	}

	r.maxStatements = defaultMaxStatements
	for k, v := range argKV {
		if isRuleOption(k, "statements") {
			stmts, ok := v.(int64)
			if !ok {
				return fmt.Errorf(`invalid configuration value for statements in "noExcessiveStatements" rule; need int64 but got %T`, v)
			}
			r.maxStatements = int(stmts)
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoExcessiveStatementsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if r.maxStatements < 0 {
		return nil // disabled
	}

	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		body := funcDecl.Body
		if body == nil || len(body.List) == 0 {
			continue
		}

		stmtCount := r.countStmts(body.List)
		if stmtCount > r.maxStatements {
			failures = append(failures, lint.Failure{
				Confidence: 1,
				Category:   lint.FailureCategoryComplexity,
				Failure:    fmt.Sprintf("function %s has too many statements (%d > %d)", funcName(funcDecl), stmtCount, r.maxStatements),
				Node:       funcDecl,
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoExcessiveStatementsRule) Name() string {
	return "noExcessiveStatements"
}

// Group returns the rule group.
func (*NoExcessiveStatementsRule) Group() string {
	return "complexity"
}

func (r *NoExcessiveStatementsRule) countStmts(stmts []ast.Stmt) int {
	count := 0
	for _, s := range stmts {
		switch stmt := s.(type) {
		case *ast.BlockStmt:
			count += r.countStmts(stmt.List)
		case *ast.IfStmt:
			count += 1 + r.countBodyStmts(stmt)
			if stmt.Else != nil {
				elseBody, ok := stmt.Else.(*ast.BlockStmt)
				if ok {
					count += r.countStmts(elseBody.List)
				}
			}
		case *ast.ForStmt, *ast.RangeStmt,
			*ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
			count += 1 + r.countBodyStmts(stmt)
		case *ast.CaseClause:
			count += r.countStmts(stmt.Body)
		case *ast.AssignStmt:
			count += 1 + r.countFuncLitStmts(stmt.Rhs[0])
		case *ast.GoStmt:
			count += 1 + r.countFuncLitStmts(stmt.Call.Fun)
		case *ast.DeferStmt:
			count += 1 + r.countFuncLitStmts(stmt.Call.Fun)
		default:
			count++
		}
	}
	return count
}

func (r *NoExcessiveStatementsRule) countFuncLitStmts(expr ast.Expr) int {
	if block, ok := expr.(*ast.FuncLit); ok {
		return r.countStmts(block.Body.List)
	}
	return 0
}

func (r *NoExcessiveStatementsRule) countBodyStmts(t any) int {
	i := reflect.ValueOf(t).Elem().FieldByName("Body").Elem().FieldByName("List").Interface()
	return r.countStmts(i.([]ast.Stmt))
}

// funcName returns the name representation of a function or method:
// "(Type).Name" for methods or simply "Name" for functions.
func funcName(fn *ast.FuncDecl) string {
	declarationHasReceiver := fn.Recv != nil && fn.Recv.NumFields() > 0
	if declarationHasReceiver {
		typ := fn.Recv.List[0].Type
		return fmt.Sprintf("(%s).%s", recvString(typ), fn.Name)
	}

	return fn.Name.Name
}

// recvString returns a string representation of recv of the
// form "T", "*T", or "BADRECV" (if not a proper receiver type).
func recvString(recv ast.Expr) string {
	switch t := recv.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + recvString(t.X)
	}
	return "BADRECV"
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name lowercased and without hyphens.
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}
