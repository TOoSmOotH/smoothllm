# Specification: Smooth LLM Proxy

## Overview

Build **Smooth LLM Proxy**, a LiteLLM-inspired API key management and LLM usage tracking system on top of the smoothweb template. This application allows users to configure their LLM provider API keys (OpenAI, Anthropic, local models), create proxy API keys, and have all requests proxied transparently to the configured providers while preserving client user-agents. The system tracks detailed usage statistics and costs (user-defined) for each API key and provider.

## Workflow Type

**Type**: feature

**Rationale**: This is a substantial new feature implementation that adds LLM proxy capabilities, new database models, API endpoints, and frontend views to an existing template. It requires creating multiple new components while leveraging the existing auth, RBAC, and user management infrastructure.

## Task Scope

### Services Involved
- **backend** (primary) - Go/Gin API server: proxy endpoints, key management, usage tracking, provider routing
- **frontend** (integration) - Vue 3 UI: key management dashboard, provider configuration, usage statistics views

### This Task Will:
- [ ] Create database models for LLM providers, API keys, and usage tracking
- [ ] Implement proxy endpoints that route requests to configured LLM providers
- [ ] Build API key management endpoints (CRUD for provider keys and proxy keys)
- [ ] Preserve client User-Agent headers when proxying requests
- [ ] Implement usage tracking and cost calculation based on user-defined rates
- [ ] Create frontend views for provider configuration and key management
- [ ] Build usage statistics dashboard
- [ ] Follow LiteLLM model naming conventions for compatibility

### Out of Scope:
- OAuth/social login for LLM providers (API key only)
- Automatic cost detection from providers (user specifies costs)
- Load balancing across multiple provider keys
- Request caching or response caching
- Streaming response handling (initial implementation uses standard request/response)
- Rate limiting (can be added later)

## Service Context

### Backend

**Tech Stack:**
- Language: Go
- Framework: Gin
- Database: SQLite with GORM
- Authentication: JWT
- Authorization: Casbin RBAC

**Entry Point:** `backend/cmd/api/main.go`

**How to Run:**
```bash
cd backend && go run cmd/api/main.go
```

**Port:** 8080

**Key Directories:**
- `backend/internal/custom/` - Custom routes and handlers (PRIMARY modification area)
- `backend/internal/models/` - GORM database models
- `backend/internal/handlers/` - HTTP handlers
- `backend/internal/services/` - Business logic
- `backend/internal/auth/` - JWT authentication

### Frontend

**Tech Stack:**
- Language: TypeScript
- Framework: Vue 3
- State Management: Pinia
- Styling: Tailwind CSS
- Build Tool: Vite
- Testing: Vitest / Playwright

**Entry Point:** `frontend/src/main.ts`

**How to Run:**
```bash
cd frontend && npm run dev
```

**Port:** 5173

**Key Directories:**
- `frontend/src/custom/` - Custom routes and config (PRIMARY modification area)
- `frontend/src/views/` - Page components
- `frontend/src/components/` - Reusable UI components
- `frontend/src/stores/` - Pinia stores

## Files to Modify

| File | Service | What to Change |
|------|---------|---------------|
| `backend/internal/custom/routes.go` | backend | Register all LLM proxy and management API routes, call custom migrations |
| `frontend/src/custom/appConfig.ts` | frontend | Update branding to "Smooth LLM Proxy", add navigation items |
| `frontend/src/custom/routes.ts` | frontend | Add routes for provider config, key management, usage stats pages |

**Note:** Do NOT modify `backend/internal/database/migrations.go` directly. Instead, create `backend/internal/custom/migrations.go` and call it from `RegisterRoutes()` to run custom model migrations. This avoids merge conflicts with template updates.

## Files to Create

**Note:** Create subdirectories under `backend/internal/custom/` to isolate custom code. This keeps customizations separate from core template code per CUSTOMIZATION.md guidelines. Import paths will be `github.com/smoothweb/backend/internal/custom/models`, etc.

