# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## SDK v1.5.0 - 2026-03-17

### Added

#### API Integration
- REST API v2 support
- GraphQL endpoint (beta)
- Webhook management
- API key management
- Rate limit configuration

#### Automation
- Scheduled tasks
- Task queue management
- Task history and logs
- Retry configuration

### Improved
- API response caching
- Request batching

## [v0.9.9] - 2026-03-16

### Added

#### i18n & Localization
- Multi-language support with vue-i18n for Web GUI
- Translations for 5 languages: English, Simplified Chinese, Traditional Chinese, Japanese, Arabic
- RTL (Right-to-Left) language support for Arabic
- Automatic browser language detection
- Language switcher component
- Mobile app localization with flutter_localizations
- ARB localization files for all 5 languages

#### Accessibility
- Keyboard navigation composable for all components
- Focus trap for modal dialogs
- High contrast mode detection and support
- Reduced motion preference support
- Screen reader announcement utility
- Skip link support
- ARIA labels and semantic HTML

#### Security
- Rate limiting middleware (1000 requests/minute)
- Security headers middleware (X-Frame-Options, CSP, etc.)
- Input sanitization utilities
- Email and URL validation helpers

### Changed
- Web GUI version updated to v0.9.9
- Mobile App version updated to v0.9.9
- Desktop App version updated to v0.9.9
- SDK version updated to v1.5.0

## [v0.7.9] - 2026-03-16

### Added
- Desktop app Linux support with Wails v2
- Theme customization with 5 accent colors
- Compact mode and font size options
- Theme store with Pinia state management
- ThemeSwitcher component

## [v0.6.9] - 2026-03-16

### Added
- Mobile app store preparation
- AndroidManifest.xml with full permissions and features
- iOS Info.plist with Face ID, camera, notifications
- Firebase Cloud Messaging integration
- Deep link support

## [v0.5.9] - 2026-03-16

### Added
- Theme store with Pinia
- ThemeSwitcher component
- Dark/Light/System theme modes

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

[Unreleased]: https://github.com/AnixOps/anixops-control-center/compare/v0.9.9...HEAD
[v0.9.9]: https://github.com/AnixOps/anixops-control-center/compare/v0.7.9...v0.9.9
[v0.7.9]: https://github.com/AnixOps/anixops-control-center/compare/v0.6.9...v0.7.9
[v0.6.9]: https://github.com/AnixOps/anixops-control-center/compare/v0.5.9...v0.6.9
[v0.5.9]: https://github.com/AnixOps/anixops-control-center/compare/v0.1.0...v0.5.9
[v0.1.0]: https://github.com/AnixOps/anixops-control-center/releases/tag/v0.1.0