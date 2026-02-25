# configdiff

Configuration File Diff Tool - compares YAML, JSON, and TOML configuration files.

## Purpose

Identify differences between two configuration files including added, removed, and modified keys.

## Installation

```bash
go build -o configdiff ./cmd/configdiff
```

## Usage

```bash
configdiff <config1> <config2>
```

### Supported Formats

- YAML (.yaml, .yml)
- JSON (.json)
- TOML (auto-detected)

### Examples

```bash
# Compare two YAML files
configdiff config.dev.yaml config.prod.yaml

# Compare JSON configurations
configdiff settings.json settings.backup.json

# Compare mixed formats
configdiff config.yaml settings.json
```

## Output

```
=== CONFIG DIFF REPORT ===

Common keys (5):
==================================================
  database_url
  port
  log_level
  timeout
  max_connections

Only in config1 (2):
--------------------------------------------------
  - debug_mode
  - legacy_setting

Only in config2 (1):
--------------------------------------------------
  + new_feature

Changed values (1):
--------------------------------------------------
  database_url
    - old: postgres://localhost:5432
    + new: postgres://prod-db:5432

=== MIGRATION SCRIPT ===

# Migration from config1 to config2:
# Remove: debug_mode
# Remove: legacy_setting
# Add: new_feature
# Update: database_url
```

## Dependencies

- Go 1.21+
- github.com/fatih/color
- gopkg.in/yaml.v3

## Build and Run

```bash
# Build
go build -o configdiff ./cmd/configdiff

# Run
go run ./cmd/configdiff config1.yaml config2.yaml
```

## License

MIT