| File | Service | Purpose |
|------|---------|---------|
| `backend/internal/custom/models/provider.go` | backend | LLM provider configuration model |
| `backend/internal/custom/models/api_key.go` | backend | Proxy API key model |
| `backend/internal/custom/models/usage.go` | backend | Usage tracking model |
| `backend/internal/custom/handlers/proxy.go` | backend | LLM proxy request handler |
| `backend/internal/custom/handlers/providers.go` | backend | Provider CRUD handlers |
| `backend/internal/custom/handlers/keys.go` | backend | API key management handlers |
| `backend/internal/custom/handlers/usage.go` | backend | Usage statistics handlers |
| `backend/internal/custom/services/proxy.go` | backend | Proxy request forwarding service |
| `backend/internal/custom/services/providers.go` | backend | Provider management service |
| `backend/internal/custom/services/keys.go` | backend | API key service |
| `backend/internal/custom/services/usage.go` | backend | Usage tracking service |
| `backend/internal/custom/migrations.go` | backend | Custom model migrations (called from routes.go) |
| `frontend/src/views/Providers.vue` | frontend | LLM provider configuration page |
| `frontend/src/views/ApiKeys.vue` | frontend | Proxy API key management page |
| `frontend/src/views/Usage.vue` | frontend | Usage statistics dashboard |
| `frontend/src/stores/providers.ts` | frontend | Provider state management |
| `frontend/src/stores/apiKeys.ts` | frontend | API key state management |
| `frontend/src/stores/usage.ts` | frontend | Usage data state management |

## Files to Reference

These files show patterns to follow:

| File | Pattern to Copy |
|------|----------------|
| `backend/internal/custom/routes.go` | How to register custom routes with deps (DB, JWT, RBAC) |
| `backend/internal/handlers/auth.go` | Handler struct pattern with service injection |
| `backend/internal/models/user.go` | GORM model structure with relationships |
| `backend/internal/auth/middleware.go` | JWT auth middleware pattern, getting user context |
| `docs/CUSTOMIZATION.md` | Template customization guidelines |
| `frontend/src/views/Dashboard.vue` | Vue 3 page component pattern |
| `frontend/src/config/appConfig.ts` | App configuration structure |

## Patterns to Follow

### Custom Route Registration Pattern

From `backend/internal/custom/routes.go`:

```go
func RegisterRoutes(v1 *gin.RouterGroup, deps Dependencies) {
    // Protected routes example
    providers := v1.Group("/providers")
    providers.Use(auth.AuthMiddleware(deps.JWT))
    {
        providers.GET("", handleListProviders)
        providers.POST("", handleCreateProvider)
    }
}
```

**Key Points:**
- Use `deps.JWT` for auth middleware on protected routes
- Use `deps.DB` for database access
- Use `deps.RBAC` for role-based access control

### Handler Pattern

From `backend/internal/handlers/auth.go`:

```go
type AuthHandler struct {
    authService *auth.Service
}

func NewAuthHandler(authService *auth.Service) *AuthHandler {
    return &AuthHandler{authService: authService}
}

func (h *AuthHandler) SomeEndpoint(c *gin.Context) {
    // Get user from context
    userID := auth.GetUserID(c)

    // Parse request body
    var req models.SomeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Call service
    result, err := h.someService.DoSomething(&req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, result)
}
```

### GORM Model Pattern

From `backend/internal/models/user.go`:

```go
type Provider struct {
    gorm.Model

    UserID       uint   `gorm:"not null;index" json:"user_id"`
    Name         string `gorm:"type:varchar(100);not null" json:"name"`
    ProviderType string `gorm:"type:varchar(50);not null" json:"provider_type"`
    BaseURL      string `gorm:"type:varchar(500)" json:"base_url"`
    APIKey       string `gorm:"type:varchar(500);not null" json:"-"` // Never expose

    // Relationships
    User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
```

### Reverse Proxy Pattern

For LLM request proxying:

