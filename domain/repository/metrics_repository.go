package repository

import "github.com/suzutan/m5stack_airq_exporter/domain/entity"

// MetricsRepository defines the interface for updating metrics
type MetricsRepository interface {
	// Update updates the metrics with the given air quality data
	Update(data *entity.AirQuality)
}
