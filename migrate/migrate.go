package migrate

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var migrators = map[string]LinterMigrator{}

// RegisterMigrator registers a linter-specific migrator.
func RegisterMigrator(m LinterMigrator) {
	migrators[m.Name()] = m
}

// MigrateResult holds the full result of a migration.
type MigrateResult struct {
	// VintYAML is the generated vint.yaml content.
	VintYAML string
	// ConvertedFiles maps file paths to their converted content.
	ConvertedFiles map[string]string
	// Warnings holds non-fatal issues encountered during migration.
	Warnings []string
	// UnknownLinters lists linters that couldn't be migrated.
	UnknownLinters []string
}

// RunMigrate is the main entry point for the migrate subcommand.
func RunMigrate(args []string) {
	fs := flag.NewFlagSet("migrate", flag.ExitOnError)
	configPath := fs.String("config", ".golangci.yml", "path to golangci-lint config file")
	dir := fs.String("dir", ".", "project directory to scan for .go files")
	outputConfig := fs.String("output", "vint.yaml", "output path for vint.yaml")
	dryRun := fs.Bool("dry-run", false, "print changes without writing files")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	result, err := Migrate(*configPath, *dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Print warnings.
	for _, w := range result.Warnings {
		fmt.Fprintf(os.Stderr, "WARNING: %s\n", w)
	}

	if len(result.UnknownLinters) > 0 {
		fmt.Fprintf(os.Stderr, "Unknown linters (not migrated): %s\n",
			strings.Join(result.UnknownLinters, ", "))
	}

	if *dryRun {
		fmt.Println("=== vint.yaml ===")
		fmt.Print(result.VintYAML)
		for path, content := range result.ConvertedFiles {
			fmt.Printf("\n=== %s ===\n", path)
			fmt.Print(content)
		}
		return
	}

	// Write vint.yaml.
	if err := os.WriteFile(*outputConfig, []byte(result.VintYAML), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", *outputConfig, err)
		os.Exit(1)
	}
	fmt.Printf("Wrote %s\n", *outputConfig)

	// Write converted Go files.
	for path, content := range result.ConvertedFiles {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", path, err)
			os.Exit(1)
		}
		fmt.Printf("Converted %s\n", path)
	}
}

// Migrate performs the full migration: config conversion + nolint conversion.
func Migrate(configPath, dir string) (*MigrateResult, error) {
	// Load golangci-lint config.
	golangciCfg, err := LoadGolangciConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("load golangci config: %w", err)
	}

	// Load the rules registry.
	registry, err := LoadRegistry()
	if err != nil {
		return nil, fmt.Errorf("load registry: %w", err)
	}

	result := &MigrateResult{
		ConvertedFiles: make(map[string]string),
	}

	// Migrate config for each enabled linter.
	allConfigs := make(map[string]VintRuleConfig)
	enabledLinters := golangciCfg.EnabledLinters()

	// Track which linters contributed each rule, so exclusion presets can
	// correctly handle rules contributed by multiple linters.
	ruleSources := make(map[string]map[string]bool)
	addRuleSource := func(rulePath, linterName string) {
		if ruleSources[rulePath] == nil {
			ruleSources[rulePath] = make(map[string]bool)
		}
		ruleSources[rulePath][linterName] = true
	}

	for _, linterName := range enabledLinters {
		m, ok := migrators[linterName]
		if !ok {
			// Check if we at least know this linter from the registry.
			if len(registry.RulesForLinter(linterName)) > 0 {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("linter %q has mapped rules but no migrator — rules will be enabled with defaults", linterName))
				// Enable the mapped rules with default config, but do not
			// overwrite configs already set by a proper migrator.
				for _, mapped := range registry.RulesForLinter(linterName) {
					if mapped.FullVintPath != "" {
						if _, alreadySet := allConfigs[mapped.FullVintPath]; !alreadySet {
							allConfigs[mapped.FullVintPath] = VintRuleConfig{}
						}
						addRuleSource(mapped.FullVintPath, linterName)
					}
				}
			} else {
				result.UnknownLinters = append(result.UnknownLinters, linterName)
			}
			continue
		}

		// Get linter-specific settings from golangci config.
		settings := golangciCfg.LinterSettings(linterName)

		configs, err := m.MigrateConfig(settings)
		if err != nil {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("error migrating %q settings: %v", linterName, err))
			continue
		}

		// Collect warnings from migrators that support it.
		if wr, ok := m.(WarningReporter); ok {
			result.Warnings = append(result.Warnings, wr.Warnings()...)
		}

		for k, v := range configs {
			allConfigs[k] = v
			addRuleSource(k, linterName)
		}
	}

	// Apply exclusion presets to remove or modify rules.
	if presets := golangciCfg.Linters.Exclusions.Presets; len(presets) > 0 {
		enabledLintersSet := make(map[string]bool, len(enabledLinters))
		for _, name := range enabledLinters {
			enabledLintersSet[name] = true
		}
		presetWarnings := applyExclusionPresets(presets, enabledLintersSet, allConfigs, ruleSources)
		result.Warnings = append(result.Warnings, presetWarnings...)
	}

	// Render vint.yaml.
	yamlContent, err := RenderVintYAML(allConfigs)
	if err != nil {
		return nil, fmt.Errorf("render vint.yaml: %w", err)
	}
	result.VintYAML = yamlContent

	// Scan for .go files and convert nolint directives.
	goFiles, err := findGoFiles(dir)
	if err != nil {
		return nil, fmt.Errorf("scan go files: %w", err)
	}

	for _, goFile := range goFiles {
		content, err := os.ReadFile(goFile)
		if err != nil {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("cannot read %s: %v", goFile, err))
			continue
		}

		// Quick check if the file has any nolint directives.
		if !nolintRegexp.Match(content) {
			continue
		}

		converted, err := ConvertNolintInSource(goFile, string(content), registry, allConfigs)
		if err != nil {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("error converting %s: %v", goFile, err))
			continue
		}

		if converted.Directives > 0 {
			result.ConvertedFiles[goFile] = converted.NewContent
		}
	}

	return result, nil
}

func findGoFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			base := filepath.Base(path)
			if base == "vendor" || base == ".git" || base == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