```go
func (s *ProxyService) ProxyRequest(c *gin.Context, provider *models.Provider) error {
    targetURL, _ := url.Parse(provider.BaseURL)

    proxy := &httputil.ReverseProxy{
        Rewrite: func(r *httputil.ProxyRequest) {
            r.SetURL(targetURL)
            r.Out.Host = targetURL.Host

            // Preserve client User-Agent (use slice assignment per Go docs)
            r.Out.Header["User-Agent"] = r.In.Header["User-Agent"]

            // Set provider auth
            if provider.ProviderType == "anthropic" {
                r.Out.Header.Set("x-api-key", provider.APIKey)
                r.Out.Header.Set("anthropic-version", "2023-06-01")
            } else {
                r.Out.Header.Set("Authorization", "Bearer "+provider.APIKey)
            }
        },
    }

    proxy.ServeHTTP(c.Writer, c.Request)
    return nil
}
```

### Provider API Differences

**CRITICAL**: OpenAI and Anthropic APIs have key differences that the proxy must handle:

| Aspect | OpenAI | Anthropic |
|--------|--------|-----------|
| Endpoint | `/v1/chat/completions` | `/v1/messages` |
| Auth Header | `Authorization: Bearer KEY` | `x-api-key: KEY` |
| Required Headers | `Content-Type` | `Content-Type`, `anthropic-version` |
| `max_tokens` | Optional | **REQUIRED** - proxy must ensure this is set |
| Message Format | `role: system/user/assistant` | `role: user/assistant` (system in separate field) |

**Implementation Strategy:**
1. Accept requests in OpenAI format on `/v1/chat/completions`
2. For Anthropic providers, transform request:
   - Change endpoint to `/v1/messages`
   - Set `x-api-key` instead of `Authorization`
   - Add `anthropic-version: 2023-06-01` header
   - Ensure `max_tokens` is present (use default like 4096 if missing)
   - Transform message format if needed (move system message to `system` field)

## Requirements

### Functional Requirements

1. **Provider Configuration**
   - Description: Users can add/edit/delete LLM provider configurations (OpenAI, Anthropic, local endpoints)
   - Acceptance: CRUD operations work, API keys are stored securely, provider connection can be tested

2. **Proxy API Key Management**
   - Description: Users create proxy API keys (sk-smoothllm-xxx format) that route to configured providers
   - Acceptance: Keys can be created, listed, revoked; each key maps to a specific provider configuration

3. **Request Proxying**
   - Description: Proxy receives requests, routes to correct provider based on model name, preserves User-Agent
   - Acceptance: OpenAI-compatible endpoint at `/v1/chat/completions`, requests forwarded with client User-Agent

4. **Usage Tracking**
   - Description: Track token usage, request counts, and calculate costs based on user-defined rates
   - Acceptance: Usage visible per key, per provider, per day; costs calculated correctly

5. **Model Routing**
   - Description: Route requests based on LiteLLM-style model names (e.g., `openai/gpt-4o`, `anthropic/claude-sonnet-4`)
   - Acceptance: Model prefix determines provider, model name forwarded correctly to provider

6. **Cost Configuration**
   - Description: Users specify input/output token costs per model
   - Acceptance: Costs stored per model, spend calculated from usage * cost

### Edge Cases

1. **Invalid API Key** - Return 401 with clear error message
2. **Provider Timeout** - Return 504 with timeout details, log the failure
3. **Missing Model Config** - Return 400 explaining which model was requested but not configured
4. **Rate Limit from Provider** - Forward provider's rate limit response, log for tracking
5. **Malformed Request** - Validate request body before proxying, return helpful errors
6. **User-Agent Missing** - Use a sensible default (e.g., "SmoothLLM-Proxy/1.0")
7. **Anthropic Missing max_tokens** - If routing to Anthropic and max_tokens not provided, inject default (4096)

## Implementation Notes

### DO
- Follow the customization pattern in `docs/CUSTOMIZATION.md` - keep changes in `/custom/` directories
- Reuse existing `auth.GetUserID(c)` for getting current user in handlers
- Use `deps.JWT` middleware for protected routes
- Store API keys encrypted or hashed when possible
- Use transactions when updating usage + calculating costs
- Follow existing handler/service separation pattern

