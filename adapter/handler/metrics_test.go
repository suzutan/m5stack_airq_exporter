package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

func TestMetricsHandler_Handle(t *testing.T) {
	e := echo.New()
	registry := prometheus.NewRegistry()

	// Register a test metric
	testGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "test_metric",
		Help: "A test metric",
	})
	registry.MustRegister(testGauge)
	testGauge.Set(42.0)

	handler := NewMetricsHandler(registry)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Handle(c)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "test_metric 42") {
		t.Errorf("expected body to contain test_metric 42, got %s", body)
	}
}

func TestMetricsHandler_ContentType(t *testing.T) {
	e := echo.New()
	registry := prometheus.NewRegistry()

	handler := NewMetricsHandler(registry)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler.Handle(c)

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/plain") {
		t.Errorf("expected Content-Type to contain text/plain, got %s", contentType)
	}
}
