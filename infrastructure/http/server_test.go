package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/suzutan/m5stack_airq_exporter/adapter/gateway"
	"github.com/suzutan/m5stack_airq_exporter/adapter/handler"
	"github.com/suzutan/m5stack_airq_exporter/infrastructure/di"
	"github.com/suzutan/m5stack_airq_exporter/usecase"
)

func newTestContainer() *di.Container {
	registry := prometheus.NewRegistry()
	httpClient := &http.Client{}
	airqRepo := gateway.NewAirQHTTPGateway("http://localhost/test", httpClient)
	metricsRepo := gateway.NewPrometheusMetricsGateway(registry)
	fetchUsecase := usecase.NewFetchAirQUsecase(airqRepo, metricsRepo)
	metricsHandler := handler.NewMetricsHandler(registry)
	healthHandler := handler.NewHealthHandler()

	return &di.Container{
		Config:            &di.Config{AirQDataURL: "http://localhost/test", Port: "8080"},
		AirQRepository:    airqRepo,
		MetricsRepository: metricsRepo,
		FetchAirQUsecase:  fetchUsecase,
		MetricsHandler:    metricsHandler,
		HealthHandler:     healthHandler,
		Registry:          registry,
	}
}

func TestNewServer_RoutesRegistered(t *testing.T) {
	container := newTestContainer()
	server := NewServer(container)

	routes := server.Echo().Router().Routes()
	pathMethods := make(map[string]bool)
	for _, r := range routes {
		pathMethods[r.Method+":"+r.Path] = true
	}

	expected := []string{
		"GET:/metrics",
		"GET:/healthz",
		"GET:/readyz",
	}
	for _, e := range expected {
		if !pathMethods[e] {
			t.Errorf("expected route %s to be registered", e)
		}
	}
}

func TestServer_HealthzEndpoint(t *testing.T) {
	container := newTestContainer()
	server := NewServer(container)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	server.Echo().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestServer_ReadyzEndpoint(t *testing.T) {
	container := newTestContainer()
	server := NewServer(container)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	server.Echo().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestServer_MetricsEndpoint(t *testing.T) {
	container := newTestContainer()
	server := NewServer(container)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	server.Echo().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestServer_StartAndShutdown(t *testing.T) {
	container := newTestContainer()
	server := NewServer(container)

	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Start(ctx, ":0")
	}()

	// Cancel context to trigger graceful shutdown
	cancel()

	if err := <-errCh; err != nil {
		t.Fatalf("expected no error from Start, got %v", err)
	}
}
