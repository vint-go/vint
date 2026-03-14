package no_hardcoded_iv

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoHardcodedIvRule detects the use of hardcoded initialization vectors (IVs)
// or nonces for encryption operations.
type NoHardcodedIvRule struct{}

// Apply applies the rule to the given file.
func (r *NoHardcodedIvRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintHardcodedIV{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
		hardcodedVars: map[string]bool{},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoHardcodedIvRule) Name() string {
	return "noHardcodedIv"
}

// Group returns the rule group.
func (*NoHardcodedIvRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoHardcodedIvRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// cipherNewFunctions lists cipher package functions that take an IV as their second argument.
var cipherNewFunctions = map[string]bool{
	"NewCBCEncrypter": true,
	"NewCBCDecrypter": true,
	"NewCTR":          true,
	"NewOFB":          true,
	"NewCFBEncrypter": true,
	"NewCFBDecrypter": true,
}

type lintHardcodedIV struct {
	onFailure     func(lint.Failure)
	hardcodedVars map[string]bool
}

func (w *lintHardcodedIV) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		// Create a new scope for each function to track local variable assignments.
		funcWalker := &lintHardcodedIV{
			onFailure:     w.onFailure,
			hardcodedVars: copyMap(w.hardcodedVars),
		}
		if n.Body != nil {
			for _, stmt := range n.Body.List {
				ast.Walk(funcWalker, stmt)
			}
		}
		return nil

	case *ast.GenDecl:
		// Track package-level var declarations with hardcoded byte slices.
		if n.Tok == token.VAR {
			for _, spec := range n.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for i, name := range vs.Names {
					if i < len(vs.Values) && isHardcodedByteSlice(vs.Values[i]) {
						w.hardcodedVars[name.Name] = true
					}
				}
			}
		}

	case *ast.AssignStmt:
		// Track local variable assignments with hardcoded byte slices.
		for i, lhs := range n.Lhs {
			ident, ok := lhs.(*ast.Ident)
			if !ok {
				continue
			}
			if i < len(n.Rhs) && isHardcodedByteSlice(n.Rhs[i]) {
				w.hardcodedVars[ident.Name] = true
			} else if i < len(n.Rhs) {
				// If reassigned to a non-hardcoded value, remove it.
				delete(w.hardcodedVars, ident.Name)
			}
		}

	case *ast.CallExpr:
		w.checkCipherCall(n)
		w.checkAEADCall(n)
	}

	return w
}

// checkCipherCall checks calls like cipher.NewCBCEncrypter(block, iv)
// where iv is the second argument.
func (w *lintHardcodedIV) checkCipherCall(call *ast.CallExpr) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	// Check if it is cipher.<FuncName>
	if !astutils.IsIdent(sel.X, "cipher") {
		return
	}

	if !cipherNewFunctions[sel.Sel.Name] {
		return
	}

	// The IV is the second argument (index 1).
	if len(call.Args) < 2 {
		return
	}

	ivArg := call.Args[1]
	if w.isHardcodedValue(ivArg) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       call,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "hardcoded IV or nonce: use a cryptographically random value instead",
		})
	}
}

// checkAEADCall checks calls like gcm.Seal(dst, nonce, plaintext, additionalData)
// or gcm.Open(dst, nonce, ciphertext, additionalData) where nonce is the second argument.
func (w *lintHardcodedIV) checkAEADCall(call *ast.CallExpr) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	methodName := sel.Sel.Name
	if methodName != "Seal" && methodName != "Open" {
		return
	}

	// Check that the receiver looks like a variable (not a package).
	// We cannot fully determine the type without type checking, but
	// we can check if the method is called with the right argument count.
	// Seal(dst, nonce, plaintext, additionalData) = 4 args
	// Open(dst, nonce, ciphertext, additionalData) = 4 args
	if len(call.Args) != 4 {
		return
	}

	nonceArg := call.Args[1]
	if w.isHardcodedValue(nonceArg) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       call,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "hardcoded IV or nonce: use a cryptographically random value instead",
		})
	}
}

// isHardcodedValue checks if an expression is a hardcoded byte slice value,
// either directly or through a variable that was assigned a hardcoded value.
func (w *lintHardcodedIV) isHardcodedValue(expr ast.Expr) bool {
	// Direct hardcoded byte slice literal or string literal conversion.
	if isHardcodedByteSlice(expr) {
		return true
	}

	// Variable that was previously assigned a hardcoded value.
	if ident, ok := expr.(*ast.Ident); ok {
		return w.hardcodedVars[ident.Name]
	}

	return false
}

// isHardcodedByteSlice checks if an expression is a hardcoded byte slice.
// This includes:
//   - []byte("string literal")
//   - []byte{0, 1, 2, ...} with all constant elements
func isHardcodedByteSlice(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if ok {
		// Check for []byte("string literal") conversion.
		return isByteSliceConversion(call)
	}

	comp, ok := expr.(*ast.CompositeLit)
	if ok {
		// Check for []byte{0, 1, 2, ...} literal.
		return isByteSliceLiteral(comp)
	}

	return false
}

// isByteSliceConversion checks for []byte("hardcoded string") pattern.
func isByteSliceConversion(call *ast.CallExpr) bool {
	// The function should be []byte type conversion.
	arrayType, ok := call.Fun.(*ast.ArrayType)
	if !ok {
		return false
	}

	// Check that it is []byte (no length specified and element is byte).
	if arrayType.Len != nil {
		return false
	}
	if !astutils.IsIdent(arrayType.Elt, "byte") {
		return false
	}

	// Should have exactly one argument that is a string literal.
	if len(call.Args) != 1 {
		return false
	}

	lit, ok := call.Args[0].(*ast.BasicLit)
	return ok && lit.Kind == token.STRING
}

// isByteSliceLiteral checks for []byte{0, 1, 2} pattern with all constant elements.
func isByteSliceLiteral(comp *ast.CompositeLit) bool {
	// Check the type is []byte.
	arrayType, ok := comp.Type.(*ast.ArrayType)
	if !ok {
		return false
	}
	if arrayType.Len != nil {
		return false
	}
	if !astutils.IsIdent(arrayType.Elt, "byte") {
		return false
	}

	// All elements must be constant (basic literals).
	if len(comp.Elts) == 0 {
		return false
	}

	for _, elt := range comp.Elts {
		if _, ok := elt.(*ast.BasicLit); !ok {
			return false
		}
	}

	return true
}

func copyMap(m map[string]bool) map[string]bool {
	result := make(map[string]bool, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}
