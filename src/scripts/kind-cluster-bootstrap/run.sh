#!/bin/bash

set -euo pipefail

# Create the cluster from the config file
kind create cluster \
    --name mafia \
    --config kind-config.yaml
