# AnixOps Control Center

A unified TUI/GUI control center for managing all AnixOps products.

## Features

- **Plugin Architecture**: Modular plugin system for extensibility
- **Multiple Interfaces**: TUI, Web GUI, CLI, REST API
- **Enterprise Security**: JWT, OAuth, LDAP, SAML authentication with RBAC
- **Product Integration**:
  - Ansible automation engine
  - v2board panel management
  - V2bX node management
  - AnixOps-agent remote control

## Quick Start

```bash
# Build
go build -o anixops cmd/anixops/main.go

# Run TUI
./anixops tui

# Start API server
./anixops server -c config.yaml

# Run Ansible playbook
./anixops ansible run deploy.yml -i inventory/hosts

# List nodes
./anixops nodes list
```

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        ANIXOPS CONTROL CENTER                               │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                      │
│  │ TUI Interface │  │  Web GUI     │  │ CLI Interface│                      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘                      │
│         └─────────────────┼─────────────────┘                               │
│                           ▼                                                 │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │                        Core API Layer (Go)                            │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
│                           ▼                                                 │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │                      Plugin Manager (Core)                            │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
│                           ▼                                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │
│  │ Ansible     │  │ v2board     │  │ V2bX        │  │ Agent       │       │
│  │ Plugin      │  │ Plugin      │  │ Plugin      │  │ Plugin      │       │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘       │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Plugin Development

```go
package myplugin

import (
    "context"
    "github.com/anixops/anixops-control-center/internal/core/plugin"
)

type MyPlugin struct{}

func (p *MyPlugin) Info() plugin.PluginInfo {
    return plugin.PluginInfo{
        Name:        "myplugin",
        Version:     "1.0.0",
        Description: "My custom plugin",
    }
}

func (p *MyPlugin) Init(ctx context.Context, config map[string]interface{}) error {
    return nil
}

func (p *MyPlugin) Start(ctx context.Context) error {
    return nil
}

func (p *MyPlugin) Stop(ctx context.Context) error {
    return nil
}

func (p *MyPlugin) HealthCheck(ctx context.Context) error {
    return nil
}

func (p *MyPlugin) Capabilities() []string {
    return []string{"do_something"}
}

// Optional: Implement ExecutablePlugin for action execution
func (p *MyPlugin) Execute(ctx context.Context, action string, params map[string]interface{}) (plugin.Result, error) {
    return plugin.Result{Success: true}, nil
}
```

## API Endpoints

### Auth
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh token

### Plugins
- `GET /api/v1/plugins` - List plugins
- `GET /api/v1/plugins/:name` - Get plugin info
- `POST /api/v1/plugins/:name/execute` - Execute plugin action

### Nodes
- `GET /api/v1/nodes` - List nodes
- `POST /api/v1/nodes` - Create node
- `DELETE /api/v1/nodes/:id` - Delete node

### Playbooks
- `GET /api/v1/playbooks` - List playbooks
- `POST /api/v1/playbooks/run` - Run playbook

### Users
- `GET /api/v1/users` - List users
- `POST /api/v1/users/:id/ban` - Ban user
- `POST /api/v1/users/:id/unban` - Unban user

## Configuration

See `configs/config.yaml` for configuration options.

## License

MIT