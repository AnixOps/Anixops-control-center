# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Mobile support preparation (Android/iOS)
- GitHub Actions CI/CD workflows
- Branch protection configuration

## [v0.1.0] - 2026-03-15

### Added

- **Core Plugin System**
  - Plugin interface with lifecycle management
  - Plugin manager for registration and execution
  - Plugin registry for factory pattern

- **Plugins**
  - Ansible plugin for playbook execution
  - v2board plugin for panel management
  - V2bX plugin for proxy node management
  - AnixOps-agent plugin for remote control

- **API Layer**
  - REST API server with Gin framework
  - WebSocket hub for real-time communication
  - JWT authentication
  - RBAC authorization (admin/operator/viewer)

- **User Interfaces**
  - TUI with Bubble Tea framework
  - Web GUI with Vue 3 and Tailwind CSS

- **Security**
  - JWT token management
  - Password hashing with bcrypt
  - Role-based access control

- **Infrastructure**
  - SQLite database support
  - Event bus for pub/sub
  - Task scheduler with cron

- **Deployment**
  - Docker support
  - systemd service file
  - Multi-platform builds (Windows/macOS/Linux)

- **Testing**
  - Unit tests for core modules
  - Coverage reporting

### Security

- JWT-based authentication
- RBAC with three default roles
- Password hashing with bcrypt

---

## Version Naming Convention

- **Major (X.0.0)**: Breaking changes
- **Minor (0.X.0)**: New features, backward compatible
- **Patch (0.0.X)**: Bug fixes, backward compatible

## Branch Strategy

- **dev**: Development branch (RC versions)
- **production**: Stable releases
- **vX.X.X**: Frozen version branches (protected)

---

[Unreleased]: https://github.com/AnixOps/anixops-control-center/compare/v0.1.0...HEAD
[v0.1.0]: https://github.com/AnixOps/anixops-control-center/releases/tag/v0.1.0