package main

import (
	"os/signal"
	"syscall"
	"context"
	"os/signal"
	"syscall"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

type ConfigDiff struct {
	Common   map[string]interface{}
	OnlyLeft map[string]interface{}
	OnlyRight map[string]interface{}
	Changed map[string]map[string]string
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println(color.CyanString("configdiff - Configuration File Diff Tool"))
		fmt.Println()
		fmt.Println("Usage: configdiff <config1> <config2>")
		fmt.Println()
		fmt.Println("Supported formats: YAML, JSON, TOML")
		os.Exit(1)
	}

	config1 := os.Args[1]
	config2 := os.Args[2]

	cfg1, err := loadConfig(config1)
	if err != nil {
		color.Red("Error loading %s: %v", config1, err)
		os.Exit(1)
	}

	cfg2, err := loadConfig(config2)
	if err != nil {
		color.Red("Error loading %s: %v", config2, err)
		os.Exit(1)
	}

	diff := compareConfigs(cfg1, cfg2)
	displayDiff(diff)
}

func loadConfig(filename string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg map[string]interface{}
	
	// Try to detect format by extension
	if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		err = yaml.Unmarshal(data, &cfg)
	} else if strings.HasSuffix(filename, ".json") {
		err = json.Unmarshal(data, &cfg)
	} else {
		// Try YAML first, then JSON
		err = yaml.Unmarshal(data, &cfg)
		if err != nil {
			err = json.Unmarshal(data, &cfg)
		}
	}

	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func compareConfigs(cfg1, cfg2 map[string]interface{}) ConfigDiff {
	diff := ConfigDiff{
		Common:    make(map[string]interface{}),
		OnlyLeft:  make(map[string]interface{}),
		OnlyRight: make(map[string]interface{}),
		Changed:   make(map[string]map[string]string),
	}

	leftKeys := make(map[string]bool)
	rightKeys := make(map[string]bool)

	for k := range cfg1 {
		leftKeys[k] = true
	}

	for k := range cfg2 {
		rightKeys[k] = true
	}

	for k := range leftKeys {
		if rightKeys[k] {
			v1 := cfg1[k]
			v2 := cfg2[k]
			
			if fmt.Sprintf("%v", v1) == fmt.Sprintf("%v", v2) {
				diff.Common[k] = v1
			} else {
				diff.Changed[k] = map[string]string{
					"old": fmt.Sprintf("%v", v1),
					"new": fmt.Sprintf("%v", v2),
				}
			}
		} else {
			diff.OnlyLeft[k] = cfg1[k]
		}
	}

	for k := range rightKeys {
		if !leftKeys[k] {
			diff.OnlyRight[k] = cfg2[k]
		}
	}

	return diff
}

func displayDiff(diff ConfigDiff) {
	fmt.Println(color.CyanString("\n=== CONFIG DIFF REPORT ===\n"))

	if len(diff.Common) > 0 {
		fmt.Printf(color.GreenString("Common keys (%d):\n"), len(diff.Common))
		fmt.Println(strings.Repeat("=", 50))
		for k := range diff.Common {
			fmt.Printf("  %s\n", k)
		}
		fmt.Println()
	}

	if len(diff.OnlyLeft) > 0 {
		fmt.Printf(color.HiYellowString("Only in config1 (%d):\n"), len(diff.OnlyLeft))
		fmt.Println(strings.Repeat("-", 50))
		for k := range diff.OnlyLeft {
			fmt.Printf("  - %s\n", k)
		}
		fmt.Println()
	}

	if len(diff.OnlyRight) > 0 {
		fmt.Printf(color.HiYellowString("Only in config2 (%d):\n"), len(diff.OnlyRight))
		fmt.Println(strings.Repeat("-", 50))
		for k := range diff.OnlyRight {
			fmt.Printf("  + %s\n", k)
		}
		fmt.Println()
	}

	if len(diff.Changed) > 0 {
		fmt.Printf(color.HiRedString("Changed values (%d):\n"), len(diff.Changed))
		fmt.Println(strings.Repeat("-", 50))
		for k, v := range diff.Changed {
			fmt.Printf("  %s\n", k)
			fmt.Printf("    - old: %s\n", color.RedString(v["old"]))
			fmt.Printf("    + new: %s\n", color.GreenString(v["new"]))
		}
		fmt.Println()
	}

	fmt.Println(color.YellowString("\n=== MIGRATION SCRIPT ==="))
	printMigrationScript(diff)
}

func printMigrationScript(diff ConfigDiff) {
	fmt.Println("\n# Migration from config1 to config2:")
	
	keys := sortedKeys(diff.OnlyLeft)
	for _, k := range keys {
		fmt.Printf("# Remove: %s\n", k)
	}
	
	keys = sortedKeys(diff.OnlyRight)
	for _, k := range keys {
		fmt.Printf("# Add: %s\n", k)
	}

	keys = sortedKeysMap(diff.Changed)
	for _, k := range keys {
		fmt.Printf("# Update: %s\n", k)
	}
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedKeysMap(m map[string]map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}