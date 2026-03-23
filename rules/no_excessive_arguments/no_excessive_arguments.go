package no_excessive_arguments

import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoExcessiveArgumentsRule lints the number of arguments a function can receive.
type NoExcessiveArgumentsRule struct {
	max int
}

const defaultArgumentsLimit = 8

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *NoExcessiveArgumentsRule) Configure(arguments lint.Arguments) error {
	if len(arguments) < 1 {
		r.max = defaultArgumentsLimit
		return nil
	}

	maxArguments, ok := arguments[0].(int64) // Alt. non panicking version
	if !ok {
		return errors.New(`invalid value passed as argument number to the "noExcessiveArguments" rule`)
	}
	r.max = int(maxArguments)
	return nil
}

// Apply applies the rule to given file.
func (r *NoExcessiveArgumentsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		numParams := 0
		for _, l := range funcDecl.Type.Params.List {
			numParams += len(l.Names)
		}

		if numParams <= r.max {
			continue
		}

		failures = append(failures, lint.Failure{
			Confidence: 1,
			Failure:    fmt.Sprintf("maximum number of arguments per function exceeded; max %d but got %d", r.max, numParams),
			Node:       funcDecl.Type,
		})
	}

	return failures
}

// Name returns the rule name.
func (*NoExcessiveArgumentsRule) Name() string {
	return "noExcessiveArguments"
}

// Group returns the rule group.
func (*NoExcessiveArgumentsRule) Group() string {
	return "complexity"
}

// CacheTier returns the cache tier for this rule.
func (*NoExcessiveArgumentsRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
