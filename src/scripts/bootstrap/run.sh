#!/bin/bash

set -euo pipefail

# Safe to run from anywhere
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"


# Bootstapping the project
./kind-cluster/install.sh

echo
echo "---"
echo

./kafka/install.sh
