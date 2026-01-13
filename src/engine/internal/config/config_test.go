package config

import (
	"testing"
	"time"
)

func TestLoadConfigDefaults(t *testing.T) {
	// Test that defaults work
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig with defaults failed: %v", err)
	}

	// Check some key defaults
	if len(cfg.KafkaBrokers) != 1 || cfg.KafkaBrokers[0] != "localhost:9092" {
		t.Errorf("expected default broker localhost:9092, got %v", cfg.KafkaBrokers)
	}
	if cfg.GameMinPlayers != 6 {
		t.Errorf("expected default GameMinPlayers 6, got %d", cfg.GameMinPlayers)
	}
	if cfg.GameMaxPlayers != 12 {
		t.Errorf("expected default GameMaxPlayers 12, got %d", cfg.GameMaxPlayers)
	}
	if cfg.PhaseNightTimeout != 2*time.Minute {
		t.Errorf("expected default PhaseNightTimeout 2m, got %v", cfg.PhaseNightTimeout)
	}
	if cfg.GameIDPrefix != "game" {
		t.Errorf("expected default GameIDPrefix 'game', got %q", cfg.GameIDPrefix)
	}
}

func TestLoadConfigEnvOverrides(t *testing.T) {
	t.Setenv("ENGINE_KAFKA_BROKERS", "b1:9092,b2:9092")
	t.Setenv("ENGINE_GAME_MIN_PLAYERS", "4")
	t.Setenv("ENGINE_GAME_MAX_PLAYERS", "8")
	t.Setenv("ENGINE_KAFKA_CONSUMER_TIMEOUT", "3s")
	t.Setenv("ENGINE_ENABLE_ROLE_SECRETS", "true")
	t.Setenv("ENGINE_PHASE_DAY_TIMEOUT", "10m")
	t.Setenv("ENGINE_GAME_ID_PREFIX", "test")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if len(cfg.KafkaBrokers) != 2 {
		t.Fatalf("expected 2 kafka brokers, got %d", len(cfg.KafkaBrokers))
	}
	if cfg.GameMinPlayers != 4 || cfg.GameMaxPlayers != 8 {
		t.Fatalf("expected GameMinPlayers 4 and GameMaxPlayers 8, got %d/%d", cfg.GameMinPlayers, cfg.GameMaxPlayers)
	}
	if cfg.KafkaConsumerTimeout != 3*time.Second {
		t.Fatalf("expected KafkaConsumerTimeout 3s, got %v", cfg.KafkaConsumerTimeout)
	}
	if !cfg.EnableRoleSecrets {
		t.Fatalf("expected EnableRoleSecrets true")
	}
	if cfg.PhaseDayTimeout != 10*time.Minute {
		t.Fatalf("expected PhaseDayTimeout 10m, got %v", cfg.PhaseDayTimeout)
	}
	if cfg.GameIDPrefix != "test" {
		t.Fatalf("expected GameIDPrefix 'test', got %q", cfg.GameIDPrefix)
	}
}

func TestLoadConfigInvalidValues(t *testing.T) {
	t.Setenv("ENGINE_KAFKA_CONSUMER_TIMEOUT", "not-a-duration")
	_, err := LoadConfig()
	if err == nil {
		t.Fatalf("expected error for invalid ENGINE_KAFKA_CONSUMER_TIMEOUT, got nil")
	}
}
