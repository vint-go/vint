package depguard

import "github.com/vint-go/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the depguard golangci-lint linter.
// depguard checks that package imports comply with allow/deny lists
// defined per rule group, optionally scoped to specific files.
type Migrator struct{}

func (*Migrator) Name() string {
	return "depguard"
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// Extract the "rules" map from depguard settings.
	var rulesMap map[string]any
	if settings != nil {
		if r, ok := settings["rules"]; ok {
			if rm, ok := r.(map[string]any); ok {
				rulesMap = rm
			}
		}
	}

	// If no rules are defined, enable all three rules with empty config.
	if len(rulesMap) == 0 {
		configs["lint/correctness/noDeniedImport"] = migrate.VintRuleConfig{}
		configs["lint/correctness/noFileScopedDeniedImport"] = migrate.VintRuleConfig{}
		configs["lint/correctness/noUnallowedImport"] = migrate.VintRuleConfig{}
		return configs, nil
	}

	// Analyze the rule groups to determine which vint rules to enable.
	hasDeny := false
	hasAllow := false
	hasFileScoped := false

	for _, groupCfg := range rulesMap {
		groupMap, ok := groupCfg.(map[string]any)
		if !ok {
			continue
		}

		if denyVal, ok := groupMap["deny"]; ok {
			if denyList, ok := denyVal.([]any); ok && len(denyList) > 0 {
				hasDeny = true

				// Check if there are file patterns (file-scoped deny).
				if filesVal, ok := groupMap["files"]; ok {
					if filesList, ok := filesVal.([]any); ok && len(filesList) > 0 {
						hasFileScoped = true
					}
				}
			}
		}

		if allowVal, ok := groupMap["allow"]; ok {
			if allowList, ok := allowVal.([]any); ok && len(allowList) > 0 {
				hasAllow = true
			}
		}
	}

	// Enable noDeniedImport when deny lists are present.
	if hasDeny {
		configs["lint/correctness/noDeniedImport"] = migrate.VintRuleConfig{
			Options: rulesMap,
		}
	}

	// Enable noFileScopedDeniedImport when deny lists with file patterns are present.
	if hasFileScoped {
		configs["lint/correctness/noFileScopedDeniedImport"] = migrate.VintRuleConfig{
			Options: map[string]any{
				"rules": rulesMap,
			},
		}
	}

	// Enable noUnallowedImport when allow lists are present.
	if hasAllow {
		configs["lint/correctness/noUnallowedImport"] = migrate.VintRuleConfig{
			Options: map[string]any{
				"rules": rulesMap,
			},
		}
	}

	return configs, nil
}
