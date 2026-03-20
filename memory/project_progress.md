---
name: project_progress
description: AnixOps Workers API development progress and version history
type: project
---

# AnixOps Workers API Progress

## Latest Release
**Version**: `v2.0.0-rc.4`
**Date**: 2026-03-20
**Status**: ✅ Released
**GitHub**: https://github.com/AnixOps/Anixops-control-center/releases/tag/v2.0.0-rc.4

## Recent Changes (v2.0.0-rc.4)
- Implement Playbooks feature for mobile app
- Playbooks list with category filtering
- Playbook detail page with YAML content viewer
- Variable editor for playbook execution
- Upload custom playbooks support
- Sync built-in playbooks feature
- 6 built-in playbooks: Fail2ban, Firewall, SSH Hardening, Docker, XRay, System Update

## Recent Changes (v2.0.0-rc.3)
- Fix Windows desktop window controls (minimize/maximize/close)
- Fix language switching (now connected to localeProvider)
- Fix theme switching (now connected to themeModeProvider)
- Update VERSION to v2.0.0-rc.3

## Recent Changes (v2.0.0-rc.2)
- Update permissions: all authenticated users can add nodes
- Show Users menu only for admin role
- Display user role badge in sidebar (ADMIN/OPERATOR/VIEWER)
- Upgrade Flutter to 3.41.5, Dart 3.11.3
- Add web platform support for mobile app
- Fix deprecated APIs (CardTheme → CardThemeData, withOpacity → withValues)
- Configure custom domain api.anixops.com
- API: 342 tests passed, 85.25% coverage
- Build: APK (53.2MB), Web

## Phase Summary

| Phase | Status | Description | RC Version |
|-------|--------|-------------|------------|
| Phase 1 | ✅ | Security Hardening | rc.1 - rc.5 |
| Phase 2 | ✅ | Core Features (Agent + MFA + Ansible) | rc.6 - rc.16 |
| Phase 3 | ✅ | Enterprise Features (Multi-tenancy + Audit) | rc.17 - rc.21 |
| Phase 4 | ✅ | Scalability (Auto-scaling + Load Balancing) | rc.22 - rc.26 |

## v2.x Cloud Native Integration

| Version | Feature | Status | RC Version |
|---------|---------|--------|------------|
| v2.0 | Kubernetes Operator, Helm Charts | ✅ | rc.17 |
| v2.1 | Istio Service Mesh | ✅ | rc.22 |
| v2.2 | ELK Observability Stack | ✅ | rc.24 |
| v2.3 | Resilience (Circuit Breaker, Rate Limit) | ✅ | v1.2 |

## RC Release History

### rc.26 (2026-03-20)
- Load Balancing Service
- 6 algorithms: round-robin, weighted, least-connections, ip-hash, random, response-time
- 342 tests, 85.25% coverage

### rc.25 (2026-03-20)
- Auto-Scaling Service
- Metric-based scaling decisions
- 316 tests, 85.45% coverage

### rc.24 (2026-03-20)
- ELK Observability Configuration
- Grafana dashboard, Logstash pipeline, Elasticsearch templates
- Docker Compose for observability stack

### rc.23 (2026-03-20)
- ELK Log Integration
- Centralized logging, log search, export
- 293 tests, 84.7% coverage

### rc.22 (2026-03-20)
- Istio Service Mesh Integration
- VirtualServices, DestinationRules, Gateways
- Traffic split, circuit breaker, fault injection
- 270 tests, 84.59% coverage

## API Endpoints Summary

Total Endpoints: 100+

### Core
- Auth: 5 endpoints
- Users: 12 endpoints
- Nodes: 12 endpoints

### Features
- Playbooks: 7 endpoints
- Tasks: 6 endpoints
- Schedules: 7 endpoints
- Plugins: 3 endpoints

### Enterprise
- Tenants: 10 endpoints
- Audit: 1 endpoint
- MFA: 7 endpoints
- Agents: 6 endpoints

### Cloud Native
- Kubernetes: 11 endpoints
- Istio: 7 endpoints
- Elasticsearch: 11 endpoints

### Scalability
- Auto-scaling: 14 endpoints
- Load Balancing: 15 endpoints

## Test Coverage History

| Version | Tests | Coverage |
|---------|-------|----------|
| rc.26 | 342 | 85.25% |
| rc.25 | 316 | 85.45% |
| rc.23 | 293 | 84.7% |
| rc.22 | 270 | 84.59% |
| rc.16 | 254 | 84.59% |