package no_unallowed_import

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/vint-go/vint/lint"
)

// NoUnallowedImportRule reports when a Go source file imports a package
// that is not present in the allow list configured for the matching
// depguard rule group.
type NoUnallowedImportRule struct {
	ruleGroups map[string]*unallowedImportRuleGroup
}

type unallowedImportRuleGroup struct {
	name     string
	listMode string // "strict", "lax", or "original"
	files    []string
	allow    []string
	deny     []string
}

// Configure validates and applies the rule configuration.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *NoUnallowedImportRule) Configure(arguments lint.Arguments) error {
	r.ruleGroups = map[string]*unallowedImportRuleGroup{}

	if len(arguments) == 0 {
		return nil
	}

	args, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf("invalid argument to the noUnallowedImport rule. Expecting a map, got %T", arguments[0])
	}

	// The config can be either a "rules" key containing the groups, or groups directly
	rulesMap := args
	if rVal, ok := args["rules"]; ok {
		rm, ok := rVal.(map[string]any)
		if !ok {
			return fmt.Errorf("invalid 'rules' value in noUnallowedImport rule. Expecting a map, got %T", rVal)
		}
		rulesMap = rm
	}

	for groupName, groupCfg := range rulesMap {
		groupMap, ok := groupCfg.(map[string]any)
		if !ok {
			return fmt.Errorf("invalid rule group %q configuration. Expecting a map, got %T", groupName, groupCfg)
		}

		group := &unallowedImportRuleGroup{
			name:     groupName,
			listMode: "original",
		}

		if lm, ok := groupMap["list-mode"]; ok {
			if lmStr, ok := lm.(string); ok {
				group.listMode = lmStr
			}
		}

		if filesVal, ok := groupMap["files"]; ok {
			if filesList, ok := filesVal.([]any); ok {
				for _, f := range filesList {
					if fStr, ok := f.(string); ok {
						group.files = append(group.files, fStr)
					}
				}
			}
		}

		if allowVal, ok := groupMap["allow"]; ok {
			allowList, ok := allowVal.([]any)
			if !ok {
				return fmt.Errorf("invalid allow list in rule group %q. Expecting a list, got %T", groupName, allowVal)
			}
			for _, a := range allowList {
				switch av := a.(type) {
				case string:
					group.allow = append(group.allow, av)
				case map[string]any:
					if pkg, ok := av["pkg"]; ok {
						if pkgStr, ok := pkg.(string); ok {
							group.allow = append(group.allow, pkgStr)
						}
					}
				}
			}
		}

		if denyVal, ok := groupMap["deny"]; ok {
			denyList, ok := denyVal.([]any)
			if !ok {
				return fmt.Errorf("invalid deny list in rule group %q. Expecting a list, got %T", groupName, denyVal)
			}
			for _, d := range denyList {
				switch dv := d.(type) {
				case string:
					group.deny = append(group.deny, dv)
				case map[string]any:
					if pkg, ok := dv["pkg"]; ok {
						if pkgStr, ok := pkg.(string); ok {
							group.deny = append(group.deny, pkgStr)
						}
					}
				}
			}
		}

		r.ruleGroups[groupName] = group
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoUnallowedImportRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, importSpec := range file.AST.Imports {
		if importSpec.Path == nil {
			continue
		}
		importPath := strings.Trim(importSpec.Path.Value, `"`)

		for _, group := range r.ruleGroups {
			if !r.fileMatchesGroup(file, group) {
				continue
			}

			if !r.isImportAllowed(importPath, group) {
				msg := fmt.Sprintf("import '%s' is not allowed from list '%s'", importPath, group.name)
				failures = append(failures, lint.Failure{
					Confidence: 1,
					Failure:    msg,
					Node:       importSpec,
					Category:   lint.FailureCategoryImports,
				})
			}
		}
	}

	return failures
}

// isImportAllowed checks whether an import is allowed based on the group's
// list mode and allow/deny lists.
func (r *NoUnallowedImportRule) isImportAllowed(importPath string, group *unallowedImportRuleGroup) bool {
	switch group.listMode {
	case "strict":
		// In strict mode, every import must be explicitly allowed.
		// If not in the allow list, it is denied.
		return r.matchesAllowList(importPath, group.allow)
	case "lax":
		// In lax mode, packages are allowed by default; only the deny list is enforced.
		// The allow list is not relevant in lax mode for this rule.
		return true
	default:
		// "original" mode: if an allow list is defined, packages not matching it are denied.
		if len(group.allow) > 0 {
			return r.matchesAllowList(importPath, group.allow)
		}
		return true
	}
}

// matchesAllowList checks if the import path matches any entry in the allow list.
func (r *NoUnallowedImportRule) matchesAllowList(importPath string, allowList []string) bool {
	for _, pattern := range allowList {
		if pattern == "$gostd" {
			if unallowedImportIsGoStdLib(importPath) {
				return true
			}
			continue
		}
		if unallowedImportMatchesPackage(importPath, pattern) {
			return true
		}
	}
	return false
}

// unallowedImportMatchesPackage checks if importPath matches the package pattern.
// If pattern ends with "$", it requires exact match.
// Otherwise, it uses prefix matching.
func unallowedImportMatchesPackage(importPath, pattern string) bool {
	if strings.HasSuffix(pattern, "$") {
		// Exact match (no prefix matching)
		return importPath == strings.TrimSuffix(pattern, "$")
	}
	// Prefix matching: "foo/bar" matches "foo/bar" and "foo/bar/baz"
	return importPath == pattern || strings.HasPrefix(importPath, pattern+"/")
}

// unallowedImportIsGoStdLib checks if the import path looks like a Go standard library package.
// Standard library packages do not contain a dot in the first path segment.
func unallowedImportIsGoStdLib(importPath string) bool {
	firstSlash := strings.Index(importPath, "/")
	firstSegment := importPath
	if firstSlash >= 0 {
		firstSegment = importPath[:firstSlash]
	}
	return !strings.Contains(firstSegment, ".")
}

// fileMatchesGroup checks if the file matches any file pattern in the group.
// If no file patterns are specified, the group applies to all files.
func (r *NoUnallowedImportRule) fileMatchesGroup(file *lint.File, group *unallowedImportRuleGroup) bool {
	if len(group.files) == 0 {
		return true
	}

	hasPositivePattern := false

	for _, pattern := range group.files {
		negate := false
		p := pattern
		if strings.HasPrefix(p, "!") {
			negate = true
			p = p[1:]
		}

		if !negate {
			hasPositivePattern = true
		}

		matched := false
		switch p {
		case "$all":
			matched = true
		case "$test":
			matched = file.IsTest()
		default:
			// Use glob matching on the filename
			m, err := filepath.Match(p, filepath.Base(file.Name))
			if err == nil && m {
				matched = true
			}
			// Also try matching against the full path
			if !matched {
				m, err = filepath.Match(p, file.Name)
				if err == nil && m {
					matched = true
				}
			}
		}

		if negate {
			if matched {
				return false
			}
		} else if matched {
			return true
		}
	}

	// If only negated patterns were specified and none matched (excluded) the file,
	// the file is considered to match the group.
	if !hasPositivePattern {
		return true
	}

	return false
}

// Name returns the rule name.
func (*NoUnallowedImportRule) Name() string {
	return "noUnallowedImport"
}

// Group returns the rule group.
func (*NoUnallowedImportRule) Group() string {
	return "correctness"
}
