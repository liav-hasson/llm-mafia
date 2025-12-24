#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../../helpers/check-commands.sh"

echo "[START] Bootstrapping Kafka..."

check_commands kubectl

echo "[INFO] Creating kafka namespace..."
kubectl create ns kafka

echo "[INFO] Installing Strimzi operator..."
kubectl apply -f "https://strimzi.io/install/latest?namespace=kafka" -n kafka

echo "[INFO] Waiting for Strimzi operator to be ready..."
kubectl wait deployment/strimzi-cluster-operator \
    -n kafka \
    --for=condition=Available \
    --timeout=300s

echo "[INFO] Deploying Kafka cluster..."
kubectl apply -f "$SCRIPT_DIR/kafka-cluster.yaml"

echo "[INFO] Waiting for Kafka cluster to be ready..."
kubectl wait kafka/mafia \
    -n kafka \
    --for=condition=Ready \
    --timeout=600s

echo "[INFO] Creating Kafka topics..."
kubectl apply -f "$SCRIPT_DIR/topics/"

echo "[SUCCESS] Kafka is ready!"
echo ""
echo "Broker address: mafia-kafka-bootstrap.kafka.svc:9092"
echo ""
echo "Verify with:"
echo "  kubectl get kafka -n kafka"
echo "  kubectl get kafkatopics -n kafka"
