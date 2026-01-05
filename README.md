# M5Stack AirQ Prometheus Exporter

[![CI](https://github.com/suzutan/m5stack_airq_exporter/actions/workflows/ci.yaml/badge.svg)](https://github.com/suzutan/m5stack_airq_exporter/actions/workflows/ci.yaml)
[![Release](https://github.com/suzutan/m5stack_airq_exporter/actions/workflows/release.yaml/badge.svg)](https://github.com/suzutan/m5stack_airq_exporter/actions/workflows/release.yaml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/suzutan/m5stack_airq_exporter)](https://go.dev/)
[![License](https://img.shields.io/github/license/suzutan/m5stack_airq_exporter)](LICENSE)
[![Container Image](https://img.shields.io/badge/ghcr.io-suzutan%2Fm5stack__airq__exporter-blue)](https://github.com/suzutan/m5stack_airq_exporter/pkgs/container/m5stack_airq_exporter)

A Prometheus exporter for [M5Stack AirQ](https://docs.m5stack.com/en/unit/airq) air quality monitoring device. Fetches sensor data from the M5Stack EzData API and exposes it as Prometheus metrics.

## Features

- Exports air quality metrics from M5Stack AirQ (SEN55 + SCD40 sensors)
- 1-minute automatic data fetch interval
- Prometheus-compatible `/metrics` endpoint
- Health check endpoints (`/healthz`, `/readyz`)
- Multi-architecture Docker image (amd64, arm64)
- Helm chart with ServiceMonitor support for Prometheus Operator
- Clean Architecture with Dependency Injection
- Graceful shutdown handling

## Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `airq_pm1_0` | Gauge | PM1.0 particulate matter (μg/m³) |
| `airq_pm2_5` | Gauge | PM2.5 particulate matter (μg/m³) |
| `airq_pm4_0` | Gauge | PM4.0 particulate matter (μg/m³) |
| `airq_pm10_0` | Gauge | PM10.0 particulate matter (μg/m³) |
| `airq_humidity` | Gauge | Relative humidity from SEN55 (%) |
| `airq_temperature` | Gauge | Temperature from SEN55 (°C) |
| `airq_voc` | Gauge | Volatile organic compounds index |
| `airq_nox` | Gauge | Nitrogen oxides index |
| `airq_co2` | Gauge | CO2 concentration from SCD40 (ppm) |
| `airq_scd40_humidity` | Gauge | Relative humidity from SCD40 (%) |
| `airq_scd40_temperature` | Gauge | Temperature from SCD40 (°C) |

## Quick Start

### Prerequisites

- M5Stack AirQ device connected to [EzData](https://ezdata.m5stack.com/)
- Your EzData API URL (format: `https://ezdata2.m5stack.com/api/v2/{TOKEN}/dataMacByKey/raw`)

### Using Docker

```bash
docker run -d \
  -p 8080:8080 \
  -e AIRQ_DATA_URL=https://ezdata2.m5stack.com/api/v2/YOUR_TOKEN/dataMacByKey/raw \
  ghcr.io/suzutan/m5stack_airq_exporter:latest
```

### Using Helm

```bash
helm install m5stack-airq-exporter ./charts/m5stack-airq-exporter \
  --set config.airqDataUrl=https://ezdata2.m5stack.com/api/v2/YOUR_TOKEN/dataMacByKey/raw
```

## Configuration

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `AIRQ_DATA_URL` | Yes | - | M5Stack EzData API endpoint URL |
| `PORT` | No | `8080` | HTTP server listen port |

### Helm Values

See [values.yaml](./charts/m5stack-airq-exporter/values.yaml) for all available options.

Key configuration:

```yaml
config:
  airqDataUrl: "https://ezdata2.m5stack.com/api/v2/YOUR_TOKEN/dataMacByKey/raw"
  port: "8080"

serviceMonitor:
  enabled: true        # Enable for Prometheus Operator
  interval: 60s
  scrapeTimeout: 10s
```

## API Response Format

This exporter works with any HTTP endpoint that returns data in the following format. You can use a custom endpoint or proxy as long as it conforms to this structure.

### Expected Response Structure

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "dataToken": "...",
    "dataType": "...",
    "name": "...",
    "value": "{\"sen55\":{...},\"scd40\":{...},\"profile\":{...}}",
    "createTime": "...",
    "updateTime": "..."
  }
}
```

### Sensor Data Format (inside `data.value`)

The `value` field contains a JSON string (may be escaped) with the following structure:

```json
{
  "sen55": {
    "pm1.0": 3.0,
    "pm2.5": 4.0,
    "pm4.0": 4.0,
    "pm10.0": 4.0,
    "humidity": 47.5,
    "temperature": 25.3,
    "voc": 100,
    "nox": 1
  },
  "scd40": {
    "co2": 800,
    "humidity": 45.0,
    "temperature": 26.0
  },
  "profile": {
    "nickname": "Living Room"
  }
}
```

### Field Requirements

| Field | Required | Description |
|-------|----------|-------------|
| `code` | Yes | Must be `200` for success |
| `data.value` | Yes | JSON string containing sensor data |
| `sen55.*` | Optional | SEN55 sensor readings (PM, humidity, temperature, VOC, NOx) |
| `scd40.*` | Optional | SCD40 sensor readings (CO2, humidity, temperature) |
| `profile.nickname` | Optional | Device label (currently not exported as metric) |

Missing sensor fields will result in zero values for those metrics.

## Endpoints

| Path | Description |
|------|-------------|
| `/metrics` | Prometheus metrics endpoint |
| `/healthz` | Liveness probe endpoint |
| `/readyz` | Readiness probe endpoint |

## Development

### Prerequisites

- Go 1.25+
- [Task](https://taskfile.dev/) (optional, for task runner)

### Available Tasks

```bash
task              # Show available tasks
task build        # Build the application
task run          # Run locally
task test         # Run tests
task ci           # Run full CI pipeline (clean, tidy, fmt, vet, test, build)
task helm:lint    # Lint Helm chart
task helm:template # Render Helm templates
```

### Running Locally

```bash
export AIRQ_DATA_URL=https://ezdata2.m5stack.com/api/v2/YOUR_TOKEN/dataMacByKey/raw
task run
```

### Running Tests

```bash
task test

# With coverage
task test:coverage
```

### Creating a Release

```bash
task release VERSION=0.3.0
```

This will:
1. Run the CI pipeline locally
2. Create and push a git tag
3. Trigger GitHub Actions to build and push the Docker image
4. Automatically update the Helm chart version

## Architecture

```
.
├── cmd/exporter/          # Application entrypoint
├── domain/
│   ├── entity/            # Domain entities (AirQuality)
│   └── repository/        # Repository interfaces
├── usecase/               # Business logic (FetchAirQualityUseCase)
├── adapter/
│   ├── gateway/           # External service implementations
│   │   ├── airq_http.go   # M5Stack API client
│   │   └── prometheus_metrics.go
│   └── handler/           # HTTP handlers
├── infrastructure/
│   ├── di/                # Dependency injection container
│   ├── http/              # Echo HTTP server setup
│   └── scheduler/         # Periodic data fetch scheduler
└── charts/                # Helm chart
```

## License

[MIT License](LICENSE)
