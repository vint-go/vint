package no_ineffective_rand_call

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoIneffectiveRandCallRule detects calls like rand.Intn(1) that always return 0.
type NoIneffectiveRandCallRule struct{}

// randFuncs lists the math/rand functions that return 0 when called with argument 1.
var randFuncs = []string{
	"Intn",
	"Int31n",
	"Int63n",
	"IntN",
}

// Apply applies the rule to given file.
func (r *NoIneffectiveRandCallRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintIneffectiveRandCall{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoIneffectiveRandCallRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintIneffectiveRandCall{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoIneffectiveRandCallRule) Name() string {
	return "noIneffectiveRandCall"
}

// Group returns the rule group.
func (*NoIneffectiveRandCallRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoIneffectiveRandCallRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintIneffectiveRandCall struct {
	onFailure func(lint.Failure)
}

func (w *lintIneffectiveRandCall) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check if it's a call to one of the targeted rand functions
	var funcName string
	for _, name := range randFuncs {
		if astutils.IsPkgDotName(call.Fun, "rand", name) {
			funcName = name
			break
		}
	}

	if funcName == "" {
		return w
	}

	// Must have exactly one argument
	if len(call.Args) != 1 {
		return w
	}

	arg := call.Args[0]

	// Check if the argument is the integer literal 1
	lit, ok := arg.(*ast.BasicLit)
	if !ok {
		return w
	}

	if lit.Kind != token.INT || lit.Value != "1" {
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryLogic,
		Confidence: 1,
		Node:       call,
		Failure:    fmt.Sprintf("rand.%s(1) always returns 0, this is likely a logic error", funcName),
	})

	return w
}
