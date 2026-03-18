package no_overlapping_encoder_slice

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoOverlappingEncoderSliceRule detects overlapping byte slices passed to encoder functions.
type NoOverlappingEncoderSliceRule struct{}

// Apply applies the rule to given file.
func (r *NoOverlappingEncoderSliceRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintOverlappingEncoder{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoOverlappingEncoderSliceRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintOverlappingEncoder{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoOverlappingEncoderSliceRule) Name() string {
	return "noOverlappingEncoderSlice"
}

// Group returns the rule group.
func (*NoOverlappingEncoderSliceRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoOverlappingEncoderSliceRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintOverlappingEncoder struct {
	onFailure func(lint.Failure)
}

// encoderFunc describes a function that takes dst and src slice arguments.
type encoderFunc struct {
	pkg      string
	name     string
	dstIndex int
	srcIndex int
}

// packageLevelFuncs are package-level functions like hex.Encode(dst, src).
var packageLevelFuncs = []encoderFunc{
	{pkg: "hex", name: "Encode", dstIndex: 0, srcIndex: 1},
	{pkg: "hex", name: "Decode", dstIndex: 0, srcIndex: 1},
	{pkg: "ascii85", name: "Encode", dstIndex: 0, srcIndex: 1},
	{pkg: "ascii85", name: "Decode", dstIndex: 0, srcIndex: 1},
}

// methodNames are method names on encoding types like base64.StdEncoding.Encode(dst, src).
var methodNames = map[string]struct {
	dstIndex int
	srcIndex int
}{
	"Encode": {dstIndex: 0, srcIndex: 1},
	"Decode": {dstIndex: 0, srcIndex: 1},
}

func (w *lintOverlappingEncoder) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check package-level functions: hex.Encode, hex.Decode, ascii85.Encode, ascii85.Decode
	for _, f := range packageLevelFuncs {
		if astutils.IsPkgDotName(call.Fun, f.pkg, f.name) {
			w.checkOverlap(call, f.dstIndex, f.srcIndex)
			return w
		}
	}

	// Check method calls: base64.StdEncoding.Encode(dst, src), base32.StdEncoding.Encode(dst, src), etc.
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if info, exists := methodNames[sel.Sel.Name]; exists {
			// Check if receiver is a selector from base64 or base32 package
			if isEncodingReceiver(sel.X) {
				w.checkOverlap(call, info.dstIndex, info.srcIndex)
				return w
			}
		}
	}

	return w
}

// isEncodingReceiver checks if the expression looks like base64.StdEncoding, base32.StdEncoding, etc.
func isEncodingReceiver(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	// Check for known encoding packages
	switch ident.Name {
	case "base64", "base32":
		return true
	}
	return false
}

// checkOverlap checks if dst and src arguments potentially overlap.
func (w *lintOverlappingEncoder) checkOverlap(call *ast.CallExpr, dstIndex, srcIndex int) {
	if len(call.Args) <= dstIndex || len(call.Args) <= srcIndex {
		return
	}

	dst := call.Args[dstIndex]
	src := call.Args[srcIndex]

	dstBase := baseIdent(dst)
	srcBase := baseIdent(src)

	if dstBase == "" || srcBase == "" {
		return
	}

	if dstBase == srcBase {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       call,
			Failure:    "overlapping dst and src slices passed to an encoder",
		})
	}
}

// baseIdent extracts the base identifier name from a slice expression.
// For example:
//   - buf       -> "buf"
//   - buf[:50]  -> "buf"
//   - buf[1:10] -> "buf"
//   - x.field   -> "" (skip complex expressions)
func baseIdent(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SliceExpr:
		return baseIdent(e.X)
	case *ast.IndexExpr:
		return baseIdent(e.X)
	default:
		return ""
	}
}
