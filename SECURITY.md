# Security Policy

## 🔒 Reporting Security Vulnerabilities

We take the security of SmoothWeb seriously. If you discover a security vulnerability, please report it responsibly.

### Reporting a Vulnerability

**Do NOT** open a public issue for security vulnerabilities.

Instead, send an email to: **security@example.com** (replace with actual email)

Include the following information:

- Description of the vulnerability
- Steps to reproduce (if applicable)
- Affected versions
- Potential impact
- Proof of concept (optional)

### What Happens Next

1. **Acknowledgment**: We will acknowledge receipt within 48 hours
2. **Investigation**: We will investigate and confirm the vulnerability
3. **Fix**: We will develop a fix
4. **Release**: We will release a security update
5. **Credit**: We will credit you in the release notes (if desired)

### Vulnerability Response Time

| Severity | Response Time |
|-----------|---------------|
| Critical | 48 hours |
| High | 72 hours |
| Medium | 1 week |
| Low | 2 weeks |

## 🛡️ Security Best Practices

### For Template Users

When using SmoothWeb as a template:

1. **Change Encryption Keys**
   ```bash
   # backend/.env
   DB_ENCRYPTION_KEY=generate-unique-key-here
   JWT_SECRET=generate-unique-jwt-secret-here
   ```

2. **Update Secrets** in production:
   - GitHub Secrets (Actions)
   - Environment variables
   - Docker secrets

3. **Enable HTTPS** in production

4. **Review Dependencies** regularly:
   ```bash
   # Backend
   cd backend && go get -u ./...
   
   # Frontend
   cd frontend && npm audit
   ```

5. **Set Up RBAC** properly:
   - First user becomes admin (by design)
   - Create additional admin accounts carefully
   - Review user permissions regularly

6. **Database Backups**:
   - Regular encrypted backups
   - Store backups securely
   - Test restore procedures

### For Contributors

When contributing:

1. **Never commit secrets** or credentials
2. **Validate inputs** on all endpoints
3. **Use parameterized queries** (GORM handles this)
4. **Sanitize outputs** to prevent XSS
5. **Implement rate limiting** for public endpoints
6. **Add security tests** for new features

## 🔐 Security Features

### Implemented

- ✅ **Password Hashing**: bcrypt with cost 12
- ✅ **JWT Authentication**: Stateless tokens with expiration
- ✅ **Encrypted Database**: AES-256 encrypted SQLite (SQLCipher)
- ✅ **RBAC**: Role-based access control with Casbin
- ✅ **CORS Protection**: Configurable origins
- ✅ **Input Validation**: Strict validation on all inputs
- ✅ **SQL Injection Prevention**: GORM parameterized queries
- ✅ **XSS Prevention**: Vue's automatic escaping

### Recommended Enhancements

- 🔲 Rate limiting middleware
- 🔲 CSRF tokens
- 🔲 Security headers middleware
- 🔲 Two-factor authentication (2FA)
- 🔲 Email verification
- 🔲 Password reset flow
- 🔲 Security audit logging
- 🔲 Intrusion detection

## 📊 Dependencies

We regularly update dependencies to address security vulnerabilities:

### Automated Updates

- **Dependabot**: Automated dependency updates
- **GitHub Actions**: Security scanning workflows

### Manual Updates

Before major releases, we:
1. Review security advisories
2. Update affected dependencies
3. Test thoroughly
4. Release security patches

## 🌐 HTTPS Configuration

### Production Deployment

Always use HTTPS in production:

#### Nginx Example

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
}
```

#### Docker Compose Example

```yaml
services:
  frontend:
    environment:
      - FORCE_HTTPS=true
```

## 📋 Security Checklist

Before deploying to production:

- [ ] Changed all default passwords and keys
- [ ] Enabled HTTPS
- [ ] Configured CORS for production domain
- [ ] Set up firewall rules
- [ ] Enabled database encryption
- [ ] Configured JWT secret
- [ ] Set up rate limiting
- [ ] Tested authentication and authorization
- [ ] Configured logging and monitoring
- [ ] Set up backup strategy
- [ ] Updated dependencies
- [ ] Reviewed and tested RBAC policies

## 🔍 Security Audits

We encourage security audits of SmoothWeb. If you conduct an audit:

1. Contact us first: security@example.com
2. Scope your testing to your own instance
3. Do not exploit vulnerabilities for gain
4. Report findings responsibly
5. Allow time for fixes before disclosure

## 📞 Security Questions

For security-related questions not involving vulnerabilities:

- 📧 Email: security@example.com
- 💬 GitHub Discussions: Use "Security" tag
- 🐛 Issues: Use "security" label

## 📜 Security Version History

| Version | Date | Description |
|---------|-------|-------------|
| 1.0.0 | TBD | Initial release with security features |
| - | - | - |

---

**Last Updated**: December 2025
