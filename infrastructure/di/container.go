package di

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/suzutan/m5stack_airq_exporter/adapter/gateway"
	"github.com/suzutan/m5stack_airq_exporter/adapter/handler"
	"github.com/suzutan/m5stack_airq_exporter/domain/repository"
	"github.com/suzutan/m5stack_airq_exporter/usecase"
)

// Config holds the configuration for the application
type Config struct {
	AirQDataURL string
	Port        string
}

// Container holds all dependencies for the application
type Container struct {
	// Config
	Config *Config

	// Repositories
	AirQRepository    repository.AirQRepository
	MetricsRepository repository.MetricsRepository

	// Usecases
	FetchAirQUsecase *usecase.FetchAirQUsecase

	// Handlers
	MetricsHandler *handler.MetricsHandler
	HealthHandler  *handler.HealthHandler

	// Prometheus
	Registry *prometheus.Registry
}

// NewContainer creates a new dependency injection container
func NewContainer(config *Config) *Container {
	// Create Prometheus registry
	registry := prometheus.NewRegistry()

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create repositories
	airqRepo := gateway.NewAirQHTTPGateway(config.AirQDataURL, httpClient)
	metricsRepo := gateway.NewPrometheusMetricsGateway(registry)

	// Create usecases
	fetchAirQUsecase := usecase.NewFetchAirQUsecase(airqRepo, metricsRepo)

	// Create handlers
	metricsHandler := handler.NewMetricsHandler(registry)
	healthHandler := handler.NewHealthHandler()

	return &Container{
		Config:            config,
		AirQRepository:    airqRepo,
		MetricsRepository: metricsRepo,
		FetchAirQUsecase:  fetchAirQUsecase,
		MetricsHandler:    metricsHandler,
		HealthHandler:     healthHandler,
		Registry:          registry,
	}
}
