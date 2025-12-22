# Mafia LLM Kubernetes Game – Incremental Implementation Plan

> [!NOTE]  
> Temporary document!

## Project Goal

Build a distributed Mafia game where:
- Each player is an LLM-backed Kubernetes pod
- Players communicate, reason, vote, and eliminate each other
- Eliminated players are deleted from the cluster
- Game state is observable and debuggable
- Architecture is cloud-agnostic and production-inspired

## Open Questions (To Resolve During Implementation)

- [ ] Kafka topic partitioning strategy
- [ ] Event schema format (JSON vs Protobuf)
- [ ] Agent timeout handling
- [ ] Graceful shutdown behavior
- [ ] CRD versioning strategy
- [ ] 
---

## - DONE: Phase 0 – Architecture Lock-In (No Code Yet) 

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
- **Language:** Python for Player Agents (faster iteration, easier LLM integration), Go for Controller (industry standard for operators)
- **Repo structure:** Monorepo (single developer, tightly coupled components, atomic changes)
- **Engine model:** One engine per game (isolation, Kubernetes-native cleanup via OwnerReferences, easier debugging)
- **CRD usage:** Justified for learning operators + portfolio value (declarative API, status tracking, automatic garbage collection)



---

## Phase 1 – Local Dev Environment Bootstrap

**Time:** ~2–3 hours

### Goals
- Be able to iterate fast locally
- No cloud dependency

### Tasks
- Create a local Kubernetes cluster (`kind` or `k3d`)
- Install:
  - Prometheus
  - Grafana
  - Loki
- Create a Makefile or task runner for:
  - cluster up/down
  - deploy all core components

### Open Decisions
- `kind` vs `k3d`
- Helm vs raw manifests for infra

---

## Phase 2 – Kafka Backbone (Event Plumbing)

**Time:** ~3–4 hours

### Goals
- Establish the system’s nervous system early
- Everything meaningful becomes an event

### Tasks
- Deploy Kafka using Strimzi Operator
- Create initial topics:
  - `game.events`
  - `game.system`
  - `game.votes`
- Write a minimal producer + consumer (local binary or pod)
- Verify:
  - ordering
  - consumer groups
  - offsets

### Open Decisions
- One topic per game vs shared topic with `game_id` key
- JSON vs Protobuf event schemas

---

## Phase 3 – Game Engine Skeleton (No LLM Yet)

**Time:** ~2–3 hours

### Goals
- Central authority that enforces rules
- No intelligence yet, only mechanics

### Tasks
- Create `game-engine` service:
  - Subscribes to Kafka
  - Maintains in-memory game state
- Implement:
  - Game creation
  - Player registration
  - Day/Night phase toggling
- Emit system events back to Kafka

### Open Decisions
- Stateless engine + replay on restart vs in-memory only
- Single-threaded loop vs async workers
- Potentially fork existing mafia engines for skeleton

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
  - Turn-based message handling
  - Vote submission
- Run 3–5 agents locally and simulate a game

### Open Decisions
- HTTP vs gRPC between agent and engine
- Push vs pull turn signaling

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
  - model download
  - concurrency limits
- Modify agents:
  - Build prompts
  - Send reasoning requests to Ollama
- Capture prompts + responses as events/logs
- Ensure graceful fallback to mock mode if Ollama unavailable

### Open Decisions
- Small vs medium model

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

## Phase 11 – Observability Deep Dive

**Time:** ~2–3 hours

### Goals
- Make the system explainable
- Identify bottlenecks

### Tasks
- Add structured logging (JSON)
- Export Prometheus metrics:
  - LLM latency
  - Kafka lag
  - Phase duration
- Create Grafana dashboards:
  - Timeline view per game
  - Per-player behavior
  - Engine health

### Open Decisions
- Trace reasoning chains (Tempo) or logs only
- Sampling vs full capture

---

## Phase 12 – Game State Persistence (Optional)

**Time:** ~2–3 hours

### Goals
- Replay games
- Enable analytics and UI later

### Tasks (choose one)
- Option A: Kafka-only event sourcing
- Option B: Snapshot state to Postgres
- Option C: Hybrid (Kafka + DB projections)

### Open Decisions
- DB choice (Postgres, SQLite, none)
- Snapshot frequency

---

## Phase 13 – CI/CD & Validation

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

## Phase 14 – GitOps Deployment

**Time:** ~2–3 hours

### Goals
- Declarative, Git-driven deployments
- Practice production-grade deployment patterns

### Tasks
- Set up ArgoCD or Flux in the cluster
- Create GitOps structure:
  - `gitops/base/` – base manifests
  - `gitops/overlays/local/` – local dev overrides
  - `gitops/overlays/prod/` – production config (future)
- Configure ApplicationSet or Kustomization for:
  - Infrastructure components (Kafka, Ollama, monitoring)
  - Game controller
- Enable auto-sync for local environment
- Document GitOps workflow in README

### Decisions Made
- ArgoCD recommended (more visual, good learning experience)

---

## Phase 15 – Bootstrap Guide & Quick Start

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

## Phase 16 – Documentation & Presentation

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

## Phase 17 – Future Extensions (Explicitly Out of Scope)

- Web UI for live game visualization
- Multi-game tournaments
- Role variety (detective, doctor)
- Model comparison experiments
- Multi-cluster games

---

## Guiding Principles

- CRD = intent, not data dump
- Kafka = source of truth
- Engine = authority
- Agents = untrusted participants
- Observability first, polish later

---

## Success Criteria (v1)

- One `kubectl apply` starts a game
- Players reason and vote
- Pods are deleted when eliminated
- Full transcript observable
- Game ends deterministically

---

