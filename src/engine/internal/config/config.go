package config

// this package provides a centralized loader for runtime
// configuration used by the engine. It reads values from environment variables,
// applies sensible defaults, and validates the result. The file intentionally
// uses only the standard library so it's easy to test and understand.
// kubernetes controller values can override these values.

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all runtime configuration for the engine.
// Environment variable names are suggested in comments on each field.
type Config struct {
	// Kafka connection settings - list of broker connections
	// KAFKA_BROKERS="broker1:9092,broker2:9092"
	KafkaBrokers []string

	// KAFKA_CLIENT_ID, KAFKA_GROUP_ID (optional, useful for debugging)
	KafkaClientID string
	KafkaGroupID  string

	// Coarse topic names (no role/chat/vote prefixes)
	// ENGINE_EVENTS_TOPIC (default: engine.events)
	EngineEventsTopic string
	// PLAYER_ACTIONS_TOPIC (default: player.actions)
	PlayerActionsTopic string

	// Timeouts and durations
	// KAFKA_CONSUMER_TIMEOUT (e.g. "2s")
	KafkaConsumerTimeout time.Duration
	// KAFKA_PRODUCER_TIMEOUT (e.g. "2s")
	KafkaProducerTimeout time.Duration
	// HTTP_TIMEOUT for external calls (e.g. Ollama) (optional)
	HTTPTimeout time.Duration

	// Game settings (defaults mirror previous constants)
	// GAME_MIN_PLAYERS (default 6)
	GameMinPlayers int
	// GAME_MAX_PLAYERS (default 12)
	GameMaxPlayers int

	// Agent mode: mock or llm (AGENT_MODE, default: mock)
	AgentMode string

	// Logging / environment
	// LOG_LEVEL (optional; default: info)
	LogLevel string
	// ENV (optional; dev/prod)
	Env string

	// Feature flags (optional)
	// ENABLE_ROLE_SECRETS (optional: "true"/"false")
	EnableRoleSecrets bool
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		KafkaBrokers:         []string{"localhost:9092"},
		KafkaClientID:        "mafia-engine",
		KafkaGroupID:         "mafia-engine-group",
		EngineEventsTopic:    "engine.events",
		PlayerActionsTopic:   "player.actions",
		KafkaConsumerTimeout: 2 * time.Second,
		KafkaProducerTimeout: 2 * time.Second,
		HTTPTimeout:          5 * time.Second,
		GameMinPlayers:       6,
		GameMaxPlayers:       12,
		AgentMode:            "mock",
		LogLevel:             "info",
		Env:                  "dev",
		EnableRoleSecrets:    false,
	}
}

// LoadConfig reads environment variables, applies defaults and returns a Config.
// It returns an error for invalid values or missing required settings.

func LoadConfig() (*Config, error) {
	cfg := DefaultConfig()

	if v, ok := lookupEnvTrim("KAFKA_BROKERS"); ok && v != "" {
		cfg.KafkaBrokers = parseCommaList(v)
	}

	if v, ok := lookupEnvTrim("KAFKA_CLIENT_ID"); ok && v != "" {
		cfg.KafkaClientID = v
	}
	if v, ok := lookupEnvTrim("KAFKA_GROUP_ID"); ok && v != "" {
		cfg.KafkaGroupID = v
	}

	// Coarse topic names
	if v, ok := lookupEnvTrim("ENGINE_EVENTS_TOPIC"); ok && v != "" {
		cfg.EngineEventsTopic = v
	}
	if v, ok := lookupEnvTrim("PLAYER_ACTIONS_TOPIC"); ok && v != "" {
		cfg.PlayerActionsTopic = v
	}

	// durations
	if v, ok := lookupEnvTrim("KAFKA_CONSUMER_TIMEOUT"); ok && v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid KAFKA_CONSUMER_TIMEOUT: %w", err)
		}
		cfg.KafkaConsumerTimeout = d
	}
	if v, ok := lookupEnvTrim("KAFKA_PRODUCER_TIMEOUT"); ok && v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid KAFKA_PRODUCER_TIMEOUT: %w", err)
		}
		cfg.KafkaProducerTimeout = d
	}
	if v, ok := lookupEnvTrim("HTTP_TIMEOUT"); ok && v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid HTTP_TIMEOUT: %w", err)
		}
		cfg.HTTPTimeout = d
	}

	// integers
	if v, ok := lookupEnvTrim("GAME_MIN_PLAYERS"); ok && v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid GAME_MIN_PLAYERS: %w", err)
		}
		cfg.GameMinPlayers = n
	}
	if v, ok := lookupEnvTrim("GAME_MAX_PLAYERS"); ok && v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid GAME_MAX_PLAYERS: %w", err)
		}
		cfg.GameMaxPlayers = n
	}

	if v, ok := lookupEnvTrim("AGENT_MODE"); ok && v != "" {
		cfg.AgentMode = v
	}

	if v, ok := lookupEnvTrim("LOG_LEVEL"); ok && v != "" {
		cfg.LogLevel = v
	}
	if v, ok := lookupEnvTrim("ENV"); ok && v != "" {
		cfg.Env = v
	}

	if v, ok := lookupEnvTrim("ENABLE_ROLE_SECRETS"); ok && v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("invalid ENABLE_ROLE_SECRETS: %w", err)
		}
		cfg.EnableRoleSecrets = b
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks config sanity and returns an error for invalid settings.
func (c *Config) Validate() error {
	if len(c.KafkaBrokers) == 0 {
		return errors.New("no kafka brokers configured (KAFKA_BROKERS)")
	}
	if c.KafkaConsumerTimeout <= 0 {
		return errors.New("KAFKA_CONSUMER_TIMEOUT must be > 0")
	}
	if c.KafkaProducerTimeout <= 0 {
		return errors.New("KAFKA_PRODUCER_TIMEOUT must be > 0")
	}
	if c.GameMinPlayers <= 0 {
		return errors.New("GAME_MIN_PLAYERS must be > 0")
	}
	if c.GameMaxPlayers < c.GameMinPlayers {
		return errors.New("GAME_MAX_PLAYERS must be >= GAME_MIN_PLAYERS")
	}
	if c.EngineEventsTopic == "" {
		return errors.New("ENGINE_EVENTS_TOPIC must not be empty")
	}
	if c.PlayerActionsTopic == "" {
		return errors.New("PLAYER_ACTIONS_TOPIC must not be empty")
	}
	return nil
}

// lookupEnvTrim is a small helper that wraps os.LookupEnv and trims spaces.
func lookupEnvTrim(key string) (string, bool) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return "", false
	}
	return strings.TrimSpace(v), true
}

// parseCSV splits a comma-separated string and trims elements.
func parseCommaList(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}
