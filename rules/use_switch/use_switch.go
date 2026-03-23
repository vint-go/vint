package use_switch

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseSwitchRule detects repeated if-else chains and recommends converting them to switch statements.
type UseSwitchRule struct {
	minThreshold int
}

const defaultMinThreshold = 2

// Configure validates and applies the rule configuration.
func (r *UseSwitchRule) Configure(arguments lint.Arguments) error {
	if len(arguments) < 1 {
		r.minThreshold = defaultMinThreshold
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		// try direct int64 argument
		threshold, ok := lint.ToInt64(arguments[0])
		if !ok {
			return fmt.Errorf(`invalid argument to the "useSwitch" rule, expecting a k,v map or integer, got %T`, arguments[0])
		}
		r.minThreshold = int(threshold)
		return nil
	}

	r.minThreshold = defaultMinThreshold
	for k, v := range argKV {
		if isRuleOption(k, "minThreshold") {
			threshold, ok := lint.ToInt64(v)
			if !ok {
				return fmt.Errorf(`invalid configuration value for minThreshold in "useSwitch" rule; need integer but got %T`, v)
			}
			r.minThreshold = int(threshold)
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *UseSwitchRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintUseSwitch{
		minThreshold: r.minThreshold,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseSwitchRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintUseSwitch{
		minThreshold: r.minThreshold,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*UseSwitchRule) Name() string {
	return "useSwitch"
}

// Group returns the rule group.
func (*UseSwitchRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseSwitchRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintUseSwitch struct {
	minThreshold int
	onFailure    func(lint.Failure)
}

func (w *lintUseSwitch) Visit(node ast.Node) ast.Visitor {
	ifStmt, ok := node.(*ast.IfStmt)
	if !ok {
		return w
	}

	// Count the number of if-else-if branches in the chain
	branches := countIfElseBranches(ifStmt)

	if branches >= w.minThreshold {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Category:   lint.FailureCategoryStyle,
			Failure:    fmt.Sprintf("could replace if-else chain with %d branches by a switch statement", branches),
			Node:       ifStmt,
		})
		// Don't recurse into the else branches to avoid double-reporting
		// We only want to visit the bodies and init statements
		return nil
	}

	return w
}

// countIfElseBranches counts the total number of if/else-if branches in an if-else chain.
// It only counts chains where each if/else-if has no init statement (simple conditions).
func countIfElseBranches(ifStmt *ast.IfStmt) int {
	count := 1 // count the initial "if"
	current := ifStmt

	for {
		elseStmt := current.Else
		if elseStmt == nil {
			break
		}

		elseIf, ok := elseStmt.(*ast.IfStmt)
		if !ok {
			// It's a plain "else" block, not an "else if"
			break
		}

		count++
		current = elseIf
	}

	return count
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name lowercased and without hyphens.
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}
