package repository

import (
	"context"

	"github.com/suzutan/m5stack_airq_exporter/domain/entity"
)

// AirQRepository defines the interface for fetching air quality data
type AirQRepository interface {
	// Fetch retrieves the latest air quality data from the data source
	Fetch(ctx context.Context) (*entity.AirQuality, error)
}
