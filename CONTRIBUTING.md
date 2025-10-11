# Contributing to Go API Starter

Thank you for your interest in contributing to Go API Starter! This document provides guidelines and information for contributors.

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code.

## How to Contribute

### Reporting Issues

- Check existing issues before creating a new one
- Use clear, descriptive titles
- Include steps to reproduce the issue
- Provide system information and Go version

### Suggesting Features

- Open an issue with the "enhancement" label
- Clearly describe the feature and its benefits
- Provide examples of how it would be used

### Development Setup

1. Fork the repository
2. Clone your fork:

   ```bash
   git clone https://github.com/your-username/go-api-starter.git
   cd go-api-starter
   ```

3. Install dependencies:

   ```bash
   make deps
   make install-tools
   ```

4. Create a feature branch:

   ```bash
   git checkout -b feature/your-feature-name
   ```

5. Make your changes and test them:

   ```bash
   make test
   make lint
   ```

### Pull Request Process

1. Update documentation for any new features
2. Add tests for new functionality
3. Ensure all tests pass
4. Update CHANGELOG.md if applicable
5. Follow the commit message conventions
6. Create a pull request with:
   - Clear title and description
   - Reference any related issues
   - Include screenshots for UI changes

### Coding Standards

- Follow Go conventions and best practices
- Use `gofmt` and `goimports` for formatting
- Write clear, self-documenting code
- Add comments for complex logic
- Follow the existing project structure

### Testing

- Write unit tests for new functions
- Ensure test coverage doesn't decrease
- Test your changes thoroughly
- Include integration tests where appropriate

### Commit Messages

Follow conventional commits format:

- `feat: add new authentication method`
- `fix: resolve database connection issue`
- `docs: update API documentation`
- `refactor: improve error handling`
- `test: add unit tests for user service`

## Development Guidelines

### Project Structure

Please maintain the existing project structure:

- `cmd/` - Application entry points
- `internals/` - Private application code
- `pkg/` - Public library code
- `scripts/` - Build and deployment scripts

### Dependencies

- Prefer standard library when possible
- Keep dependencies minimal and well-maintained
- Update go.mod and go.sum when adding dependencies

### Security

- Never commit sensitive information
- Follow security best practices
- Report security issues privately

## Getting Help

- Check the documentation first
- Look through existing issues
- Join discussions in issues and PRs
- Contact maintainers if needed

Thank you for contributing!