### DON'T
- Don't modify core template files unless absolutely necessary
- Don't expose provider API keys in any API response
- Don't block on usage tracking - do it async if possible
- Don't hardcode provider endpoints - make them configurable
- Don't create new auth mechanisms - use existing JWT system

## Data Models

### Provider
```go
type Provider struct {
    gorm.Model
    UserID           uint    `gorm:"not null;index"`
    Name             string  `gorm:"type:varchar(100);not null"`
    ProviderType     string  `gorm:"type:varchar(50);not null"` // openai, anthropic, local
    BaseURL          string  `gorm:"type:varchar(500)"`
    APIKey           string  `gorm:"type:varchar(500);not null"` // encrypted
    IsActive         bool    `gorm:"default:true"`
    DefaultModel     string  `gorm:"type:varchar(100)"`
    InputCostPer1K   float64 `gorm:"default:0"`
    OutputCostPer1K  float64 `gorm:"default:0"`
}
```

### ProxyAPIKey
```go
type ProxyAPIKey struct {
    gorm.Model
    UserID       uint      `gorm:"not null;index"`
    ProviderID   uint      `gorm:"not null;index"`
    KeyHash      string    `gorm:"type:varchar(255);uniqueIndex;not null"`
    KeyPrefix    string    `gorm:"type:varchar(20);not null"` // sk-smoothllm-xxx (visible part)
    Name         string    `gorm:"type:varchar(100)"`
    IsActive     bool      `gorm:"default:true"`
    LastUsedAt   *time.Time
    ExpiresAt    *time.Time
}
```

### UsageRecord
```go
type UsageRecord struct {
    gorm.Model
    UserID          uint    `gorm:"not null;index"`
    ProxyKeyID      uint    `gorm:"not null;index"`
    ProviderID      uint    `gorm:"not null;index"`
    Model           string  `gorm:"type:varchar(100);index"`
    InputTokens     int     `gorm:"default:0"`
    OutputTokens    int     `gorm:"default:0"`
    TotalTokens     int     `gorm:"default:0"`
    Cost            float64 `gorm:"default:0"`
    RequestDuration int     // milliseconds
    StatusCode      int
    ErrorMessage    string  `gorm:"type:text"`
}
```

## API Endpoints

### Provider Management
- `GET /api/v1/providers` - List user's providers
- `POST /api/v1/providers` - Create provider
- `GET /api/v1/providers/:id` - Get provider details
- `PUT /api/v1/providers/:id` - Update provider
- `DELETE /api/v1/providers/:id` - Delete provider
- `POST /api/v1/providers/:id/test` - Test provider connection

### API Key Management
- `GET /api/v1/keys` - List user's proxy API keys
- `POST /api/v1/keys` - Create new proxy key (returns full key once)
- `GET /api/v1/keys/:id` - Get key details (not the key itself)
- `PUT /api/v1/keys/:id` - Update key (name, active status)
- `DELETE /api/v1/keys/:id` - Revoke/delete key

### Proxy Endpoints
- `POST /v1/chat/completions` - OpenAI-compatible chat endpoint (authenticated via proxy key)
- `GET /v1/models` - List available models

### Usage Statistics
- `GET /api/v1/usage` - Get usage summary
- `GET /api/v1/usage/daily` - Get daily breakdown
- `GET /api/v1/usage/by-key` - Get usage by proxy key
- `GET /api/v1/usage/by-provider` - Get usage by provider

## Development Environment

### Start Services

```bash
# Terminal 1: Backend
cd backend && go run cmd/api/main.go

# Terminal 2: Frontend
cd frontend && npm run dev
```

### Service URLs
- Backend API: http://localhost:8080
- Frontend: http://localhost:5173

