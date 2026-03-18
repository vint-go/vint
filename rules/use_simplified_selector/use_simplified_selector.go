package use_simplified_selector

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseSimplifiedSelectorRule detects selector expressions where an embedded
// field name can be omitted.  For example, s.Mutex.Lock() can be written
// as s.Lock() when Mutex is an embedded (anonymous) field and there is no
// ambiguity.
type UseSimplifiedSelectorRule struct{}

// Apply applies the rule to the given file.
func (r *UseSimplifiedSelectorRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintSimplifiedSelector{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseSimplifiedSelectorRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintSimplifiedSelector{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*UseSimplifiedSelectorRule) Name() string {
	return "useSimplifiedSelector"
}

// Group returns the rule group.
func (*UseSimplifiedSelectorRule) Group() string {
	return "style"
}

// RequiresTypecheck returns true because the rule needs type information.
func (*UseSimplifiedSelectorRule) RequiresTypecheck() bool {
	return true
}

// CacheTier returns the cache tier for this rule.
func (*UseSimplifiedSelectorRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

type lintSimplifiedSelector struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintSimplifiedSelector) Visit(node ast.Node) ast.Visitor {
	outerSel, ok := node.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	// We need outerSel.X to also be a SelectorExpr (the embedded field access).
	innerSel, ok := outerSel.X.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	info := w.pkg.TypesInfo()
	if info == nil || info.Selections == nil {
		return w
	}

	// Check if the inner selector is a field access on an embedded (anonymous) field.
	innerSelection, ok := info.Selections[innerSel]
	if !ok {
		return w
	}

	// Must be a field value (not a method value or method expression).
	if innerSelection.Kind() != types.FieldVal {
		return w
	}

	// The field must be embedded (anonymous).
	obj := innerSelection.Obj()
	field, ok := obj.(*types.Var)
	if !ok || !field.Embedded() {
		return w
	}

	// Now check if the outer selector (the method/field accessed on the embedded type)
	// can also be accessed directly on the receiver type (i.e., it is promoted).
	receiverType := w.pkg.TypesInfo().TypeOf(innerSel.X)
	if receiverType == nil {
		return w
	}

	targetName := outerSel.Sel.Name

	// Look up the field/method directly on the receiver type.
	directObj, _, _ := types.LookupFieldOrMethod(receiverType, true, w.pkg.TypesPkg(), targetName)
	if directObj == nil {
		return w
	}

	// Verify it resolves to the same underlying object as the explicit path.
	outerSelection, ok := info.Selections[outerSel]
	if !ok {
		return w
	}
	if directObj != outerSelection.Obj() {
		// Ambiguity or different resolution — don't suggest simplification.
		return w
	}

	embeddedName := innerSel.Sel.Name

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       outerSel,
		Category:   lint.FailureCategoryStyle,
		Failure:    fmt.Sprintf("embedded field %s can be omitted from selector", embeddedName),
	})

	return w
}
