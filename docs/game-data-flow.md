# Game Logic & LLM Interaction

> [!NOTE]
> Draft document!

## Game Flow

> [!NOTE]
> TODO: Change to actual diagram
```
GAME START
    │
    ├── Engine assigns roles (Villager/Mafia)
    ├── Engine creates Kafka topics
    └── Players initialize with base context
           │
           ▼
    ┌──────────────────────────────────────┐
    │            ROUND LOOP                │
    │                                      │
    │   NIGHT ─────────────────────────    │
    │      │                               │
    │      └── Mafia coordinates           │
    │      └── Mafia submits kill vote     │
    │      └── Engine announces death      │
    │             │                        │
    │             ▼                        │
    │   DAY ───────────────────────────    │
    │      │                               │
    │      └── Players discuss freely      │
    │      └── Players speak when ready    │
    │             │                        │
    │             ▼                        │
    │   VOTING ────────────────────────    │
    │      │                               │
    │      └── Each player votes           │
    │      └── Engine tallies & eliminates │
    │             │                        │
    │             ▼                        │
    │   ROUND END ─────────────────────    │
    │      │                               │
    │      └── Each player summarizes      │
    │      └── Context resets for round    │
    │                                      │
    └──────────────────────────────────────┘
           │
           ▼
    GAME END (Mafia or Villagers win)
```

---

## Player Context (What the LLM Sees)

Each player maintains three layers of context:

| Layer | Content | Lifecycle |
|-------|---------|-----------|
| **Base** | Rules, role, player list, personality | Set once, never changes |
| **Summaries** | Player's summary of each past round | Grows by ~300 tokens/round |
| **Current Round** | All messages this round | Cleared after each round |

**Full context sent to LLM:**
```
[Base Context]
[All Round Summaries]
[Current Round Messages]
[Query: What do you say? / Who do you vote for?]
```

---

## LLM Queries

### During Discussion
- **Trigger:** New messages arrive + cooldown expired
- **Context:** Base + Summaries + Current Round
- **Prompt:** "What do you want to say? Reply [SILENT] if nothing."
- **Cooldown:** Wait 3+ seconds after speaking before speaking again

### During Voting
- **Trigger:** Voting phase starts
- **Context:** Base + Summaries + Current Round
- **Prompt:** "Who do you vote for and why?"

### At Round End
- **Trigger:** Round ends
- **Context:** Base + Current Round only
- **Prompt:** "Summarize this round. Who is suspicious? What's your strategy?"
- **Output:** Saved as this round's summary, current round cleared

---

## Reactive Speaking (Not Turn-Based)

Players decide when to speak. No one forces them.

**Triggers to speak:**
- 2-3 new messages since I last spoke
- Someone accused me
- Random small chance (organic flow)

**Collision prevention:**
- Random delay (1-5 seconds) before responding
- Cooldown after speaking
- Check if someone else already responded

---

## Mafia Coordination

Mafia players have a private channel (`game.<id>.mafia`).

During night:
1. Mafia see each other's identity
2. Mafia discuss on private channel
3. Mafia submit kill vote
4. Engine processes majority vote

---

## Win Conditions

- **Villagers win:** All Mafia eliminated
- **Mafia wins:** Mafia count ≥ Villager count
