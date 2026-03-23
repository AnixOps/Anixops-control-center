---
name: project_progress
description: AnixOps Workers API development progress and version history
type: project
---

# AnixOps Workers API Progress

## Latest Release
**Version**: `v2.0.0-rc.24`
**Date**: 2026-03-23
**Status**: ✅ Released
**GitHub**: https://github.com/AnixOps/Anixops-control-center/releases/tag/v2.0.0-rc.24

## Test Summary
| Platform | Tests | Status |
|----------|-------|--------|
| Web Frontend | 92 | ✅ Passing |
| Mobile App | 141 | ✅ Passing |
| Workers API | 342 | ✅ Passing |
| **Total** | **575** | ✅ |

## Recent Changes (v2.0.0-rc.24)
- Add login.store.test.js with 9 tests
- Add login form validation tests
- Web: 92, Mobile: 141, Workers: 342 tests

## Recent Changes (v2.0.0-rc.23)
- Add auth_provider_test.dart with 7 tests
- Add mfa_provider_test.dart with 10 tests
- Web: 83, Mobile: 141, Workers: 342 tests

## Recent Changes (v2.0.0-rc.22)
- Add playbooks_provider_test.dart with 10 tests
- Add tasks_provider_test.dart with 12 tests
- Web: 83, Mobile: 124, Workers: 342 tests

## Recent Changes (v2.0.0-rc.21)
- Add agents.store.test.js with 8 tests
- Add settings.store.test.js with 12 tests
- Web: 83, Mobile: 102, Workers: 342 tests

## Recent Changes (v2.0.0-rc.20)
- Add logs.store.test.js with 7 tests
- Add dashboard.store.test.js with 8 tests
- Web: 63, Mobile: 102, Workers: 342 tests

## Recent Changes (v2.0.0-rc.19)
- Add schedules.store.test.js with 9 tests
- Add notifications.store.test.js with 11 tests
- Web: 48, Mobile: 102, Workers: 342 tests

## Recent Changes (v2.0.0-rc.18)
- Add playbooks.store.test.js with 8 tests
- Web: 28, Mobile: 102, Workers: 342 tests

## Recent Changes (v2.0.0-rc.17)
- Add users.store.test.js with 10 tests
- Web: 20, Mobile: 102, Workers: 342 tests

## Recent Changes (v2.0.0-rc.16)
- Add web frontend testing with vitest
- Add auth.store.test.js with 5 tests
- Web: 5, Mobile: 102, Workers: 342 tests

## Recent Changes (v2.0.0-rc.14)
- Add tokens_provider_test.dart with model tests
- Test ApiToken with expiry checks
- Test Session model with device parsing
- Flutter: 102 tests, Workers: 342 tests

## Recent Changes (v2.0.0-rc.13)
- Add backup_provider_test.dart with model tests
- Test Backup and BackupStatus models
- Flutter: 92 tests, Workers: 342 tests

## Recent Changes (v2.0.0-rc.12)
- Add backup_provider.dart for Flutter mobile
- Add notifications_provider.dart with tests
- Web: Notifications page and updated Sidebar
- Flutter: 86 tests, Workers: 342 tests

## Recent Changes (v2.0.0-rc.11)
- Add Notifications page for Web Frontend
- Add notifications_provider for Flutter
- Add notifications_provider_test.dart
- Update Sidebar with Tasks, Schedules, Notifications
- Flutter: 86 tests, Workers: 342 tests

## Recent Changes (v2.0.0-rc.10)
- Add Tasks and Schedules pages to Web Frontend
- Implement full logs page with filtering and streaming
- Add NotificationsApi and BackupApi services
- Fix Schedule enabled parsing for bool/int types
- Add schedules provider tests
- Flutter: 81 tests, Workers: 342 tests

## Recent Changes (v2.0.0-rc.8)
- Implement full logs page with filtering and streaming
- Add NotificationsApi for push notifications
- Add BackupApi for backup management
- Add model classes for ApiToken, Session, Notification, Backup
- Flutter: 81 tests, Workers: 342 tests

## Recent Changes (v2.0.0-rc.7)
- Fix logs_provider_test.dart to match actual LogsState API
- Fix dashboard_page_test.dart async timer issues
- Add logs_provider.dart for log management
- Fix schedules_provider.dart import paths and add enabled parameter
- Simplify widget_test.dart to basic smoke test
- Flutter: 74 tests, Workers: 342 tests

## Recent Changes (v2.0.0-rc.6)
- Implement MFA (Two-Factor Authentication) settings for mobile app
- MFA setup dialog with QR code and recovery codes
- Enable/disable MFA with verification
- Update Android compileSdk to 36 for latest dependencies

## Recent Changes (v2.0.0-rc.5)
- Implement Tasks management feature for mobile app
- TasksApi with Dio client integration
- TasksPage with status filtering
- TaskDetailPage with execution logs
- Cancel and retry task operations
- Fix Android build: compileSdk 35, ndkVersion 27.0.12077973
- Fix keystore handling for missing secrets
- Update Flutter version to 3.41.5 in CI

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