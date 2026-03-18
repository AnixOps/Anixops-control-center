# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build, Test, and Development Commands

```bash
# Build
make build              # Build CLI binary (anixops)
make build-tui          # Build TUI binary (anixops-tui)
make build-all          # Build both

# Run
./anixops tui           # Launch TUI interface
./anixops server -c config.yaml   # Start API server
./anixops ansible run deploy.yml -i inventory/hosts  # Run playbook

# Test
make test               # Run all tests
make test-coverage      # Generate coverage report (coverage.html)
make test-critical      # Run critical function tests (required 100% for production)
make test-race          # Run with race detection

# Run specific test
go test -v -run TestName ./path/to/package
go test -v -run "TestRegister|TestStartPlugin" ./internal/core/plugin/...

# Coverage summary
go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | grep total

# Quality
make lint               # Run golangci-lint
make fmt                # Format code
make vet                # Run go vet
make check              # All quality checks (fmt + vet + lint + test)

# Dependencies
make deps               # Download and tidy dependencies
go mod tidy             # Tidy modules
```

## Test Coverage Requirements

| Branch | Minimum | Critical Tests |
|--------|---------|----------------|
| dev | 10% | Not required |
| production | 80% | 100% required |
| v*.*.* (frozen) | 80% | 100% required |

**Commercial SDK Standards:**
- Core packages: ≥90%
- Utility packages: ≥80%
- Minimum: ≥70%

## Architecture Overview

```
cmd/
├── anixops/          # CLI entry point
└── anixops-tui/      # TUI entry point

internal/
├── core/             # SDK core (≥90% coverage target)
│   ├── anixops/      # SDK unified entry point
│   ├── config/       # Configuration loading/validation
│   ├── container/    # Dependency injection container
│   ├── errors/       # Structured error codes (E0000-E0899)
│   ├── eventbus/     # Event pub/sub for component communication
│   ├── health/       # Health checks (ReadyChecker, LiveChecker for K8s)
│   ├── hotreload/    # Configuration hot reload with rollback
│   ├── lifecycle/    # Lifecycle hooks (OnStart, OnStop, etc.)
│   ├── logger/       # Structured logging (JSON/text)
│   ├── metrics/      # Metrics collection (Counter, Gauge, Histogram)
│   ├── plugin/       # Plugin system (interface.go, manager.go, registry.go)
│   ├── scheduler/    # Cron-based task scheduling
│   ├── shutdown/     # Graceful shutdown handling
│   ├── state/        # State machine (10 states, strict transitions)
│   └── tracing/      # Distributed tracing (OpenTelemetry-style)
│
├── plugins/          # Built-in plugins
│   ├── ansible/      # Ansible automation
│   ├── v2board/      # v2board panel management
│   ├── v2bx/         # V2bX node management
│   └── agent/        # AnixOps-agent remote control
│
├── api/              # API layer
│   ├── rest/         # REST API (Gin)
│   └── websocket/    # WebSocket hub
│
├── storage/          # Data layer
│   └── sqlite/       # SQLite storage (GORM)
│
└── security/         # Security (auth/ in critical tests)
```

## Plugin Development

All plugins implement the `Plugin` interface from `internal/core/plugin/interface.go`:

```go
type Plugin interface {
    Info() PluginInfo
    Init(ctx context.Context, config map[string]interface{}) error
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    HealthCheck(ctx context.Context) error
    Capabilities() []string
}
```

Optional interfaces:
- `ExecutablePlugin` - for action execution
- `ObservablePlugin` - for status reporting and event streaming

## Key Design Patterns

1. **Plugin-first**: All extensibility through plugins
2. **Event-driven**: Components communicate via EventBus
3. **Dependency injection**: Use constructor injection via Container
4. **Interface segregation**: Small, focused interfaces
5. **State machine**: 10 states with strict transition rules (None→Created→Initializing→Initialized→Starting→Running→Stopping→Stopped, Error, Destroyed)

## Error Codes

- E0000-E0099: General errors
- E0100-E0199: Plugin errors
- E0200-E0299: Authentication errors
- E0300-E0399: Database errors
- E0400-E0499: Node errors
- E0500-E0599: Task errors
- E0600-E0699: Configuration errors
- E0700-E0799: Tracing errors
- E0800-E0899: Rate limiting errors

## Branch Strategy

- `dev`: Development, RC releases, 1 review required
- `production`: Stable releases, 2 reviews + CODEOWNER, 80% coverage
- `v*.*.*`: Frozen immutable version branches

## Critical Functions (Must Pass 100%)

From DEVELOPMENT.md, these tests require 100% pass rate for production:
- Plugin: `TestRegister`, `TestStartPlugin`, `TestStopPlugin`
- Auth: `TestValidateToken`, `TestCheckPassword`, `TestHasPermission`
- Config: `TestDefaultConfig`, `TestLoadFromString`

## Go Version

Go 1.24 is required (specified in go.mod and CI workflows).