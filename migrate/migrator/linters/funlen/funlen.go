package funlen

import "github.com/vint-go/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the funlen golangci-lint linter.
// funlen checks for long functions (by line count and statement count).
type Migrator struct{}

func (*Migrator) Name() string {
	return "funlen"
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// noLongFunctions is always enabled when funlen is enabled.
	// golangci-lint defaults: lines=60, ignore-comments=true
	// vint defaults:          lines=60, ignoreComments=false
	// We always pass ignoreComments to match golangci-lint's default of true.
	longFuncOpts := map[string]any{
		"ignoreComments": true,
	}

	if settings != nil {
		if v, ok := settings["lines"]; ok {
			longFuncOpts["lines"] = toInt64(v)
		}
		if v, ok := settings["ignore-comments"]; ok {
			longFuncOpts["ignoreComments"] = v
		}
	}

	configs["lint/complexity/noLongFunctions"] = migrate.VintRuleConfig{
		Options: longFuncOpts,
	}

	// noExcessiveStatements is always enabled when funlen is enabled.
	stmtsOpts := map[string]any{}

	if settings != nil {
		if v, ok := settings["statements"]; ok {
			stmtsOpts["statements"] = toInt64(v)
		}
	}

	configs["lint/complexity/noExcessiveStatements"] = migrate.VintRuleConfig{
		Options: stmtsOpts,
	}

	return configs, nil
}

// toInt64 converts numeric values to int64. YAML parsers typically decode
// integers as int, but vint rules expect int64.
func toInt64(v any) int64 {
	switch n := v.(type) {
	case int:
		return int64(n)
	case int64:
		return n
	case float64:
		return int64(n)
	default:
		return 0
	}
}
