# Game Engine

> [!NOTE]
> Draft document!

Central authority that enforces game rules and manages state. Written in Go.

---

## Project Structure

```
src/engine/
├── cmd/engine/main.go           # Entry point
├── internal/
│   ├── config/                  # Configuration
│   ├── domain/                  # Pure game logic
│   ├── engine/                  # Orchestration
│   ├── events/                  # Event types
│   └── kafka/                   # Kafka adapter
├── go.mod
├── Dockerfile
└── README.md
```

### Package Responsibilities

| Package | Pure | Description |
|---------|------|-------------|
| `cmd/engine` | - | Entry point, wires dependencies together |
| `config` | Yes | Load configuration from environment variables |
| `domain` | Yes | Game rules, state, validations (no I/O) |
| `engine` | No | Event loop, handlers, timers (has side effects) |
| `events` | Yes | Event struct definitions, JSON serialization |
| `kafka` | No | Kafka consumer/producer wrappers |

### Why This Separation?

**domain/ is pure:**
- No Kafka imports, no I/O, no external dependencies
- Can be unit tested without mocking anything
- Game rules are isolated and easy to verify

**engine/ orchestrates:**
- Consumes from Kafka, updates domain state, produces events
- Contains the main event loop and timer logic
- Depends on domain, events, and kafka packages

**kafka/ is an adapter:**
- Wraps segmentio/kafka-go library
- If we switch message brokers, only this package changes
- Translates between Kafka messages and internal event types

---

