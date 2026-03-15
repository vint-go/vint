package no_constant_result

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoConstantResultRule reports named result parameters that always return
// the same constant value across all execution paths. By default, only
// unexported (private) functions are checked. Set check-exported to true
// to also analyze exported functions.
type NoConstantResultRule struct {
	checkExported bool
}

// Configure validates and applies the rule configuration.
func (r *NoConstantResultRule) Configure(arguments lint.Arguments) error {
	if len(arguments) < 1 {
		r.checkExported = false
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noConstantResult" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		if isRuleOption(k, "check-exported") {
			val, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for check-exported in "noConstantResult" rule; need bool but got %T`, v)
			}
			r.checkExported = val
		}
	}

	return nil
}

// resultInfo tracks a single named result parameter of a function declaration.
type resultInfo struct {
	field     *ast.Field
	nameIdent *ast.Ident
	name      string
	index     int // positional index among all individual result names
}

// Apply applies the rule to the given file.
func (r *NoConstantResultRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}

		// Skip exported functions unless check-exported is enabled
		if !r.checkExported && funcDecl.Name != nil && ast.IsExported(funcDecl.Name.Name) {
			continue
		}

		results := collectNamedResults(funcDecl)
		if len(results) == 0 {
			continue
		}

		// Collect all return statements in the function body.
		returnStmts := collectReturnStmts(funcDecl.Body)

		// If there are no explicit return statements, skip. This can
		// happen for functions that use bare returns or panic.
		if len(returnStmts) == 0 {
			continue
		}

		// For each named result, check if all return statements provide
		// the same constant value at its position.
		totalResults := countTotalResults(funcDecl)
		for _, ri := range results {
			constVal, isConst := checkConstantResult(returnStmts, ri.index, totalResults)
			if !isConst {
				continue
			}

			failures = append(failures, lint.Failure{
				Confidence: 1,
				Node:       funcDecl, // report at the function declaration
				Category:   lint.FailureCategoryBadPractice,
				Failure:    fmt.Sprintf("result %s is always %s", ri.name, constVal),
			})
		}
	}

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoConstantResultRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok || funcDecl.Body == nil {
		return nil
	}

	// Skip exported functions unless check-exported is enabled
	if !r.checkExported && funcDecl.Name != nil && ast.IsExported(funcDecl.Name.Name) {
		return nil
	}

	results := collectNamedResults(funcDecl)
	if len(results) == 0 {
		return nil
	}

	returnStmts := collectReturnStmts(funcDecl.Body)
	if len(returnStmts) == 0 {
		return nil
	}

	var failures []lint.Failure
	totalResults := countTotalResults(funcDecl)
	for _, ri := range results {
		constVal, isConst := checkConstantResult(returnStmts, ri.index, totalResults)
		if !isConst {
			continue
		}

		failures = append(failures, lint.Failure{
			Confidence: 1,
			Node:       funcDecl,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    fmt.Sprintf("result %s is always %s", ri.name, constVal),
		})
	}

	return failures
}

// Name returns the rule name.
func (*NoConstantResultRule) Name() string {
	return "noConstantResult"
}

// Group returns the rule group.
func (*NoConstantResultRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoConstantResultRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// collectNamedResults extracts a flat list of named result parameters
// from a function declaration.
func collectNamedResults(fn *ast.FuncDecl) []resultInfo {
	if fn.Type.Results == nil {
		return nil
	}
	var result []resultInfo
	idx := 0
	for _, field := range fn.Type.Results.List {
		if len(field.Names) == 0 {
			idx++
			continue
		}
		for _, name := range field.Names {
			if name.Name == "_" {
				idx++
				continue
			}
			result = append(result, resultInfo{
				field:     field,
				nameIdent: name,
				name:      name.Name,
				index:     idx,
			})
			idx++
		}
	}
	return result
}

// countTotalResults counts the total number of result values (named and
// unnamed) in a function's result list.
func countTotalResults(fn *ast.FuncDecl) int {
	if fn.Type.Results == nil {
		return 0
	}
	count := 0
	for _, field := range fn.Type.Results.List {
		if len(field.Names) == 0 {
			count++
		} else {
			count += len(field.Names)
		}
	}
	return count
}

// collectReturnStmts collects all explicit return statements in a
// function body (not inside nested function literals). Only returns
// with explicit result expressions are collected; bare returns are
// ignored (they imply named results keep their current value, which
// complicates analysis).
func collectReturnStmts(body *ast.BlockStmt) []*ast.ReturnStmt {
	var stmts []*ast.ReturnStmt
	ast.Inspect(body, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.FuncLit:
			// Don't descend into nested function literals.
			return false
		}
		ret, ok := n.(*ast.ReturnStmt)
		if ok && len(ret.Results) > 0 {
			stmts = append(stmts, ret)
		}
		return true
	})
	return stmts
}

// checkConstantResult checks whether a specific result position always
// holds the same constant value across all return statements.
// It returns the rendered constant value and true, or ("", false) if it varies.
func checkConstantResult(returns []*ast.ReturnStmt, resultIndex, totalResults int) (string, bool) {
	if len(returns) == 0 {
		return "", false
	}

	var seenValue string
	first := true

	for _, ret := range returns {
		// If return doesn't have the expected number of results, skip
		// analysis (could be a multi-valued call). This is conservative.
		if len(ret.Results) != totalResults {
			return "", false
		}
		if resultIndex >= len(ret.Results) {
			return "", false
		}

		expr := ret.Results[resultIndex]
		if !isConstantExpr(expr) {
			return "", false
		}

		rendered := astutils.GoFmt(expr)
		if first {
			seenValue = rendered
			first = false
		} else if rendered != seenValue {
			return "", false
		}
	}

	if first {
		// No returns processed
		return "", false
	}

	return seenValue, true
}

// isConstantExpr returns true if the expression is a compile-time constant
// literal (basic lit, unary of basic lit, identifiers like true/false/nil).
func isConstantExpr(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return true
	case *ast.Ident:
		// true, false, nil are constant identifiers
		return e.Name == "true" || e.Name == "false" || e.Name == "nil"
	case *ast.UnaryExpr:
		return isConstantExpr(e.X)
	case *ast.ParenExpr:
		return isConstantExpr(e.X)
	default:
		return false
	}
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name lowercased and without hyphens.
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}
