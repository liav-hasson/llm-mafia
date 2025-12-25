# Kafka

> [!NOTE]
> Draft document!

Event backbone for game communication. All game state changes and player interactions flow through Kafka.

---

## Architecture

> [!NOTE]
> TODO: Change to actual diagram
```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Engine    │────▶│    Kafka    │◀────│   Players   │
│  (Go)       │◀────│   Broker    │────▶│  (Python)   │
└─────────────┘     └─────────────┘     └─────────────┘
```

- Engine publishes game events, subscribes to votes
- Players publish messages/votes, subscribe to events
- No direct communication between components

---

## Topics

| Topic | Direction | Purpose |
|-------|-----------|---------|
| `game.events` | Engine → Players | Phase changes, deaths, game state |
| `game.chat` | Players → All | Public discussion |
| `game.mafia` | Mafia → Mafia | Private night coordination |
| `game.votes` | Players → Engine | Vote submissions |
| `game.thoughts` | Players → Logging | LLM reasoning (observability only) |
| `game.doctor` | Doctor → Engine | Night protection target |
| `game.sheriff` | Sheriff → Engine | Night investigation target |

All topics use single partition for strict message ordering.

### Topic Isolation

Special roles have dedicated topics to prevent information leakage:
- Mafia cannot see who doctor protects
- Mafia cannot see who sheriff investigates
- Engine resolves all night actions privately

---

## Operator

Uses [Strimzi](https://strimzi.io/) operator for Kafka management.

| Resource | Purpose |
|----------|---------|
| `Kafka` | Cluster definition (brokers, listeners, storage) |
| `KafkaNodePool` | Broker and controller node configuration |
| `KafkaTopic` | Topic definitions as Kubernetes resources |

Benefits:
- Declarative configuration via CRDs
- Kubernetes-native lifecycle management
- Automatic topic creation from YAML

---

## Cluster Configuration

- **Mode:** KRaft (Zookeeper deprecated)
- **Brokers:** 1 (TODO: scale to 3)
- **Storage:** Persistent on node (survives pod restarts)
- **Listener:** Plain on port 9092 (TODO: implement TLS)

Broker address from within cluster:
```
mafia-kafka-kafka-bootstrap.kafka.svc:9092
```

---

## Files

- [kafka-cluster.yaml](../src/scripts/bootstrap/kafka/kafka-cluster.yaml) - Cluster and node pool definitions
- [topics/game-topics.yaml](../src/scripts/bootstrap/kafka/topics/game-topics.yaml) - Topic definitions
- [install.sh](../src/scripts/bootstrap/kafka/install.sh) - Deployment script

---

## Design Decisions

**Why Kafka over Redis Streams or direct HTTP?**
- Decoupled architecture (components don't know about each other)
- Persistent message log (replay possible)
- Industry-standard for event-driven systems

**Why separate topics vs single topic?**
- Access control (mafia topic is private)
- Simpler consumer logic
- Clear semantic boundaries

**Why 1 partition?**
- Guarantees message ordering within a game
- Single game at a time (no parallelism needed)
