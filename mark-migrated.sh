#!/usr/bin/env bash
#
# Mark a linter as having its migrator coded in the linters registry.
#
# Usage:
#   ./mark-migrated.sh <linter_name>
#
# Examples:
#   ./mark-migrated.sh bodyclose
#   ./mark-migrated.sh errcheck

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REGISTRY_DIR="$SCRIPT_DIR/../linters_registry"

cd "$REGISTRY_DIR"
go run . migrated "$@"
