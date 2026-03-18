package use_buffer_string_or_bytes

import (
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseBufferStringOrBytesRule detects string(buf.Bytes()) conversions on bytes.Buffer
// that can be simplified to buf.String().
type UseBufferStringOrBytesRule struct{}

// Apply applies the rule to given file.
func (r *UseBufferStringOrBytesRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintBufferStringOrBytes{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseBufferStringOrBytesRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintBufferStringOrBytes{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*UseBufferStringOrBytesRule) Name() string {
	return "useBufferStringOrBytes"
}

// Group returns the rule group.
func (*UseBufferStringOrBytesRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseBufferStringOrBytesRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule needs type information.
func (*UseBufferStringOrBytesRule) RequiresTypecheck() bool {
	return true
}

type lintBufferStringOrBytes struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintBufferStringOrBytes) Visit(node ast.Node) ast.Visitor {
	// Look for string(...) type conversion calls
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check if it's a string() type conversion
	ident, ok := call.Fun.(*ast.Ident)
	if !ok {
		return w
	}
	if ident.Name != "string" {
		return w
	}
	if len(call.Args) != 1 {
		return w
	}

	// Check if the argument is a method call to .Bytes()
	innerCall, ok := call.Args[0].(*ast.CallExpr)
	if !ok {
		return w
	}

	sel, ok := innerCall.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	if sel.Sel.Name != "Bytes" {
		return w
	}

	// Check if the receiver is of type bytes.Buffer or *bytes.Buffer
	recvType := w.pkg.TypeOf(sel.X)
	if recvType == nil {
		return w
	}

	if isBytesBuffer(recvType) {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryStyle,
			Confidence: 1,
			Node:       call,
			Failure:    "use buf.String() instead of string(buf.Bytes())",
		})
	}

	return w
}

// isBytesBuffer checks if the type is bytes.Buffer or *bytes.Buffer.
func isBytesBuffer(t types.Type) bool {
	// Dereference pointer if needed
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}

	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	obj := named.Obj()
	if obj == nil {
		return false
	}

	return obj.Name() == "Buffer" && obj.Pkg() != nil && obj.Pkg().Path() == "bytes"
}
