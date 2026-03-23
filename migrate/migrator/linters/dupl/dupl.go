package dupl

import "github.com/vint-go/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the dupl golangci-lint linter.
// dupl detects duplicate fragments of code.
type Migrator struct{}

func (*Migrator) Name() string {
	return "dupl"
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// noDuplicateCode is always enabled when dupl is enabled.
	opts := map[string]any{}

	if settings != nil {
		if v, ok := settings["threshold"]; ok {
			// Convert to int64 since the vint rule's Configure expects int64,
			// but YAML v3 decodes integers as int.
			switch n := v.(type) {
			case int:
				opts["threshold"] = int64(n)
			case int64:
				opts["threshold"] = n
			default:
				opts["threshold"] = v
			}
		}
	}

	configs["lint/complexity/noDuplicateCode"] = migrate.VintRuleConfig{
		Options: opts,
	}

	return configs, nil
}
