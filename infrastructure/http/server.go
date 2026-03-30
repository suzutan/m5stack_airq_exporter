package http

import (
	"context"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
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

	// Middleware
	e.Use(middleware.RequestLogger())
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

// Start starts the HTTP server with graceful shutdown support.
// It blocks until the context is canceled and the server shuts down.
func (s *Server) Start(ctx context.Context, address string) error {
	sc := echo.StartConfig{
		Address:    address,
		HideBanner: true,
	}
	return sc.Start(ctx, s.echo)
}

// Echo returns the underlying echo instance (for testing)
func (s *Server) Echo() *echo.Echo {
	return s.echo
}
