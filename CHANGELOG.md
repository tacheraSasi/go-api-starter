# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Three-service architecture: `core` (API), `minion` (asynq worker), and `overseer` (cron scheduler)
- Dedicated Dockerfiles and Compose stacks for each service
- `internals/minions` package for task types and handlers
- Service logger and asynq logger adapter in `pkg/logger`
- `DATABASE_URL`, `DB_SSL_MODE`, and `TIMEZONE` configuration
- Initial Go API starter template
- JWT authentication system
- Multi-database support (SQLite, PostgreSQL, MySQL)
- Structured logging with Logrus
- Docker and Docker Compose configuration
- CORS middleware
- Request validation
- Health check endpoints
- Hot reload development setup
- Comprehensive documentation
- Contributing guidelines
- MIT License

### Security
- Password hashing with bcrypt
- JWT token management
- Token blacklist functionality

## [1.0.0] - 2025-10-11

### Added
- Initial release of Go API Starter
- Complete REST API foundation
- Authentication and authorization
- Database integration with GORM
- Clean architecture implementation
- Development and production Docker setups
- Comprehensive README and documentation