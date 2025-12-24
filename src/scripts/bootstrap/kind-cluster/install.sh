#!/bin/bash

set -euo pipefail

# Source command checker helper
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../../helpers/check-commands.sh"

echo "[START] Bootstrapping Kind cluster..."

check_commands kind

# Create the cluster from the config file
echo "[INFO] Creating Kind cluster..."
kind create cluster \
    --name "mafia" \
    --config "$SCRIPT_DIR/kind-config.yaml"

echo "[SUCCESS] Kind cluster deployed!"
