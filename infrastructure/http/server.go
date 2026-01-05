package http

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/suzutan/m5stack_airq_exporter/infrastructure/di"
)

// Server represents the HTTP server
type Server struct {
	echo      *echo.Echo
	container *di.Container
}

// NewServer creates a new HTTP server with the given container
func NewServer(container *di.Container) *Server {
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/metrics", container.MetricsHandler.Handle)
	e.GET("/healthz", container.HealthHandler.HandleLiveness)
	e.GET("/readyz", container.HealthHandler.HandleReadiness)

	return &Server{
		echo:      e,
		container: container,
	}
}

// Start starts the HTTP server
func (s *Server) Start(address string) error {
	return s.echo.Start(address)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}

// Echo returns the underlying echo instance (for testing)
func (s *Server) Echo() *echo.Echo {
	return s.echo
}
