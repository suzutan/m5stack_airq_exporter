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
go test -v -run TestFunctionName .

# Helm install
helm install m5stack-airq-exporter ./charts/m5stack-airq-exporter --set config.airqDataUrl=<URL>
```

## Release Process

Releases are automated via conventional commits. When a PR is merged to master:

1. `auto-release.yaml` parses commit messages since the last tag
2. Determines version bump: `feat:` → minor, `fix:` → patch, `!`/`BREAKING CHANGE` → major
3. If a bump is needed: creates tag, builds/pushes Docker image, updates Helm chart, creates GitHub Release
4. Commits without `feat:`/`fix:` prefix (e.g. `chore:`, `docs:`, `ci:`) do **not** trigger a release

```bash
# Manual release (alternative, triggers release.yaml via tag)
task release VERSION=0.3.0
```

Helm chart uses `appVersion` as the default image tag.

## Architecture

Flat single-package structure using `net/http` (no framework).

```
├── main.go           # Entrypoint + HTTP server + scheduler
├── airq.go           # AirQuality struct + API fetch function
├── airq_test.go      # API fetch tests (7 scenarios)
├── metrics.go        # Prometheus metrics definition & update
├── metrics_test.go   # Metrics tests (2 scenarios)
├── charts/m5stack-airq-exporter/  # Helm chart for Kubernetes
└── .github/workflows/             # CI/CD pipelines
```

### Data Flow

1. **Scheduler goroutine** (1-minute interval) calls `FetchAirQuality()` to fetch JSON from M5Stack API
2. **`Metrics.Update()`** updates 11 Prometheus gauge metrics
3. **`net/http` server** serves `/metrics`, `/healthz`, `/readyz` endpoints

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `AIRQ_DATA_URL` | Yes | M5Stack AirQ data endpoint URL |
| `PORT` | No | HTTP server port (default: 8080) |

## Testing Strategy

- `httptest.Server` for API fetch tests
- `prometheus/testutil` for metrics tests
