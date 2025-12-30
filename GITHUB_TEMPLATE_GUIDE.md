# SmoothWeb - GitHub Template Guide

## 📚 Overview

This guide outlines all the additional files, configurations, and best practices needed to make SmoothWeb a production-ready GitHub repository template.

---

## 🎯 Template Goals

1. **One-Click Setup**: Users can click "Use this template" and get a working project
2. **Easy Customization**: Clear guides for branding, features, and configuration
3. **Professional Documentation**: Comprehensive setup, deployment, and customization guides
4. **Active Maintenance**: GitHub Actions for template health and dependency updates
5. **Community Ready**: Issue templates, PR templates, and contribution guidelines

---

## 📁 Additional Project Structure (Template-Specific)

```
smoothweb/
├── .github/
│   ├── ISSUE_TEMPLATE/              # Issue form templates
│   │   ├── config.yml              # Issue template configuration
│   │   ├── bug_report.yml         # Bug report form
│   │   ├── feature_request.yml    # Feature request form
│   │   └── question.yml          # Question form
│   ├── PULL_REQUEST_TEMPLATE/     # PR templates (optional)
│   │   ├── pr-template.md
│   │   └── feature-pr.md
│   ├── workflows/                 # GitHub Actions workflows
│   │   ├── ci.yml               # Continuous integration
│   │   ├── template-check.yml     # Template health checks
│   │   ├── dependabot-automerge.yml
│   │   ├── notify-updates.yml    # Notify template users
│   │   └── deploy.yml           # Deployment workflow
│   ├── dependabot.yml            # Dependabot configuration
│   ├── CODEOWNERS                # Code ownership rules
│   └── FUNDING.yml              # Sponsorship configuration (optional)
│
├── scripts/                      # Setup and automation scripts
│   ├── setup.sh                  # Bash setup script
│   ├── setup.py                  # Python setup script
│   ├── setup.mjs                # Node.js setup script
│   └── replace-placeholders.sh   # Placeholder replacement
│
├── docs/                        # Additional documentation
│   ├── SETUP.md                 # Step-by-step setup guide
│   ├── CUSTOMIZATION.md          # Customization guide
│   ├── DEPLOYMENT.md            # Deployment instructions
│   ├── API.md                   # API documentation
│   ├── ARCHITECTURE.md          # Architecture overview
│   └── FAQ.md                  # Frequently asked questions
│
├── config/branding.js           # Centralized branding configuration
├── .gitattributes              # Git attributes for templates
├── .editorconfig               # Editor configuration
├── LICENSE                     # MIT License
├── README.md                   # Main documentation (template-focused)
├── CONTRIBUTING.md             # Contribution guidelines
├── CODE_OF_CONDUCT.md         # Code of conduct
├── SUPPORT.md                  # Support information
├── SECURITY.md                 # Security policy
├── CHANGELOG.md               # Changelog template
└── USAGE.md                   # Usage examples and patterns
```

---

## 📝 Template-Specific Files

### 1. Issue Templates

#### `.github/ISSUE_TEMPLATE/config.yml`
```yaml
blank_issues_enabled: true
contact_links:
  - name: 📖 Documentation
    url: https://github.com/username/smoothweb/blob/main/docs/README.md
    about: Please read the documentation first
  - name: 💬 Discussions
    url: https://github.com/username/smoothweb/discussions
    about: Ask questions and discuss features
```

