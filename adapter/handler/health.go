package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HandleLiveness handles the /healthz endpoint for liveness probe
func (h *HealthHandler) HandleLiveness(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// HandleReadiness handles the /readyz endpoint for readiness probe
func (h *HealthHandler) HandleReadiness(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
