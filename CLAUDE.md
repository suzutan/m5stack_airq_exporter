# CLAUDE.md

M5Stack AirQ Prometheus Exporter — Go application, flat single-package structure.

## Build & Test

Commands via [Task](https://taskfile.dev/). `task --list` for all options.

```bash
task ci          # Full CI: clean, tidy, fmt, vet, test, build
task test        # Run all tests
task build       # Build binary

# Run single test
go test -v -run TestFunctionName .
```

## Release Process

Automated via conventional commits on master merge (`auto-release.yaml`):
- `feat:` → minor, `fix:` → patch, `!`/`BREAKING CHANGE` → major
- `chore:`, `docs:`, `ci:` → no release
- Manual alternative: `task release VERSION=x.y.z` (triggers `release.yaml` via tag)

## Gotchas

- M5Stack API returns double-escaped JSON in the `value` field — `FetchAirQuality()` handles unescaping via `strconv.Unquote`
- Helm chart `appVersion` is used as the default image tag — kept in sync by CI