#### `.github/ISSUE_TEMPLATE/bug_report.yml`
```yaml
name: 🐛 Bug Report
description: Report a bug to help us improve SmoothWeb
title: "[Bug]: "
labels: ["bug", "triage"]
assignees: []
type: bug
body:
  - type: markdown
    attributes:
      value: |
        Thanks for taking the time to fill out this bug report!
        Please ensure you've read the documentation first.

  - type: input
    id: version
    attributes:
      label: Version
      description: What version of SmoothWeb are you using?
      placeholder: e.g., v1.0.0
    validations:
      required: true

  - type: dropdown
    id: component
    attributes:
      label: Component
      description: Which component is affected?
      options:
        - Authentication
        - User Profiles
        - Admin Panel
        - API
        - Database
        - Docker
        - Frontend
        - Other
    validations:
      required: true

  - type: textarea
    id: description
    attributes:
      label: What happened?
      description: Describe the bug in detail. What did you expect to happen?
      placeholder: Tell us what you see!
    validations:
      required: true

  - type: textarea
    id: reproduction
    attributes:
      label: Steps to reproduce
      description: How can we reproduce this behavior?
      placeholder: |
        1. Go to '...'
        2. Click on '....'
        3. Scroll down to '....'
        4. See error
    validations:
      required: true

  - type: textarea
    id: environment
    attributes:
      label: Environment
      description: |
        Please provide your environment details.
      value: |
        - OS: [e.g., Ubuntu 20.04, macOS 14, Windows 11]
        - Docker: [version]
        - Go: [version]
        - Node: [version]
        - Browser: [e.g., Chrome 120, Firefox 121, Safari 17]
      render: markdown
    validations:
      required: false

  - type: textarea
    id: logs
    attributes:
      label: Logs/Errors
      description: Please paste any relevant logs or error messages here
      render: shell
```

#### `.github/ISSUE_TEMPLATE/feature_request.yml`
```yaml
name: ✨ Feature Request
description: Suggest a new feature or enhancement
title: "[Feature]: "
labels: ["enhancement", "triage"]
assignees: []
type: bug
body:
  - type: markdown
    attributes:
      value: |
        Thanks for suggesting a new feature!

  - type: textarea
    id: problem
    attributes:
      label: What problem does this solve?
      description: A clear and concise description of what the problem is
    validations:
      required: true

  - type: textarea
    id: solution
    attributes:
      label: What is the proposed solution?
      description: Describe your proposed solution in detail
    validations:
      required: true

  - type: textarea
    id: alternatives
    attributes:
      label: What alternatives have you considered?
      description: Describe any alternative solutions or features you've considered
    validations:
      required: false

  - type: checkboxes
    id: contribution
    attributes:
      label: Would you like to contribute?
      description: Check if you're willing to implement this feature
      options:
        - label: I'm willing to submit a pull request for this feature
          required: false
```

### 2. Pull Request Template

#### `.github/pull_request_template.md`
```markdown
## 📝 Description
Please include a summary of the changes and the related issue.

Fixes #(issue_number)

## 🔄 Type of Change
Please delete options that are not relevant.

- [ ] 🐛 Bug fix (non-breaking change which fixes an issue)
- [ ] ✨ New feature (non-breaking change which adds functionality)
- [ ] 💥 Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] 📚 Documentation update
- [ ] 🎨 Code style update (formatting, renaming)
- [ ] ♻️ Code refactoring (no functional change)
- [ ] ⚡ Performance improvement
- [ ] ✅ Test additions or updates
- [ ] 🌍 Internationalization
- [ ] 🔧 Configuration changes

## 🧪 How Has This Been Tested?

Please describe the tests that you ran to verify your changes.

- [ ] Unit tests pass: `make test-coverage` (backend) or `npm test` (frontend)
- [ ] Integration tests pass: `make test-integration`
- [ ] Manual testing: (describe what you tested)
- [ ] E2E tests pass: `npm run test:e2e`

## ✅ Checklist

- [ ] My code follows the style guidelines of this project
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] Any dependent changes have been merged and published in downstream modules
- [ ] I have updated the CHANGELOG.md with my changes

## 📸 Screenshots (if applicable)

Add screenshots to help explain your changes:

Before:
![Before](url)

After:
![After](url)
```

### 3. Documentation Files

