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

	// Serialize each function body into a token sequence, compute histograms,
	// and pre-compute rolling hash sets for the pair comparison phase.
	type serializedDecl struct {
		decl    *ast.FuncDecl
		tokens  []duplToken
		freq    [numTokenTypes]int
		hashSet map[uint64]struct{} // rolling hash fingerprints of all windows of size threshold
	}

	// Pre-compute base^(threshold-1) for rolling hash removal.
	const hashBase uint64 = 257
	basePow := uint64(1)
	for k := 0; k < r.threshold-1; k++ {
		basePow *= hashBase
	}

	serialized := make([]serializedDecl, len(decls))
	for i, d := range decls {
		tokens := serializeAST(d.Body)
		var freq [numTokenTypes]int
		for _, t := range tokens {
			freq[t]++
		}
		sd := serializedDecl{
			decl:   d,
			tokens: tokens,
			freq:   freq,
		}
		// Build hash set for functions large enough to participate in comparison.
		if len(tokens) >= r.threshold {
			sd.hashSet = buildHashSet(tokens, r.threshold, hashBase, basePow)
		}
		serialized[i] = sd
	}

	// Compare each pair of functions for structural similarity using rolling hash.
	reported := map[token.Pos]bool{}
	for i := 0; i < len(serialized); i++ {
		if serialized[i].hashSet == nil {
			continue
		}
		for j := i + 1; j < len(serialized); j++ {
			// Skip if this function was already reported as a duplicate.
			if reported[serialized[j].decl.Pos()] {
				continue
			}
			if serialized[j].hashSet == nil {
				continue
			}
			// Quick histogram pre-filter: if the token frequency overlap is
			// below threshold, no common substring of that length can exist.
			if tokenOverlap(&serialized[i].freq, &serialized[j].freq) < r.threshold {
				continue
			}
			if hasCommonWindow(serialized[i].tokens, serialized[j].hashSet, r.threshold, hashBase, basePow) {
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
							r.threshold,
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
// Uses uint8 since there are fewer than 256 token types, keeping token arrays compact.
type duplToken uint8

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

	numTokenTypes // sentinel: total number of token types
)

// serializeAST converts an AST node into a sequence of tokens representing
// structural shape, ignoring concrete values like identifiers and literals.
func serializeAST(node ast.Node) []duplToken {
	tokens := make([]duplToken, 0, 128)
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

// tokenOverlap computes the sum of min(freqA[t], freqB[t]) over all token
// types. This is an upper bound on the longest common substring length.
func tokenOverlap(a, b *[numTokenTypes]int) int {
	total := 0
	for i := 0; i < int(numTokenTypes); i++ {
		va, vb := a[i], b[i]
		if va < vb {
			total += va
		} else {
			total += vb
		}
	}
	return total
}

// buildHashSet computes rolling hash fingerprints for all windows of size
// windowSize in the token sequence, returning a set of hashes.
func buildHashSet(tokens []duplToken, windowSize int, base, basePow uint64) map[uint64]struct{} {
	n := len(tokens)
	numWindows := n - windowSize + 1
	set := make(map[uint64]struct{}, numWindows)

	// Compute hash for the first window.
	var h uint64
	for i := 0; i < windowSize; i++ {
		h = h*base + uint64(tokens[i])
	}
	set[h] = struct{}{}

	// Slide the window, updating the hash in O(1).
	for i := 1; i < numWindows; i++ {
		h = (h-uint64(tokens[i-1])*basePow)*base + uint64(tokens[i+windowSize-1])
		set[h] = struct{}{}
	}

	return set
}

// hasCommonWindow checks if token sequence a has any window of windowSize
// whose rolling hash appears in setB (pre-computed from another sequence).
func hasCommonWindow(a []duplToken, setB map[uint64]struct{}, windowSize int, base, basePow uint64) bool {
	numWindows := len(a) - windowSize + 1
	if numWindows <= 0 {
		return false
	}

	// Compute hash for the first window.
	var h uint64
	for i := 0; i < windowSize; i++ {
		h = h*base + uint64(a[i])
	}
	if _, ok := setB[h]; ok {
		return true
	}

	// Slide the window.
	for i := 1; i < numWindows; i++ {
		h = (h-uint64(a[i-1])*basePow)*base + uint64(a[i+windowSize-1])
		if _, ok := setB[h]; ok {
			return true
		}
	}

	return false
}
