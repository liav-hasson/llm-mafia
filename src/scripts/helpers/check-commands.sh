#!/bin/bash

set -e

# Helper function to verify commands exists
check_commands() {

    local script_name="$(basename "${BASH_SOURCE[1]}")"
    local script_parent_dir=$(basename "$(dirname "${BASH_SOURCE[1]}")")
    
    echo "[INFO] - Checking commands in $script_parent_dir/$script_name."

    # Count missing commands
    local FAILED=0

    # Array of all commands parsed to strings
    for command in "$@"; do
        if ! command -v "$command" >/dev/null 2>&1; then
            echo "[ERROR] - Command '$command' not found."
            ((FAILED++))
        fi 
    done

    if [[ $FAILED -gt 0 ]]; then
        return 1
    fi

    echo "[SUCCESS] - All commands found."
}