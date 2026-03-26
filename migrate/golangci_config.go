package migrate

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// GolangciConfig represents the relevant parts of a .golangci.yml file.
type GolangciConfig struct {
	Version string     `yaml:"version"`
	Linters LintersCfg `yaml:"linters"`
	// V1 only: top-level linters-settings.
	LintersSettings map[string]map[string]any `yaml:"linters-settings"`
}

// LintersCfg holds linter enable/disable configuration.
type LintersCfg struct {
	Enable     []string `yaml:"enable"`
	Disable    []string `yaml:"disable"`
	DisableAll bool     `yaml:"disable-all"`
	EnableAll  bool     `yaml:"enable-all"`
	// V2 only: settings nested under linters.
	Settings   map[string]map[string]any `yaml:"settings"`
	Exclusions ExclusionsCfg             `yaml:"exclusions"`
}

// ExclusionRuleEntry represents a single entry in the exclusions.rules array.
type ExclusionRuleEntry struct {
	Linters []string `yaml:"linters"`
	Path    string   `yaml:"path"`
	Text    string   `yaml:"text"`
}

// ExclusionsCfg holds the exclusion configuration from golangci-lint.
type ExclusionsCfg struct {
	Presets []string             `yaml:"presets"`
	Rules   []ExclusionRuleEntry `yaml:"rules"`
}

// LoadGolangciConfig reads and parses a .golangci.yml file.
// Only version "2" configs are supported.
func LoadGolangciConfig(path string) (*GolangciConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read golangci config: %w", err)
	}

	var cfg GolangciConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse golangci config: %w", err)
	}

	if cfg.Version != "2" {
		return nil, fmt.Errorf("unsupported golangci-lint config version %q (only version \"2\" is supported)", cfg.Version)
	}

	return &cfg, nil
}

// EnabledLinters returns the set of linters that are explicitly enabled.
func (c *GolangciConfig) EnabledLinters() []string {
	return c.Linters.Enable
}

// LinterSettings returns the settings for a specific linter.
// In v2, settings are under linters.settings.<linter>.
func (c *GolangciConfig) LinterSettings(linterName string) map[string]any {
	if c.Linters.Settings != nil {
		return c.Linters.Settings[linterName]
	}
	return nil
}
