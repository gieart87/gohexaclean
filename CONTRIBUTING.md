# Contributing to GoHexaClean

Thank you for your interest in contributing to GoHexaClean! This document provides guidelines and instructions for contributing.

## Code of Conduct

Please be respectful and constructive in your interactions with other contributors.

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in Issues
2. If not, create a new issue with:
   - Clear title and description
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details (Go version, OS, etc.)

### Suggesting Features

1. Open an issue with the "enhancement" label
2. Describe the feature and its use case
3. Explain why it fits the project's goals

### Pull Requests

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Follow the architecture principles**
   - Maintain separation of concerns
   - Follow SOLID principles
   - Keep domain logic pure and framework-independent
   - Use dependency injection

4. **Write tests**
   - Unit tests for business logic
   - Integration tests for adapters
   - Maintain or improve code coverage

5. **Follow Go best practices**
   - Run `make fmt` before committing
   - Run `make lint` to check code quality
   - Ensure `make test` passes

6. **Commit your changes**
   ```bash
   git commit -m "feat: add amazing feature"
   ```

   Use conventional commit messages:
   - `feat:` for new features
   - `fix:` for bug fixes
   - `docs:` for documentation
   - `refactor:` for code refactoring
   - `test:` for adding tests
   - `chore:` for maintenance tasks

7. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```

8. **Create a Pull Request**
   - Provide a clear description
   - Reference related issues
   - Include screenshots if applicable

## Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/gohexaclean.git

# Add upstream remote
git remote add upstream https://github.com/gieart87/gohexaclean.git

# Install dependencies
make deps

# Generate proto files
make proto

# Run tests
make test
```

## Architecture Guidelines

### Domain Layer
- Pure business logic, no framework dependencies
- Domain entities should be self-validating
- Use value objects for domain concepts

### Application Layer
- Orchestrate use cases
- Coordinate between domain and ports
- Handle transaction boundaries

### Ports
- Define clear interfaces
- Keep interfaces small and focused
- Use dependency inversion

### Adapters
- Implement port interfaces
- Handle framework-specific details
- Keep adapters thin

## Testing Guidelines

- **Unit Tests**: Test business logic in isolation
- **Integration Tests**: Test adapters with real dependencies
- **E2E Tests**: Test complete flows

```bash
# Run unit tests
make test-unit

# Run integration tests
make test-integration

# Generate coverage report
make test-coverage
```

## Documentation

- Add comments for exported functions
- Update README.md if adding new features
- Document architectural decisions in docs/

## Questions?

Feel free to open an issue for any questions or clarifications.

Thank you for contributing! 🎉
