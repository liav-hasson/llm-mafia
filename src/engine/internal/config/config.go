package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	// -------------
	// Kafka
	// -------------

	// ENGINE_KAFKA_BROKERS="broker1:9092,broker2:9092"
	// envDefault is used instead of default (doesnt work for this lib)
	KafkaBrokers []string `env:"ENGINE_KAFKA_BROKERS" envSeparator:"," envDefault:"localhost:9092"`

	// Optional, useful for debugging
	KafkaClientID string `env:"ENGINE_KAFKA_CLIENT_ID" envDefault:"mafia-engine"`
	KafkaGroupID  string `env:"ENGINE_KAFKA_GROUP_ID" envDefault:"mafia-engine-group"`

	// NOTE: Topic names are constants in kafka/topics.go (single source of truth)
	// Do NOT add topic configuration here to avoid mismatch bugs.

	// Timeouts
	KafkaConsumerTimeout time.Duration `env:"ENGINE_KAFKA_CONSUMER_TIMEOUT" envDefault:"2s"`
	KafkaProducerTimeout time.Duration `env:"ENGINE_KAFKA_PRODUCER_TIMEOUT" envDefault:"2s"`

	// External calls (e.g. Ollama)
	HTTPTimeout time.Duration `env:"ENGINE_HTTP_TIMEOUT" envDefault:"5s"`

	// -------------
	// Game
	// -------------

	GameMinPlayers int      `env:"ENGINE_GAME_MIN_PLAYERS" envDefault:"6"`
	GameMaxPlayers int      `env:"ENGINE_GAME_MAX_PLAYERS" envDefault:"12"`
	GameIDPrefix   string   `env:"ENGINE_GAME_ID_PREFIX" envDefault:"game"`
	PlayerNames    []string `env:"ENGINE_PLAYER_NAMES" envSeparator:"," envDefault:"Gilbert McDonald,Dorothy Bird,Ernest Preston,Vincent Schultz,Joanne Sloan,Lana Moran,Adrienne Fuller,Greg Bennett,Curt Simon,Rachel McMillan,Dustin Eastman,Willard Mendez"`

	// Phase timeouts (how long each phase lasts before auto-advancing)
	PhaseNightTimeout  time.Duration `env:"ENGINE_PHASE_NIGHT_TIMEOUT" envDefault:"2m"`
	PhaseDayTimeout    time.Duration `env:"ENGINE_PHASE_DAY_TIMEOUT" envDefault:"5m"`
	PhaseVotingTimeout time.Duration `env:"ENGINE_PHASE_VOTING_TIMEOUT" envDefault:"1m"`

	// mock | llm
	AgentMode string `env:"ENGINE_AGENT_MODE" envDefault:"mock"`

	// -------------
	// Logging
	// -------------

	LogLevel string `env:"ENGINE_LOG_LEVEL" envDefault:"info"`
	Env      string `env:"ENGINE_ENV" envDefault:"dev"`

	// Feature flags
	EnableRoleSecrets bool `env:"ENGINE_ENABLE_ROLE_SECRETS" envDefault:"false"`
}

// Load loads configuration from environment variables,
// applies defaults, and validates the result.
func LoadConfig() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate checks config sanity and returns an error for invalid settings.
func (c *Config) Validate() error {
	if len(c.KafkaBrokers) == 0 {
		return errors.New("ENGINE_KAFKA_BROKERS must not be empty")
	}

	if c.KafkaConsumerTimeout <= 0 {
		return errors.New("ENGINE_KAFKA_CONSUMER_TIMEOUT must be > 0")
	}

	if c.KafkaProducerTimeout <= 0 {
		return errors.New("ENGINE_KAFKA_PRODUCER_TIMEOUT must be > 0")
	}

	if c.HTTPTimeout <= 0 {
		return errors.New("ENGINE_HTTP_TIMEOUT must be > 0")
	}

	if c.GameMinPlayers <= 0 {
		return errors.New("ENGINE_GAME_MIN_PLAYERS must be > 0")
	}

	if c.GameMaxPlayers < c.GameMinPlayers {
		return errors.New("ENGINE_GAME_MAX_PLAYERS must be >= ENGINE_GAME_MIN_PLAYERS")
	}

	if c.GameIDPrefix == "" {
		return errors.New("ENGINE_GAME_ID_PREFIX must not be empty")
	}

	if len(c.PlayerNames) == 0 {
		return errors.New("ENGINE_PLAYER_NAMES must not be empty")
	}

	if c.PhaseNightTimeout <= 0 {
		return errors.New("ENGINE_PHASE_NIGHT_TIMEOUT must be > 0")
	}

	if c.PhaseDayTimeout <= 0 {
		return errors.New("ENGINE_PHASE_DAY_TIMEOUT must be > 0")
	}

	if c.PhaseVotingTimeout <= 0 {
		return errors.New("ENGINE_PHASE_VOTING_TIMEOUT must be > 0")
	}

	switch c.AgentMode {
	case "mock", "llm":
		// ok
	default:
		return fmt.Errorf("ENGINE_AGENT_MODE must be one of [mock, llm], got %q", c.AgentMode)
	}

	return nil
}