#### `README.md` (Template-Focused)
```markdown
# ⚡ SmoothWeb

> A production-ready, full-stack web application template with a cyberpunk UI, encrypted database, and comprehensive user management

[![Build Status](https://github.com/username/smoothweb/workflows/CI/badge.svg)](https://github.com/username/smoothweb/actions)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![Vue Version](https://img.shields.io/badge/Vue-3.4+-4FC08D?logo=vue.js)](https://vuejs.org/)

## ✨ Features

- 🔐 **Secure Authentication**: JWT-based auth with email/password (OAuth-ready)
- 👤 **User Management**: RBAC with Admin/User roles
- 🎨 **Cyberpunk UI**: Beautiful neon aesthetic with glow effects
- 🔒 **Encrypted Database**: AES-256 encrypted SQLite (SQLCipher)
- 🐳 **Docker Ready**: Multi-container setup for development and production
- 📊 **Comprehensive Profiles**: Customizable user profiles with privacy controls
- 🧪 **Full Test Coverage**: Unit, integration, and E2E tests
- 📚 **Well Documented**: Extensive guides for setup, customization, and deployment

## 🚀 Quick Start

### Using as Template

1. Click the green **"Use this template"** button above
2. Name your new repository
3. Clone your new repository:
   ```bash
   git clone https://github.com/YOUR_USERNAME/YOUR_REPO.git
   cd YOUR_REPO
   ```

4. Install dependencies:
   ```bash
   # Backend
   cd backend && go mod download

   # Frontend
   cd frontend && npm install
   ```

5. Start development:
   ```bash
   make dev
   ```

6. Open your browser:
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080
   - API Health: http://localhost:8080/health

## 📚 Documentation

- **[SETUP Guide](docs/SETUP.md)** - Step-by-step setup instructions
- **[CUSTOMIZATION Guide](docs/CUSTOMIZATION.md)** - How to customize for your needs
- **[DEPLOYMENT Guide](docs/DEPLOYMENT.md)** - Deployment instructions
- **[CONTRIBUTING Guide](CONTRIBUTING.md)** - How to contribute
- **[API Documentation](docs/API.md)** - API endpoint documentation
- **[Architecture Overview](docs/ARCHITECTURE.md)** - System architecture details

## 🎨 Screenshots

<!-- Add screenshots of your application here -->

## 🔧 Technology Stack

### Backend
- **Framework**: [Gin](https://github.com/gin-gonic/gin) v1.11.0
- **ORM**: [GORM](https://github.com/go-gorm/gorm) v1.31.1
- **Database**: [SQLCipher](https://github.com/mutecomm/go-sqlcipher) (Encrypted SQLite)
- **Auth**: [golang-jwt/jwt](https://github.com/golang-jwt/jwt) v5
- **RBAC**: [Casbin](https://github.com/casbin/casbin) v3.4.1

### Frontend
- **Framework**: [Vue 3](https://vuejs.org/) with Composition API
- **Build Tool**: [Vite](https://vitejs.dev/) 5.0+
- **UI Library**: [shadcn-vue](https://www.shadcn-vue.com/) + Tailwind CSS
- **State**: [Pinia](https://pinia.vuejs.org/) 2.1+
- **Testing**: [Vitest](https://vitest.dev/) + Vue Test Utils

## 📦 Scripts

### Project-wide
```bash
make dev          # Start development environment (Docker)
make prod         # Start production environment
make build        # Build all services
make test         # Run all tests
make clean        # Clean build artifacts
```

### Backend
```bash
cd backend
make run-dev     # Run with hot reload (Air)
make test         # Run unit tests
make test-coverage # Run tests with coverage
make build        # Build binary
```

### Frontend
```bash
cd frontend
npm run dev      # Start dev server (Vite)
npm run build     # Build for production
npm run test      # Run unit tests
npm run test:e2e  # Run E2E tests
npm run lint      # Lint code
```

## 🐳 Docker

```bash
# Development (with hot reload)
docker-compose -f docker-compose.dev.yml up

# Production
docker-compose -f docker-compose.prod.yml up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f
```

## 🧪 Testing

```bash
# Run all tests
make test

# Backend tests only
cd backend && make test

# Frontend tests only
cd frontend && npm run test

# E2E tests
cd frontend && npm run test:e2e
```

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) first.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Gin](https://github.com/gin-gonic/gin) and [Vue 3](https://vuejs.org/)
- UI inspired by cyberpunk design principles
- Icons by [Lucide](https://lucide.dev/)
- Fonts: [Orbitron](https://fonts.google.com/specimen/Orbitron), [Exo 2](https://fonts.google.com/specimen/Exo+2)

## 📞 Support

- 📖 [Documentation](docs/)
- 💬 [GitHub Discussions](https://github.com/username/smoothweb/discussions)
- 🐛 [Issue Tracker](https://github.com/username/smoothweb/issues)

## ⭐ Show Your Support

If this template helps you, please consider giving it a ⭐ on GitHub!

---

**Built with ❤️ using Go, Vue 3, and Tailwind CSS**
