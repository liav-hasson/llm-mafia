package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mafia-engine/internal/config"
	"mafia-engine/internal/domain"
	"mafia-engine/internal/engine"
	"mafia-engine/internal/kafka"
)

func main() {
	// Load configuration from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting Mafia Engine with config: brokers=%v, topic=%s, groupID=%s, maxPlayers=%d",
		cfg.KafkaBrokers, cfg.PlayerActionsTopic, cfg.KafkaGroupID, cfg.GameMaxPlayers)

	// Create Kafka producer for publishing authoritative events
	producer, err := kafka.NewKafkaProducer(cfg.KafkaBrokers, cfg.EngineEventsTopic)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	log.Printf("Kafka producer created for topic: %s", cfg.EngineEventsTopic)

	// Create Kafka consumer for receiving player actions
	// Consumer needs topics as a slice
	consumer, err := kafka.NewKafkaConsumer(
		cfg.KafkaBrokers,
		[]string{cfg.PlayerActionsTopic},
		cfg.KafkaGroupID,
	)
	if err != nil {
		if closeErr := producer.Close(); closeErr != nil {
			log.Printf("Error closing producer: %v", closeErr)
		}
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	log.Printf("Kafka consumer created for topic: %s, group: %s", cfg.PlayerActionsTopic, cfg.KafkaGroupID)

	// Initialize game state with configuration
	// Start in Waiting phase (players can join)
	gameState := &domain.GameState{
		ID:      fmt.Sprintf("game-%d", time.Now().Unix()),
		Phase:   domain.PhaseWaiting,
		Round:   0,
		Winner:  domain.WinnerNone,
		Players: make(map[string]*domain.Player),
		Votes:   make(map[string]string),
	}
	log.Printf("Game state initialized: id=%s, phase=%s", gameState.ID, gameState.Phase)

	// Create the game engine
	eng, err := engine.NewEngine(gameState, producer)
	if err != nil {
		if closeErr := consumer.Close(); closeErr != nil {
			log.Printf("Error closing consumer: %v", closeErr)
		}
		if closeErr := producer.Close(); closeErr != nil {
			log.Printf("Error closing producer: %v", closeErr)
		}
		log.Fatalf("Failed to create engine: %v", err)
	}
	log.Println("Game engine created")

	// Start the engine event loop
	eng.Start()
	log.Println("Engine started")

	// Create context for coordinating shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consuming player actions and feeding them to the engine
	// This runs in a goroutine and blocks until context is canceled
	go func() {
		log.Println("Starting consumer loop...")
		if err := consumer.Consume(ctx, eng.HandleMessage); err != nil {
			log.Printf("Consumer error: %v", err)
			cancel() // Signal shutdown on consumer error
		}
	}()

	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	log.Println("Engine is running. Press Ctrl+C to stop.")
	<-sigCh
	log.Println("Shutdown signal received, initiating graceful shutdown...")

	// Graceful shutdown sequence:
	// 1. Cancel context to stop consumer and engine loops
	cancel()

	// 2. Give components time to finish current work
	time.Sleep(1 * time.Second)

	// 3. Stop engine (drains queues, cancels timers)
	log.Println("Stopping engine...")
	eng.Stop()

	// 4. Close Kafka connections
	log.Println("Closing Kafka consumer...")
	if err := consumer.Close(); err != nil {
		log.Printf("Error closing consumer: %v", err)
	}

	log.Println("Closing Kafka producer...")
	if err := producer.Close(); err != nil {
		log.Printf("Error closing producer: %v", err)
	}

	log.Println("Shutdown complete")
	fmt.Println("Mafia Engine stopped successfully")
}
