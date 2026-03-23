package migrate

import (
	_ "embed"
	"fmt"
	"sort"

	"github.com/BurntSushi/toml"
	"github.com/vint-go/vint/config"
	"github.com/vint-go/vint/lint"
)

//go:embed registry/rules_registry.toml
var registryData string

// registryEntry matches the TOML structure of rules_registry.toml.
type registryEntry struct {
	Linter        string `toml:"linter"`
	Source        string `toml:"source"`
	ExtractorRule string `toml:"extractor_rule"`
	SubsumedBy    string `toml:"subsumed_by"`
}

type registryFile struct {
	Rules map[string]registryEntry `toml:"rules"`
}

// MappedRule represents a vint rule that corresponds to a golangci-lint linter.
type MappedRule struct {
	ArchName      string    // camelCase rule name, e.g. "noUncheckedError"
	FullVintPath  string    // e.g. "lint/correctness/noUncheckedError"
	ExtractorRule string    // e.g. "unchecked-error"
	Rule          lint.Rule // resolved rule instance (may be nil if not found)
}

// RuleRegistry provides lookups from golangci-lint linter names to vint rules.
type RuleRegistry struct {
	byLinter map[string][]MappedRule // linter name → mapped rules
}

// LoadRegistry parses the embedded rules_registry.toml and resolves rule
// instances from vint's compiled rule set.
func LoadRegistry() (*RuleRegistry, error) {
	var rf registryFile
	if err := toml.Unmarshal([]byte(registryData), &rf); err != nil {
		return nil, fmt.Errorf("parse rules registry: %w", err)
	}

	// Build a lookup from rule Name() → (fullPath, rule instance).
	allRules := config.GetAllRules()
	ruleByName := make(map[string]lint.Rule, len(allRules))
	for _, r := range allRules {
		ruleByName[r.Name()] = r
	}

	reg := &RuleRegistry{
		byLinter: make(map[string][]MappedRule),
	}

	// First pass: register non-subsumed rules.
	for archName, entry := range rf.Rules {
		if entry.SubsumedBy != "" {
			continue
		}

		var fullPath string
		var ruleInstance lint.Rule
		if r, ok := ruleByName[archName]; ok {
			fullPath = lint.FullRuleName(r)
			ruleInstance = r
		}

		reg.byLinter[entry.Linter] = append(reg.byLinter[entry.Linter], MappedRule{
			ArchName:      archName,
			FullVintPath:  fullPath,
			ExtractorRule: entry.ExtractorRule,
			Rule:          ruleInstance,
		})
	}

	// Second pass: register subsumed rules under their linter, resolved to
	// the subsumer's rule instance. This ensures that nolint conversion for
	// a linter whose rules are subsumed by another linter's rules still
	// works correctly.
	for _, entry := range rf.Rules {
		if entry.SubsumedBy == "" {
			continue
		}

		subsumerName := entry.SubsumedBy
		var fullPath string
		var ruleInstance lint.Rule
		if r, ok := ruleByName[subsumerName]; ok {
			fullPath = lint.FullRuleName(r)
			ruleInstance = r
		}

		// Only add if the subsumer resolves to an actual rule and isn't
		// already registered for this linter.
		if fullPath == "" {
			continue
		}
		alreadyRegistered := false
		for _, existing := range reg.byLinter[entry.Linter] {
			if existing.FullVintPath == fullPath {
				alreadyRegistered = true
				break
			}
		}
		if !alreadyRegistered {
			reg.byLinter[entry.Linter] = append(reg.byLinter[entry.Linter], MappedRule{
				ArchName:      subsumerName,
				FullVintPath:  fullPath,
				ExtractorRule: entry.ExtractorRule,
				Rule:          ruleInstance,
			})
		}
	}

	// Sort rules within each linter for deterministic output.
	for linter := range reg.byLinter {
		sort.Slice(reg.byLinter[linter], func(i, j int) bool {
			return reg.byLinter[linter][i].FullVintPath < reg.byLinter[linter][j].FullVintPath
		})
	}

	return reg, nil
}

// RulesForLinter returns the vint rules mapped to the given golangci-lint linter.
func (r *RuleRegistry) RulesForLinter(linterName string) []MappedRule {
	return r.byLinter[linterName]
}

// AllLinters returns a sorted list of all known golangci-lint linter names.
func (r *RuleRegistry) AllLinters() []string {
	names := make([]string, 0, len(r.byLinter))
	for name := range r.byLinter {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// RuleInstancesForLinter returns just the lint.Rule instances for a linter.
func (r *RuleRegistry) RuleInstancesForLinter(linterName string) []lint.Rule {
	mapped := r.byLinter[linterName]
	var rules []lint.Rule
	for _, m := range mapped {
		if m.Rule != nil {
			rules = append(rules, m.Rule)
		}
	}
	return rules
}

// AllRuleInstances returns all resolved lint.Rule instances across all linters.
func (r *RuleRegistry) AllRuleInstances() []lint.Rule {
	seen := make(map[string]bool)
	var rules []lint.Rule
	for _, mapped := range r.byLinter {
		for _, m := range mapped {
			if m.Rule != nil && !seen[m.FullVintPath] {
				seen[m.FullVintPath] = true
				rules = append(rules, m.Rule)
			}
		}
	}
	return rules
}
