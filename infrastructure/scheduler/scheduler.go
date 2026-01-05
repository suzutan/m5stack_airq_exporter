package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/suzutan/m5stack_airq_exporter/usecase"
)

// Scheduler handles periodic task execution
type Scheduler struct {
	fetchUsecase *usecase.FetchAirQUsecase
	interval     time.Duration
}

// NewScheduler creates a new scheduler with the given usecase and interval
func NewScheduler(fetchUsecase *usecase.FetchAirQUsecase, interval time.Duration) *Scheduler {
	return &Scheduler{
		fetchUsecase: fetchUsecase,
		interval:     interval,
	}
}

// Start begins the periodic execution of the fetch task
func (s *Scheduler) Start(ctx context.Context) {
	// Execute immediately on start
	s.execute(ctx)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduler stopped")
			return
		case <-ticker.C:
			s.execute(ctx)
		}
	}
}

func (s *Scheduler) execute(ctx context.Context) {
	if err := s.fetchUsecase.Execute(ctx); err != nil {
		log.Printf("Failed to fetch air quality data: %v", err)
	} else {
		log.Println("Successfully fetched air quality data")
	}
}
