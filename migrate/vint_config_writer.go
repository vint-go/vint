package migrate

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// vintYAMLOutput is the top-level structure written to vint.yaml.
type vintYAMLOutput struct {
	Settings map[string]VintRuleConfig `yaml:"settings"`
}

// WriteVintYAML writes the collected rule configs to the given path as vint.yaml.
func WriteVintYAML(path string, configs map[string]VintRuleConfig) error {
	data, err := RenderVintYAML(configs)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		return fmt.Errorf("write vint.yaml: %w", err)
	}
	return nil
}

// RenderVintYAML returns the vint.yaml content as a string.
func RenderVintYAML(configs map[string]VintRuleConfig) (string, error) {
	out := vintYAMLOutput{Settings: configs}
	data, err := yaml.Marshal(&out)
	if err != nil {
		return "", fmt.Errorf("marshal vint.yaml: %w", err)
	}
	return string(data), nil
}
