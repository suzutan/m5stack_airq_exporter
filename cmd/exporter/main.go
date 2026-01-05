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

	// Create context that will be canceled on shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start scheduler in background
	go sched.Start(ctx)

	// Handle graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down...")

		// Create shutdown context with timeout
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		// Shutdown HTTP server
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}

		// Cancel main context to stop scheduler
		cancel()
	}()

	// Start HTTP server
	log.Printf("Starting server on :%s", config.Port)
	if err := server.Start(":" + config.Port); err != nil {
		if err.Error() != "http: Server closed" {
			log.Fatalf("Failed to start server: %v", err)
		}
	}

	log.Println("Server stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
