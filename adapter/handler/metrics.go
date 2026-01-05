package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler handles the /metrics endpoint
type MetricsHandler struct {
	handler echo.HandlerFunc
}

// NewMetricsHandler creates a new MetricsHandler with the given registry
func NewMetricsHandler(registry prometheus.Gatherer) *MetricsHandler {
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	return &MetricsHandler{
		handler: echo.WrapHandler(h),
	}
}

// Handle processes the metrics request
func (h *MetricsHandler) Handle(c echo.Context) error {
	return h.handler(c)
}
