package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/suzutan/m5stack_airq_exporter/infrastructure/di"
	"github.com/suzutan/m5stack_airq_exporter/infrastructure/http"
	"github.com/suzutan/m5stack_airq_exporter/infrastructure/scheduler"
)

func main() {
	// Load configuration from environment variables
	config := &di.Config{
		AirQDataURL: getEnv("AIRQ_DATA_URL", ""),
		Port:        getEnv("PORT", "8080"),
	}

	if config.AirQDataURL == "" {
		log.Fatal("AIRQ_DATA_URL environment variable is required")
	}

	// Create dependency injection container
	container := di.NewContainer(config)

	// Create HTTP server
	server := http.NewServer(container)

	// Create scheduler for periodic data fetch (1 minute interval)
	sched := scheduler.NewScheduler(container.FetchAirQUsecase, 1*time.Minute)

	// Create context that will be canceled on shutdown signal
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Start scheduler in background
	go sched.Start(ctx)

	// Start HTTP server (blocks until context is canceled and graceful shutdown completes)
	log.Printf("Starting server on :%s", config.Port)
	if err := server.Start(ctx, ":"+config.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Println("Server stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
