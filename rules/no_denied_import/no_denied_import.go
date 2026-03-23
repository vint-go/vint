package no_denied_import

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/vint-go/vint/lint"
)

// NoDeniedImportRule reports when a Go source file imports a package
// that appears on a deny list configured for a matching rule group.
type NoDeniedImportRule struct {
	ruleGroups map[string]*deniedImportRuleGroup
}

type deniedImportEntry struct {
	pkg  string
	desc string
}

type deniedImportRuleGroup struct {
	name     string
	listMode string // "lax" (default) or "strict"
	files    []string
	deny     []deniedImportEntry
	allow    []deniedImportEntry
}

// Configure validates and applies the rule configuration.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *NoDeniedImportRule) Configure(arguments lint.Arguments) error {
	r.ruleGroups = map[string]*deniedImportRuleGroup{}

	if len(arguments) == 0 {
		return nil
	}

	args, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf("invalid argument to the noDeniedImport rule. Expecting a map, got %T", arguments[0])
	}

	for groupName, groupCfg := range args {
		groupMap, ok := groupCfg.(map[string]any)
		if !ok {
			return fmt.Errorf("invalid rule group %q configuration. Expecting a map, got %T", groupName, groupCfg)
		}

		group := &deniedImportRuleGroup{
			name:     groupName,
			listMode: "lax",
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

		if denyVal, ok := groupMap["deny"]; ok {
			denyList, ok := denyVal.([]any)
			if !ok {
				return fmt.Errorf("invalid deny list in rule group %q. Expecting a list, got %T", groupName, denyVal)
			}
			for _, d := range denyList {
				dMap, ok := d.(map[string]any)
				if !ok {
					return fmt.Errorf("invalid deny entry in rule group %q. Expecting a map, got %T", groupName, d)
				}
				entry := deniedImportEntry{}
				if pkg, ok := dMap["pkg"]; ok {
					if pkgStr, ok := pkg.(string); ok {
						entry.pkg = pkgStr
					}
				}
				if desc, ok := dMap["desc"]; ok {
					if descStr, ok := desc.(string); ok {
						entry.desc = descStr
					}
				}
				group.deny = append(group.deny, entry)
			}
		}

		if allowVal, ok := groupMap["allow"]; ok {
			allowList, ok := allowVal.([]any)
			if !ok {
				return fmt.Errorf("invalid allow list in rule group %q. Expecting a list, got %T", groupName, allowVal)
			}
			for _, a := range allowList {
				aMap, ok := a.(map[string]any)
				if !ok {
					return fmt.Errorf("invalid allow entry in rule group %q. Expecting a map, got %T", groupName, a)
				}
				entry := deniedImportEntry{}
				if pkg, ok := aMap["pkg"]; ok {
					if pkgStr, ok := pkg.(string); ok {
						entry.pkg = pkgStr
					}
				}
				if desc, ok := aMap["desc"]; ok {
					if descStr, ok := desc.(string); ok {
						entry.desc = descStr
					}
				}
				group.allow = append(group.allow, entry)
			}
		}

		r.ruleGroups[groupName] = group
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoDeniedImportRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, importSpec := range file.AST.Imports {
		if importSpec.Path == nil {
			continue
		}
		// Remove quotes from import path
		importPath := strings.Trim(importSpec.Path.Value, `"`)

		for _, group := range r.ruleGroups {
			if !r.fileMatchesGroup(file, group) {
				continue
			}

			if group.listMode == "strict" {
				// In strict mode: denied unless a more specific allow overrides
				if entry, denied := r.isDenied(importPath, group.deny); denied {
					if !r.isAllowedMoreSpecific(importPath, entry.pkg, group.allow) {
						msg := fmt.Sprintf("import '%s' is not allowed from list '%s'", importPath, group.name)
						if entry.desc != "" {
							msg += ": " + entry.desc
						}
						failures = append(failures, lint.Failure{
							Confidence: 1,
							Failure:    msg,
							Node:       importSpec,
							Category:   lint.FailureCategoryImports,
						})
					}
				}
			} else {
				// lax mode (default): all allowed unless denied
				if entry, denied := r.isDenied(importPath, group.deny); denied {
					msg := fmt.Sprintf("import '%s' is not allowed from list '%s'", importPath, group.name)
					if entry.desc != "" {
						msg += ": " + entry.desc
					}
					failures = append(failures, lint.Failure{
						Confidence: 1,
						Failure:    msg,
						Node:       importSpec,
						Category:   lint.FailureCategoryImports,
					})
				}
			}
		}
	}

	return failures
}

// isDenied checks if the importPath matches any entry in the deny list.
// Returns the matching entry and whether a match was found.
func (r *NoDeniedImportRule) isDenied(importPath string, denyList []deniedImportEntry) (deniedImportEntry, bool) {
	for _, entry := range denyList {
		if matchesPackage(importPath, entry.pkg) {
			return entry, true
		}
	}
	return deniedImportEntry{}, false
}

// isAllowedMoreSpecific checks if the importPath matches a more specific allow entry
// than the deny entry that matched.
func (r *NoDeniedImportRule) isAllowedMoreSpecific(importPath, denyPkg string, allowList []deniedImportEntry) bool {
	denyPrefix := strings.TrimSuffix(denyPkg, "$")
	for _, entry := range allowList {
		if matchesPackage(importPath, entry.pkg) {
			allowPrefix := strings.TrimSuffix(entry.pkg, "$")
			// The allow entry is more specific if it's longer (more specific path)
			if len(allowPrefix) > len(denyPrefix) {
				return true
			}
		}
	}
	return false
}

// matchesPackage checks if importPath matches the package pattern.
// If pattern ends with "$", it requires exact match.
// Otherwise, it uses prefix matching.
func matchesPackage(importPath, pattern string) bool {
	if strings.HasSuffix(pattern, "$") {
		// Exact match (no prefix matching)
		return importPath == strings.TrimSuffix(pattern, "$")
	}
	// Prefix matching: "foo/bar" matches "foo/bar" and "foo/bar/baz"
	return importPath == pattern || strings.HasPrefix(importPath, pattern+"/")
}

// fileMatchesGroup checks if the file matches any file pattern in the group.
// If no file patterns are specified, the group applies to all files.
func (r *NoDeniedImportRule) fileMatchesGroup(file *lint.File, group *deniedImportRuleGroup) bool {
	if len(group.files) == 0 {
		return true
	}

	for _, pattern := range group.files {
		negate := false
		p := pattern
		if strings.HasPrefix(p, "!") {
			negate = true
			p = p[1:]
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

	return false
}

// Name returns the rule name.
func (*NoDeniedImportRule) Name() string {
	return "noDeniedImport"
}

// Group returns the rule group.
func (*NoDeniedImportRule) Group() string {
	return "correctness"
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name lowercased and without hyphens.
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}
