# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

M5Stack AirQ Prometheus Exporter - A Go application that fetches air quality data from M5Stack AirQ devices and exposes it as Prometheus metrics.

## Build & Development Commands

Commands are managed via [Task](https://taskfile.dev/). Run `task --list` to see all available tasks.

```bash
task ci             # Run full CI pipeline (clean, tidy, fmt, vet, test, build)

task build          # Build the application
task run            # Run locally with default AIRQ_DATA_URL
task test           # Run all tests
task test:coverage  # Run tests with coverage report
task vet            # Run go vet
task lint           # Run golangci-lint (requires installation)
task fmt            # Format code
task tidy           # Tidy go modules
task clean          # Clean build artifacts

task docker:build   # Build Docker image
task docker:run     # Run Docker container

task helm:lint      # Lint Helm chart
task helm:template  # Render Helm chart templates

task release VERSION=x.y.z  # Create a new release
```

### Direct Commands

```bash
# Run single test
go test -v -run TestFunctionName ./path/to/package

# Helm install
helm install m5stack-airq-exporter ./charts/m5stack-airq-exporter --set config.airqDataUrl=<URL>
```

## Release Process

```bash
# 1. Create release (updates Chart.yaml, commits, and tags)
task release VERSION=0.1.0

# 2. Push to trigger GitHub Actions
git push && git push --tags
```

GitHub Actions will:
- Run tests
- Build and push Docker image to `ghcr.io/suzutan/m5stack_airq_exporter:v0.1.0`

Helm chart uses `appVersion` as the default image tag.

## Architecture (Clean Architecture + DI)

```
├── cmd/exporter/              # Application entrypoint with DI wiring
├── domain/
│   ├── entity/                # Business entities (AirQuality)
│   └── repository/            # Repository interfaces (ports)
├── usecase/                   # Application business logic (interactors)
├── adapter/
│   ├── gateway/               # Repository implementations (AirQ HTTP, Prometheus)
│   └── handler/               # HTTP handlers (metrics, health)
├── infrastructure/
│   ├── di/                    # Dependency injection container
│   ├── http/                  # Echo server setup
│   └── scheduler/             # Periodic task scheduler
├── charts/m5stack-airq-exporter/  # Helm chart for Kubernetes
└── .github/workflows/         # CI/CD pipelines
```

### Layer Dependencies

```
cmd/exporter → infrastructure/di → usecase → domain
                                 → adapter  → domain
```

### Data Flow

1. **Scheduler** (1-minute interval) triggers `FetchAirQUsecase.Execute()`
2. **AirQHTTPGateway** (implements `AirQRepository`) fetches JSON from API
3. **PrometheusMetricsGateway** (implements `MetricsRepository`) updates gauge metrics
4. **Echo Server** serves `/metrics` endpoint via `MetricsHandler`

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `AIRQ_DATA_URL` | Yes | M5Stack AirQ data endpoint URL |
| `PORT` | No | HTTP server port (default: 8080) |

## Testing Strategy

- Repository interfaces enable easy mocking in usecase tests
- `httptest.Server` for HTTP gateway tests
- `prometheus/testutil` for metrics gateway tests
