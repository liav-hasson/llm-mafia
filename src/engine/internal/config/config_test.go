package config

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.EngineEventsTopic != "engine.events" {
		t.Fatalf("expected default EngineEventsTopic 'engine.events', got %q", cfg.EngineEventsTopic)
	}
	if cfg.PlayerActionsTopic != "player.actions" {
		t.Fatalf("expected default PlayerActionsTopic 'player.actions', got %q", cfg.PlayerActionsTopic)
	}
	if cfg.GameMinPlayers != 6 {
		t.Fatalf("expected default GameMinPlayers 6, got %d", cfg.GameMinPlayers)
	}
	if cfg.KafkaConsumerTimeout != 2*time.Second {
		t.Fatalf("expected default KafkaConsumerTimeout 2s, got %v", cfg.KafkaConsumerTimeout)
	}
}

func TestLoadConfigEnvOverrides(t *testing.T) {
	t.Setenv("KAFKA_BROKERS", "b1:9092,b2:9092")
	t.Setenv("ENGINE_EVENTS_TOPIC", "custom.events")
	t.Setenv("PLAYER_ACTIONS_TOPIC", "custom.actions")
	t.Setenv("GAME_MIN_PLAYERS", "4")
	t.Setenv("GAME_MAX_PLAYERS", "8")
	t.Setenv("KAFKA_CONSUMER_TIMEOUT", "3s")
	t.Setenv("ENABLE_ROLE_SECRETS", "true")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if len(cfg.KafkaBrokers) != 2 {
		t.Fatalf("expected 2 kafka brokers, got %d", len(cfg.KafkaBrokers))
	}
	if cfg.EngineEventsTopic != "custom.events" {
		t.Fatalf("expected EngineEventsTopic 'custom.events', got %q", cfg.EngineEventsTopic)
	}
	if cfg.PlayerActionsTopic != "custom.actions" {
		t.Fatalf("expected PlayerActionsTopic 'custom.actions', got %q", cfg.PlayerActionsTopic)
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
}

func TestLoadConfigInvalidValues(t *testing.T) {
	t.Setenv("KAFKA_CONSUMER_TIMEOUT", "not-a-duration")
	_, err := LoadConfig()
	if err == nil {
		t.Fatalf("expected error for invalid KAFKA_CONSUMER_TIMEOUT, got nil")
	}
}
