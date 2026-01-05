package usecase

import (
	"context"
	"fmt"

	"github.com/suzutan/m5stack_airq_exporter/domain/repository"
)

// FetchAirQUsecase handles the business logic for fetching air quality data
type FetchAirQUsecase struct {
	airqRepo    repository.AirQRepository
	metricsRepo repository.MetricsRepository
}

// NewFetchAirQUsecase creates a new FetchAirQUsecase with the given dependencies
func NewFetchAirQUsecase(
	airqRepo repository.AirQRepository,
	metricsRepo repository.MetricsRepository,
) *FetchAirQUsecase {
	return &FetchAirQUsecase{
		airqRepo:    airqRepo,
		metricsRepo: metricsRepo,
	}
}

// Execute fetches air quality data and updates the metrics
func (u *FetchAirQUsecase) Execute(ctx context.Context) error {
	data, err := u.airqRepo.Fetch(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch air quality data: %w", err)
	}

	u.metricsRepo.Update(data)
	return nil
}
