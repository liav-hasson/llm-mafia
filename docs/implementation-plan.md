# Mafia LLM Kubernetes Game – Incremental Implementation Plan

> [!NOTE]
> Temporary document! do not add to main documentation!

## Project Goal

Build a distributed Mafia game where:
- Each player is an LLM-backed Kubernetes pod
- Players communicate, reason, vote, and eliminate each other
- Eliminated players are deleted from the cluster
- Game state is observable and debuggable
- Architecture is cloud-agnostic and production-inspired

---

## DONE: Phase 0 – Architecture Lock-In (No Code Yet)

**Time:** ~2–3 hours

### Goals
- Avoid premature implementation mistakes
- Define strict boundaries between components

### Tasks
- Write a short architecture doc describing:
  - Game Engine responsibilities
  - Player Agent responsibilities
  - Kafka as event backbone
  - CRD as control plane only
- Decide naming conventions:
  - `game_id`
  - `player_id`
  - topic names
- Decide repo structure (mono vs multi-repo)

### Decisions Made
- **Language:** Go for Controller and Game Engine (industry standard for operators, tightly coupled components), Python for Player Agents (faster iteration, easier LLM integration)
- **Repo structure:** Monorepo (single developer, tightly coupled components, atomic changes)
- **Engine model:** One engine per game (isolation, Kubernetes-native cleanup via OwnerReferences, easier debugging)
- **CRD usage:** Justified for learning operators + portfolio value (declarative API, status tracking, automatic garbage collection)
- **Event schema:** JSON with strongly-typed models (human-readable, easier debugging, sufficient for project scale)

---

## DONE: Phase 1 – Local Dev Environment Bootstrap

**Time:** ~1 hour

### Goals
- Be able to iterate fast locally
- No cloud dependency
- Environment/cloud agnostic

### Tasks
- Create a local Kubernetes cluster using Kind
- Create bootstrap scripts for cluster up/down
- Make all configurations easily adjustable

### Decisions Made
- **Cluster:** Kind (single node for local dev)
- **Monitoring:** Deferred to Phase 11 (not needed initially)
- **Observability (minimal):** `kubectl top` for cluster health, Kafka consumer for game events
- **Cloud agnostic:** All manifests work on any Kubernetes cluster
- **Test configuration:** 5 players, small Ollama model (llama3.2:1b or 3b)
- **Future:** Experiment on cloud VM with larger hardware once game is functional

### Resource Estimates (Local Dev)
- CPU: 4-6 cores
- RAM: 8-16GB
- Ollama: Small model (1-3B parameters) for testing

---

## DONE: Phase 2 – Kafka Backbone & Dev Tooling

**Time:** ~3–4 hours

### Goals
- Establish the system's nervous system early
- Everything meaningful becomes an event
- Set up quality gates before writing code

### Tasks
- **Pre-commit setup:**
  - Python: ruff (lint + format), bandit (security), mypy (types)
  - Go: golangci-lint, gosec
  - Pre-commit hooks configuration
- **GitHub Actions CI:**
  - Lint and security checks
  - Unit tests
  - Docker image builds
- Deploy Kafka using Strimzi Operator
- Create initial topics:
  - `game.events` - public game events (deaths, phase changes)
  - `game.chat.<game_id>` - public discussion
  - `game.mafia.<game_id>` - private mafia channel
  - `game.votes.<game_id>` - vote submissions
- Write a minimal producer + consumer (local binary or pod)
- Verify:
  - ordering
  - consumer groups
  - offsets

### Decisions Made
- **Topic strategy:** Per-game topics with game_id in topic name (isolation, easier cleanup)
- **Event format:** JSON with Pydantic models (Python) / structs with tags (Go)
- **Python environment:** Python 3.12.3, venv at `src/.venv`

---

## Phase 3 – Game Engine Skeleton (No LLM Yet)

**Time:** ~2–3 hours

### Goals
- Central authority that enforces rules
- No intelligence yet, only mechanics

### Tasks
- Create `game-engine` service in **Go**:
  - Subscribes to Kafka
  - Maintains in-memory game state
- Implement:
  - Game creation
  - Player registration
  - Day/Night phase toggling
- Emit system events back to Kafka
- Unit tests for game logic

### Open Decisions
- Stateless engine + replay on restart vs in-memory only
- Single-threaded loop vs async workers
- Potentially fork existing mafia engines for skeleton for game engine and frontend design

---

## Phase 4 – Player Agent Skeleton (Dumb Players)

**Time:** ~2–3 hours

### Goals
- Validate communication model
- No LLMs yet, randomize/hardcoded actions

### Tasks
- Create `player-agent` service:
  - Subscribes to assigned topics
  - Publishes structured responses
- Implement:
  - Reactive message handling (not turn-based)
  - Vote submission
- Run 5 agents locally and simulate a game

### Decisions Made
- **Player count for testing:** 5 players (configurable)
- **Communication:** All via Kafka (no direct HTTP/gRPC between agent and engine)

---

## Phase 5 – Hello World Operator (Learning Phase)

**Time:** ~2–3 hours

### Goals
- Learn Kubernetes operator patterns in isolation
- Build confidence with CRDs and controllers before game logic

### Tasks
- Create a trivial `Greeting` CRD:
  - Spec: `name`, `message`
  - Status: `configMapName`, `lastUpdated`
- Implement a Go controller using Kubebuilder that:
  - Watches `Greeting` resources
  - Creates a ConfigMap with the greeting message
  - Updates status on the CR
- Test the full reconciliation loop
- Practice `kubectl` commands for CRs

### Decisions Made
- **Language:** Go with Kubebuilder (industry standard)

---

## Phase 6 – CRD Design (Control Plane Only)

