# Contributing to SmoothWeb

Thank you for your interest in contributing to SmoothWeb! This document provides guidelines and instructions for contributing.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How to Contribute](#how-to-contribute)
- [Development Setup](#development-setup)
- [Coding Standards](#coding-standards)
- [Commit Messages](#commit-messages)
- [Pull Request Process](#pull-request-process)
- [Testing Guidelines](#testing-guidelines)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Enhancements](#suggesting-enhancements)

## 🤝 Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment. Please read our full [Code of Conduct](CODE_OF_CONDUCT.md).

## 🚀 How to Contribute

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. If you find a bug:

1. Use the [Bug Report template](.github/ISSUE_TEMPLATE/bug_report.yml)
2. Provide as much detail as possible
3. Include steps to reproduce
4. Describe expected vs actual behavior
5. Include environment details

### Suggesting Enhancements

We welcome feature requests! Please:

1. Use the [Feature Request template](.github/ISSUE_TEMPLATE/feature_request.yml)
2. Clearly describe the problem it solves
3. Provide use cases
4. Consider if this feature benefits the broader community

### Submitting Pull Requests

We love pull requests! Here's how to get started:

## 💻 Development Setup

### Prerequisites

- Go 1.21+ 
- Node.js 18+ and npm
- Docker and Docker Compose
- Git

### Fork and Clone

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/smoothweb.git
   cd smoothweb
   ```

3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/USERNAME/smoothweb.git
   ```

### Install Dependencies

```bash
# Backend
cd backend
go mod download

# Frontend
cd ../frontend
npm install
```

### Start Development

```bash
# Start all services with Docker
make dev

# Or start individually:
cd backend && make run-dev    # Backend with hot reload
cd frontend && npm run dev     # Frontend dev server
```

## 📏 Coding Standards

### Backend (Go)

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Export functions and types with clear documentation
- Write meaningful commit messages
- Add comments for complex logic

```go
// Good
func GetUserByID(id uint) (*User, error) {
    // Get user from database
    var user User
    if err := db.First(&user, id).Error; err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    return &user, nil
}
```

### Frontend (Vue 3 + TypeScript)

- Use Composition API with `<script setup>`
- TypeScript strict mode
- Follow Vue 3 style guide
- Use descriptive component names
- Write comments for complex logic

```vue
<!-- Good -->
<script setup lang="ts">
import { ref, computed } from 'vue'

interface Props {
  title: string
  count: number
}

const props = defineProps<Props>()
const isCompleted = computed(() => props.count > 0)
</script>
```

## 📝 Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, no logic change)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Test additions or changes
- `chore`: Maintenance tasks

### Examples

```
feat(auth): add OAuth provider for GitHub

Add GitHub OAuth authentication with token exchange.
Update Casbin policy to allow OAuth endpoints.

Fixes #123
```

```
fix(profile): prevent SQL injection in profile search

Use GORM parameterized queries instead of string interpolation.

Fixes #456
```

## 🔄 Pull Request Process

### Before Submitting

1. **Create a branch** for your feature:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following coding standards

3. **Write tests** for new functionality

4. **Update documentation** if needed

5. **Commit your changes** with clear messages

6. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

7. **Create Pull Request**:
   - Use the [PR template](.github/pull_request_template.md)
   - Reference related issues
   - Describe changes clearly
   - Add screenshots for UI changes

### PR Review Process

1. Automated checks must pass (CI, tests, linting)
2. At least one maintainer review required
3. Address all review feedback
4. Squash commits when approved
5. Merge into main branch

## 🧪 Testing Guidelines

### Backend Tests

```bash
cd backend

# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific test
go test -run TestAuthService_Register -v
```

**Requirements:**
- Unit tests for all services
- Integration tests for API endpoints
- Minimum 85% code coverage
- Tests must be deterministic

### Frontend Tests

```bash
cd frontend

# Run unit tests
npm run test

# Run E2E tests
npm run test:e2e

# Run with coverage
npm run test:coverage
```

**Requirements:**
- Unit tests for components and composables
- Component tests for UI interactions
- E2E tests for critical user flows
- Minimum 80% code coverage

## 🐛 Reporting Bugs

Use the [Bug Report template](.github/ISSUE_TEMPLATE/bug_report.yml) and include:

- Clear description
- Steps to reproduce
- Expected vs actual behavior
- Environment details
- Screenshots/videos if applicable
- Logs and error messages

## ✨ Suggesting Enhancements

Use the [Feature Request template](.github/ISSUE_TEMPLATE/feature_request.yml) and include:

- Problem statement
- Proposed solution
- Use cases
- Alternative approaches considered
- Implementation suggestions (if technical)

## 📞 Getting Help

- Check [Documentation](docs/)
- Search [existing issues](https://github.com/USERNAME/smoothweb/issues)
- Start a [Discussion](https://github.com/USERNAME/smoothweb/discussions)
- Join our community chat (if applicable)

## 🎉 Recognition

Contributors will be:
- Listed in README.md
- Mentioned in release notes
- Eligible for contributor badges
- Invited to join maintainers team for significant contributions

Thank you for contributing to SmoothWeb! 🚀
