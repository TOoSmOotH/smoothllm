# Changelog

All notable changes to SmoothWeb will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned Features
- Email verification
- Password reset flow
- Two-factor authentication (2FA)
- OAuth providers (Google, GitHub)
- PostgreSQL/MySQL migration scripts
- Rate limiting middleware

### Known Issues
- None currently

## [1.0.0] - TBD

### Added
- Initial release of SmoothWeb template
- JWT-based authentication with email/password
- RBAC with Admin and User roles
- First registered user becomes administrator
- Comprehensive user profiles with privacy controls
- Social links management
- Cyberpunk UI theme with neon effects
- Encrypted SQLite database (SQLCipher)
- Docker support for development and production
- Full test coverage (unit, integration, E2E)
- GitHub Actions CI/CD workflows
- Issue and PR templates
- Comprehensive documentation

### Backend Features
- Gin web framework
- GORM ORM with SQLite, PostgreSQL, MySQL support
- Casbin RBAC implementation
- JWT token generation and validation
- Password hashing with bcrypt
- File upload service (avatar, cover photo)
- API versioning (/api/v1/*)
- Health check endpoints
- CORS middleware
- Request logging
- Error handling middleware

### Frontend Features
- Vue 3 with Composition API
- TypeScript support
- Pinia state management
- Vue Router with route guards
- Tailwind CSS with cyberpunk theme
- shadcn-vue component library
- Custom cyberpunk components:
  - GlitchText
  - NeonBorder
  - CyberCard
  - CyberInput
  - HolographicCard
  - Scanlines
  - DataCard
  - TechLabel
- Axios client with interceptors
- Form validation with vee-validate
- Responsive design
- Hot module replacement (HMR)

### Documentation
- Implementation plan
- GitHub template guide
- Setup guide
- Customization guide
- Deployment guide
- API documentation
- Architecture overview
- Contributing guidelines
- Security policy
- Code of conduct

### Development Tools
- Multi-stage Docker builds
- Docker Compose configurations
- Air hot reload (backend)
- Vite dev server (frontend)
- Makefile for common tasks
- ESLint and Prettier (frontend)
- Go vet and fmt (backend)

### Testing
- Unit tests (testify)
- Integration tests (httpexpect)
- Component tests (Vue Test Utils)
- E2E tests (Playwright)
- Coverage reporting

### Security
- AES-256 encrypted database
- JWT with refresh tokens
- Password hashing (bcrypt, cost 12)
- RBAC (Admin/User)
- Input validation
- SQL injection prevention (GORM)
- XSS prevention (Vue)
- CORS protection

## [Unreleased Template Updates]

### Documentation
- GitHub template setup guide
- Issue templates (bug report, feature request)
- Pull request template
- CONTRIBUTING.md guidelines
- SECURITY.md policy
- Code of Conduct

### GitHub Configuration
- GitHub Actions workflows (CI, template checks)
- Dependabot configuration
- CODEOWNERS file
- .gitattributes for templates
- .editorconfig for consistency

### Developer Experience
- Setup scripts (Bash, Python, Node.js)
- Placeholder replacement scripts
- Centralized branding configuration
- Environment variable templates

---

## Versioning Scheme

SmoothWeb follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backwards compatible)
- **PATCH**: Bug fixes (backwards compatible)

Example: `1.0.0` → `1.1.0` (new feature) → `1.1.1` (bug fix) → `2.0.0` (breaking change)

## Release Notes

Release notes will include:

- New features
- Breaking changes
- Bug fixes
- Security updates
- Migration notes (if needed)
- Upgrade instructions

## How to Update

As a template user, to update your project with changes from the template:

```bash
# Add template as remote
git remote add template https://github.com/USERNAME/smoothweb.git

# Fetch latest changes
git fetch template main

# Merge changes
git merge template/main

# Resolve any conflicts
# (resolve conflicts manually)

# Test thoroughly
make test
make dev

# Commit and push
git add .
git commit -m "chore: update from template v1.0.0"
git push origin main
```

## Archive

### Previous Versions

(No previous versions yet)

---

**For detailed commit history, see [GitHub Commits](https://github.com/USERNAME/smoothweb/commits/main)**