**Time:** ~2–3 hours

### Goals
- Formalize intent-driven game creation
- Avoid stuffing logic into the CRD

### Tasks
- Define `MafiaGame` CRD:
  - Spec:
    - players.count
    - rules
    - modelProfile
  - Status:
    - phase
    - alivePlayers
    - eliminatedPlayers
    - winner
- Apply CRD to cluster
- Validate with `kubectl`

### Open Decisions
- CRD versioning strategy

---

## Phase 7 – Controller (Lifecycle Automation)

**Time:** ~3–4 hours

### Goals
- Kubernetes-native lifecycle management
- Pods appear and disappear based on game state

### Tasks
- Implement controller reconciliation loop:
  - On new `MafiaGame` → spawn engine + players
  - On elimination event → delete player pod
  - On game end → cleanup
- Update CRD status safely

### Open Decisions
- One controller managing all games vs controller-per-game
- OwnerReferences vs explicit cleanup logic

---

## Phase 8 – Mock Agent Mode (Fast Iteration)

**Time:** ~1–2 hours

### Goals
- Enable rapid testing without LLM overhead
- Speed up development iteration dramatically

### Tasks
- Implement mock/rule-based agent mode:
  - Random voting with weighted probabilities
  - Simple keyword-based responses
  - Configurable via environment variable (`AGENT_MODE=mock|llm`)
- Add test scenarios that run with mock agents
- Ensure game can complete end-to-end with mocks

### Decisions Made
- Mock mode is the default for local development

---

## Phase 9 – Introduce Ollama (Centralized)

**Time:** ~2–3 hours

### Goals
- Replace dumb agents with real LLM-backed reasoning
- Keep resource usage predictable

### Tasks
- Deploy Ollama as a centralized service
- Configure:
  - model download (configurable via environment)
  - concurrency limits
- Modify agents:
  - Build prompts
  - Send reasoning requests to Ollama
- Capture prompts + responses as events/logs
- Ensure graceful fallback to mock mode if Ollama unavailable

### Decisions Made
- **Local testing:** Small model (llama3.2:1b or 3b)
- **Cloud/production:** Configurable, experiment with larger models
- **Model selection:** Environment variable, easy to swap

---

## Phase 10 – Private Reasoning & Mafia Channels

**Time:** ~2–3 hours

### Goals
- Enable deception and hidden coordination
- Preserve observability

### Tasks
- Add Kafka topics:
  - `game.thoughts.<player_id>`
  - `game.mafia`
- Enforce access rules in engine:
  - Who can publish where
- Log all messages but tag visibility

### Open Decisions
- Kafka ACLs vs engine-enforced rules
- Encrypt private topics or not
- Different player traits - aggressive, talkative, reserved, etc...

---

## Phase 11 – CI/CD & Validation

**Time:** ~2–3 hours

### Goals
- Reproducible builds
- Confidence in changes

### Tasks
- GitHub Actions:
  - build images
  - run unit tests
  - spin kind cluster
  - deploy a 3-player game (using mock agents for speed)
- Fail pipeline if game does not terminate cleanly

### Open Decisions
- Integration tests vs simulation-only tests
- Release strategy (tags vs main-only)

---

## Phase 12 – Bootstrap Guide & Quick Start

**Time:** ~2–3 hours

### Goals
- Enable anyone to run the project with minimal friction
- Provide clear onboarding for contributors

### Tasks
- Create `README.md` with:
  - Project overview and architecture diagram
  - Prerequisites (Docker, kind/k3d, kubectl, Go, Python)
  - Quick start commands (copy-paste friendly)
- Create `BOOTSTRAP.md` with detailed steps:
  - Clone repo
  - Set up local cluster
  - Deploy infrastructure (Kafka, monitoring)
  - Run a sample game
  - Verify game completion
- Create `Makefile` or `Taskfile.yaml` with common commands:
  - `make cluster-up` / `make cluster-down`
  - `make deploy-infra`
  - `make deploy-game`
  - `make run-demo`
  - `make logs`
  - `make clean`
- Add troubleshooting section for common issues

### Deliverables
- `README.md` – project overview
- `BOOTSTRAP.md` – step-by-step setup guide
- `Makefile` or `Taskfile.yaml` – automation commands

---

## Phase 13 – Testing & Experimenting

### Goals
- Test the optimal configurations for the best game and AI quality

### Tasks
- Run the app on different providors
- RUn the app on different settings / configurations

### Deliverables
- Create document of the experiments findings

---

## Phase 14 – Documentation & Presentation

**Time:** ~3–4 hours

### Goals
- Create portfolio-ready documentation
- Prepare for interviews and demos

### Tasks
- Write `docs/ARCHITECTURE.md`:
  - System components and responsibilities
  - Data flow diagrams
  - Technology choices and rationale
- Write `docs/DECISIONS.md` (ADR – Architecture Decision Records):
  - Why Kafka over Redis Streams
  - Why Go for controller, Python for agents
  - Why CRD-based game management
- Create presentation materials:
  - 5-10 slide deck (PDF or Google Slides)
  - Architecture diagram (draw.io or Excalidraw)
  - Demo script with talking points
- Record a short demo video (optional but impressive):
  - Show `kubectl apply` starting a game
  - Show pods appearing/disappearing
  - Show Grafana dashboard
  - Show game transcript
- Update GitHub repo:
  - Add badges (build status, license)
  - Add screenshots/GIFs to README
  - Add "What I Learned" section

### Deliverables
- `docs/ARCHITECTURE.md`
- `docs/DECISIONS.md`
- Presentation slide deck
- (Optional) Demo video

---

## Phase 15 – Future Extensions (Optional)

- implement gitops
- game persistense in a db
- add prometheus + grafana + loki
---
