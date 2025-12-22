# LLM-Mafia

> [!NOTE]  
> Draft document!

Kubernetes pods playing Mafia using AI reasoning.

---

## What is this?

A distributed Mafia game where each player is an LLM-backed Kubernetes pod. Players communicate via Kafka, reason with Ollama, and eliminated players get deleted from the cluster. The game is managed declaratively using a custom CRD.

---

## Components

| Component | Description | Docs |
|-----------|-------------|------|
| **Controller** | Kubernetes operator that manages game lifecycle | `cmd/controller/README.md` |
| **Engine** | Game authority that enforces rules and phases | `engine/README.md` |
| **Agent** | LLM-backed player that reasons and votes | `agent/README.md` |
| **CRD** | Declarative game definition | `manifests/crds/README.md` |

---

## Quick Start

```bash

```

---

## Documentation

- [Game Logic & LLM Interaction](docs/ARCHITECTURE.md)

---

## Project Structure

```
mafia/
├── cmd/controller/     # Go operator
├── engine/             # Python game engine
├── agent/              # Python player agent
├── manifests/          # CRDs and sample games
├── gitops/             # ArgoCD applications
├── scripts/            # Helper scripts
├── docs/               # Documentation
└── Makefile
```