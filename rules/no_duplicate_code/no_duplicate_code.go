package no_duplicate_code

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/strowk/vint/lint"
)

// NoDuplicateCodeRule detects duplicate fragments of code.
type NoDuplicateCodeRule struct {
	threshold int
}

const defaultDuplicateCodeThreshold = 150

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *NoDuplicateCodeRule) Configure(arguments lint.Arguments) error {
	if len(arguments) < 1 {
		r.threshold = defaultDuplicateCodeThreshold
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		// try direct int64 argument
		threshold, ok := arguments[0].(int64)
		if !ok {
			return fmt.Errorf(`invalid argument to the "noDuplicateCode" rule, expecting a k,v map or int64, got %T`, arguments[0])
		}
		r.threshold = int(threshold)
		return nil
	}

	for k, v := range argKV {
		if isRuleOption(k, "threshold") {
			threshold, ok := v.(int64)
			if !ok || threshold < 0 {
				return fmt.Errorf(`invalid configuration value for threshold in "noDuplicateCode" rule; need positive int64 but got %T`, v)
			}
			r.threshold = int(threshold)
		}
	}

	if r.threshold == 0 {
		r.threshold = defaultDuplicateCodeThreshold
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoDuplicateCodeRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if r.threshold == 0 {
		r.threshold = defaultDuplicateCodeThreshold
	}

	var failures []lint.Failure

	// Collect all top-level function/method declarations.
	var decls []*ast.FuncDecl
	for _, d := range file.AST.Decls {
		if fn, ok := d.(*ast.FuncDecl); ok && fn.Body != nil {
			decls = append(decls, fn)
		}
	}

	if len(decls) < 2 {
		return nil
	}

	// Serialize each function body into a token sequence.
	type serializedDecl struct {
		decl   *ast.FuncDecl
		tokens []duplToken
	}

	serialized := make([]serializedDecl, len(decls))
	for i, d := range decls {
		serialized[i] = serializedDecl{
			decl:   d,
			tokens: serializeAST(d.Body),
		}
	}

	// Compare each pair of functions for structural similarity.
	reported := map[token.Pos]bool{}
	for i := 0; i < len(serialized); i++ {
		for j := i + 1; j < len(serialized); j++ {
			lcsLen := longestCommonSubsequenceLen(serialized[i].tokens, serialized[j].tokens)
			if lcsLen >= r.threshold {
				// Report on the second (later) function to avoid double-reporting.
				pos := serialized[j].decl.Pos()
				if !reported[pos] {
					reported[pos] = true
					failures = append(failures, lint.Failure{
						Confidence: 1,
						Category:   lint.FailureCategoryComplexity,
						Failure: fmt.Sprintf(
							"duplicate code detected: %s and %s share %d tokens of identical structure",
							serialized[i].decl.Name.Name,
							serialized[j].decl.Name.Name,
							lcsLen,
						),
						Node: serialized[j].decl,
					})
				}
			}
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoDuplicateCodeRule) Name() string {
	return "noDuplicateCode"
}

// Group returns the rule group.
func (*NoDuplicateCodeRule) Group() string {
	return "complexity"
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name lowercased and without hyphens.
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}

// duplToken represents a serialized AST node type, abstracting away concrete values.
type duplToken int

const (
	tokIdent duplToken = iota
	tokBasicLit
	tokCompositeLit
	tokSelectorExpr
	tokIndexExpr
	tokSliceExpr
	tokTypeAssertExpr
	tokCallExpr
	tokStarExpr
	tokUnaryExpr
	tokBinaryExpr
	tokParenExpr
	tokKeyValueExpr
	tokArrayType
	tokMapType
	tokChanType
	tokFuncType
	tokFuncLit
	tokAssignStmt
	tokBlockStmt
	tokBranchStmt
	tokCaseClause
	tokCommClause
	tokDeclStmt
	tokDeferStmt
	tokEmptyStmt
	tokExprStmt
	tokForStmt
	tokGoStmt
	tokIfStmt
	tokIncDecStmt
	tokLabeledStmt
	tokRangeStmt
	tokReturnStmt
	tokSelectStmt
	tokSendStmt
	tokSwitchStmt
	tokTypeSwitchStmt
	tokValueSpec
	tokField
	tokFieldList
	tokEllipsis
	tokInterfaceType
	tokStructType
	tokFuncDecl
	tokGenDecl
)

// serializeAST converts an AST node into a sequence of tokens representing
// structural shape, ignoring concrete values like identifiers and literals.
func serializeAST(node ast.Node) []duplToken {
	var tokens []duplToken
	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		switch n.(type) {
		case *ast.Ident:
			tokens = append(tokens, tokIdent)
		case *ast.BasicLit:
			tokens = append(tokens, tokBasicLit)
		case *ast.CompositeLit:
			tokens = append(tokens, tokCompositeLit)
		case *ast.SelectorExpr:
			tokens = append(tokens, tokSelectorExpr)
		case *ast.IndexExpr:
			tokens = append(tokens, tokIndexExpr)
		case *ast.SliceExpr:
			tokens = append(tokens, tokSliceExpr)
		case *ast.TypeAssertExpr:
			tokens = append(tokens, tokTypeAssertExpr)
		case *ast.CallExpr:
			tokens = append(tokens, tokCallExpr)
		case *ast.StarExpr:
			tokens = append(tokens, tokStarExpr)
		case *ast.UnaryExpr:
			tokens = append(tokens, tokUnaryExpr)
		case *ast.BinaryExpr:
			tokens = append(tokens, tokBinaryExpr)
		case *ast.ParenExpr:
			tokens = append(tokens, tokParenExpr)
		case *ast.KeyValueExpr:
			tokens = append(tokens, tokKeyValueExpr)
		case *ast.ArrayType:
			tokens = append(tokens, tokArrayType)
		case *ast.MapType:
			tokens = append(tokens, tokMapType)
		case *ast.ChanType:
			tokens = append(tokens, tokChanType)
		case *ast.FuncType:
			tokens = append(tokens, tokFuncType)
		case *ast.FuncLit:
			tokens = append(tokens, tokFuncLit)
		case *ast.AssignStmt:
			tokens = append(tokens, tokAssignStmt)
		case *ast.BlockStmt:
			tokens = append(tokens, tokBlockStmt)
		case *ast.BranchStmt:
			tokens = append(tokens, tokBranchStmt)
		case *ast.CaseClause:
			tokens = append(tokens, tokCaseClause)
		case *ast.CommClause:
			tokens = append(tokens, tokCommClause)
		case *ast.DeclStmt:
			tokens = append(tokens, tokDeclStmt)
		case *ast.DeferStmt:
			tokens = append(tokens, tokDeferStmt)
		case *ast.EmptyStmt:
			tokens = append(tokens, tokEmptyStmt)
		case *ast.ExprStmt:
			tokens = append(tokens, tokExprStmt)
		case *ast.ForStmt:
			tokens = append(tokens, tokForStmt)
		case *ast.GoStmt:
			tokens = append(tokens, tokGoStmt)
		case *ast.IfStmt:
			tokens = append(tokens, tokIfStmt)
		case *ast.IncDecStmt:
			tokens = append(tokens, tokIncDecStmt)
		case *ast.LabeledStmt:
			tokens = append(tokens, tokLabeledStmt)
		case *ast.RangeStmt:
			tokens = append(tokens, tokRangeStmt)
		case *ast.ReturnStmt:
			tokens = append(tokens, tokReturnStmt)
		case *ast.SelectStmt:
			tokens = append(tokens, tokSelectStmt)
		case *ast.SendStmt:
			tokens = append(tokens, tokSendStmt)
		case *ast.SwitchStmt:
			tokens = append(tokens, tokSwitchStmt)
		case *ast.TypeSwitchStmt:
			tokens = append(tokens, tokTypeSwitchStmt)
		case *ast.ValueSpec:
			tokens = append(tokens, tokValueSpec)
		case *ast.Field:
			tokens = append(tokens, tokField)
		case *ast.FieldList:
			tokens = append(tokens, tokFieldList)
		case *ast.Ellipsis:
			tokens = append(tokens, tokEllipsis)
		case *ast.InterfaceType:
			tokens = append(tokens, tokInterfaceType)
		case *ast.StructType:
			tokens = append(tokens, tokStructType)
		case *ast.FuncDecl:
			tokens = append(tokens, tokFuncDecl)
		case *ast.GenDecl:
			tokens = append(tokens, tokGenDecl)
		}
		return true
	})
	return tokens
}

// longestCommonSubsequenceLen computes the length of the longest common
// contiguous subsequence (substring) between two token sequences.
func longestCommonSubsequenceLen(a, b []duplToken) int {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}

	maxLen := 0

	// Use a rolling-row DP approach for longest common substring.
	// dp[j] = length of the longest common suffix ending at a[i-1] and b[j-1].
	dp := make([]int, len(b)+1)

	for i := 1; i <= len(a); i++ {
		// Process in reverse to avoid overwriting values we still need.
		for j := len(b); j >= 1; j-- {
			if a[i-1] == b[j-1] {
				if i == 1 || j == 1 {
					dp[j] = 1
				} else {
					// We need dp[j-1] from the previous row (i-1).
					// Since we process j in reverse, dp[j-1] still has the value from row i-1.
					dp[j] = dp[j-1] + 1
				}
				if dp[j] > maxLen {
					maxLen = dp[j]
				}
			} else {
				dp[j] = 0
			}
		}
	}

	return maxLen
}