### Required Environment Variables
- `SERVER_PORT`: Backend port (default: 8080)
- `GIN_MODE`: Gin mode (debug/release)
- `DB_PATH`: SQLite database path
- `JWT_SECRET`: Secret for JWT signing (sensitive)
- `JWT_EXPIRATION`: Token expiration duration
- `CORS_ORIGINS`: Allowed CORS origins
- `LOG_LEVEL`: Logging level
- `VITE_API_URL`: Frontend API URL (http://localhost:8080)

## Success Criteria

The task is complete when:

1. [ ] Users can add LLM provider configurations with API keys
2. [ ] Users can create proxy API keys mapped to providers
3. [ ] Proxy endpoint accepts requests and forwards to correct provider
4. [ ] Client User-Agent is preserved in proxied requests
5. [ ] Usage is tracked per request (tokens, cost, duration)
6. [ ] Usage statistics are visible in frontend dashboard
7. [ ] Model naming follows LiteLLM convention (provider/model)
8. [ ] No console errors in browser or server
9. [ ] Existing template tests still pass
10. [ ] New functionality verified via browser and API testing

## QA Acceptance Criteria

**CRITICAL**: These criteria must be verified by the QA Agent before sign-off.

### Unit Tests

| Test | File | What to Verify |
|------|------|----------------|
| Provider CRUD | `backend/internal/custom/services/providers_test.go` | Create, read, update, delete operations work |
| API Key Generation | `backend/internal/custom/services/keys_test.go` | Keys are generated with proper format, hashed correctly |
| Usage Calculation | `backend/internal/custom/services/usage_test.go` | Cost calculation from tokens and rates is accurate |
| Proxy Routing | `backend/internal/custom/services/proxy_test.go` | Correct provider selected based on model name |

### Integration Tests

| Test | Services | What to Verify |
|------|----------|----------------|
| Provider Creation Flow | backend API | POST /providers creates provider, GET /providers lists it |
| Key Management Flow | backend API | POST /keys returns full key once, GET /keys shows prefix only |
| Auth Integration | backend auth + custom | Custom routes properly protected by JWT |
| Proxy Authentication | backend proxy | Requests with valid proxy key succeed, invalid fails |

### End-to-End Tests

| Flow | Steps | Expected Outcome |
|------|-------|------------------|
| Provider Setup | 1. Login 2. Navigate to Providers 3. Add OpenAI config 4. Test connection | Provider saved, test succeeds |
| Key Creation | 1. Login 2. Navigate to API Keys 3. Create key for provider 4. Copy key | Key displayed once, usable for proxy |
| Proxy Request | 1. Send request to /v1/chat/completions with proxy key 2. Check response | Request proxied, response returned |
| Usage Tracking | 1. Make several proxy requests 2. Check usage dashboard | Usage shown with token counts and costs |

### Browser Verification (if frontend)

| Page/Component | URL | Checks |
|----------------|-----|--------|
| Home Page | `http://localhost:5173/` | Branding shows "Smooth LLM Proxy" |
| Providers Page | `http://localhost:5173/providers` | Can add/edit/delete providers |
| API Keys Page | `http://localhost:5173/keys` | Can create/revoke keys |
| Usage Dashboard | `http://localhost:5173/usage` | Shows usage statistics and costs |
| Navigation | All pages | Sidebar includes new navigation items |

### Database Verification

| Check | Query/Command | Expected |
|-------|---------------|----------|
| Tables exist | `.schema providers` | Provider, ProxyAPIKey, UsageRecord tables created |
| Foreign keys | Check constraints | Proper cascading deletes configured |
| Indexes | `.indices providers` | UserID indexed for performance |

### API Verification

| Endpoint | Method | Expected |
|----------|--------|----------|
| `/api/v1/providers` | GET | Returns user's providers only (user isolation) |
| `/v1/chat/completions` | POST | Returns proxied response or helpful error |
| `/v1/models` | GET | Returns list of configured models |

### QA Sign-off Requirements
- [ ] All unit tests pass
- [ ] All integration tests pass
- [ ] All E2E tests pass
- [ ] Browser verification complete
- [ ] Database state verified
- [ ] No regressions in existing functionality
- [ ] Code follows established patterns (custom directories, handler/service separation)
- [ ] No security vulnerabilities introduced (API keys not exposed, proper auth)
- [ ] Provider API keys never returned in any API response
- [ ] User-Agent forwarding verified in proxy requests
