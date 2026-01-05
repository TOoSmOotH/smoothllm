package services

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/smoothweb/backend/internal/custom/models"
)

const (
	// Anthropic OAuth configuration
	// Using Claude Code CLI's client ID for OAuth (this is public knowledge)
	AnthropicClientID     = "9d1c250a-e61b-44d9-88ed-5944d1962f5e"
	AnthropicAuthURL      = "https://claude.ai/oauth/authorize"
	AnthropicTokenURL     = "https://console.anthropic.com/api/oauth/token"
	AnthropicScope        = "user:inference user:profile"
	AnthropicResponseType = "code"

	// Token expiry buffer - refresh tokens 5 minutes before expiry
	TokenExpiryBuffer = 5 * time.Minute

	// Access token default lifetime (8 hours)
	DefaultTokenLifetime = 8 * time.Hour
)

// OAuthService handles OAuth flows for Claude Max subscriptions
type OAuthService struct {
	db              *gorm.DB
	providerService *ProviderService
	frontendURL     string

	// In-memory storage for PKCE state (code verifiers mapped by state)
	// In production, consider using Redis or database storage
	stateMu     sync.RWMutex
	stateStore  map[string]*OAuthState
}

// OAuthState holds the PKCE parameters for an OAuth flow
type OAuthState struct {
	CodeVerifier string
	ProviderID   uint
	UserID       uint
	CreatedAt    time.Time
}

// TokenResponse represents the response from the OAuth token endpoint
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // seconds
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// NewOAuthService creates a new OAuthService instance
func NewOAuthService(db *gorm.DB, providerService *ProviderService, frontendURL string) *OAuthService {
	svc := &OAuthService{
		db:              db,
		providerService: providerService,
		frontendURL:     frontendURL,
		stateStore:      make(map[string]*OAuthState),
	}

	// Start a goroutine to clean up expired states
	go svc.cleanupExpiredStates()

	return svc
}

// GenerateAuthURL generates the OAuth authorization URL with PKCE
func (s *OAuthService) GenerateAuthURL(userID, providerID uint, redirectURI string) (string, error) {
	// Verify the provider exists and belongs to the user
	provider, err := s.providerService.getProviderByID(userID, providerID)
	if err != nil {
		return "", fmt.Errorf("provider not found: %w", err)
	}

	if provider.ProviderType != models.ProviderTypeAnthropicMax {
		return "", fmt.Errorf("provider is not an Anthropic Max provider")
	}

	// Generate PKCE parameters
	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		return "", fmt.Errorf("failed to generate code verifier: %w", err)
	}

	codeChallenge := generateCodeChallenge(codeVerifier)
	state := generateState()

	// Store the state and code verifier
	s.stateMu.Lock()
	s.stateStore[state] = &OAuthState{
		CodeVerifier: codeVerifier,
		ProviderID:   providerID,
		UserID:       userID,
		CreatedAt:    time.Now(),
	}
	s.stateMu.Unlock()

	// Build the authorization URL
	params := url.Values{}
	params.Set("client_id", AnthropicClientID)
	params.Set("response_type", AnthropicResponseType)
	params.Set("redirect_uri", redirectURI)
	params.Set("scope", AnthropicScope)
	params.Set("state", state)
	params.Set("code_challenge", codeChallenge)
	params.Set("code_challenge_method", "S256")

	authURL := fmt.Sprintf("%s?%s", AnthropicAuthURL, params.Encode())
	return authURL, nil
}

// ExchangeCode exchanges the authorization code for tokens
func (s *OAuthService) ExchangeCode(code, state, redirectURI string) (*models.Provider, error) {
	// Retrieve and validate the state
	s.stateMu.Lock()
	oauthState, exists := s.stateStore[state]
	if exists {
		delete(s.stateStore, state)
	}
	s.stateMu.Unlock()

	if !exists {
		return nil, fmt.Errorf("invalid or expired state")
	}

	// Check if state is too old (15 minutes max)
	if time.Since(oauthState.CreatedAt) > 15*time.Minute {
		return nil, fmt.Errorf("state expired")
	}

	// Exchange the code for tokens
	tokens, err := s.exchangeCodeForTokens(code, oauthState.CodeVerifier, redirectURI)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Calculate token expiry time
	expiresAt := time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)

	// Update the provider with the tokens
	provider, err := s.providerService.getProviderByID(oauthState.UserID, oauthState.ProviderID)
	if err != nil {
		return nil, fmt.Errorf("provider not found: %w", err)
	}

	updates := map[string]interface{}{
		"access_token":    tokens.AccessToken,
		"refresh_token":   tokens.RefreshToken,
		"token_expires_at": expiresAt,
		"oauth_connected": true,
	}

	if err := s.db.Model(provider).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update provider tokens: %w", err)
	}

	// Refresh the provider data
	if err := s.db.First(provider, provider.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to refresh provider: %w", err)
	}

	return provider, nil
}

// RefreshAccessToken refreshes the access token using the refresh token
func (s *OAuthService) RefreshAccessToken(provider *models.Provider) error {
	if provider.RefreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	// Make the token refresh request
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", provider.RefreshToken)
	data.Set("client_id", AnthropicClientID)

	req, err := http.NewRequest("POST", AnthropicTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token refresh failed: %s", string(body))
	}

	var tokens TokenResponse
	if err := json.Unmarshal(body, &tokens); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	// Calculate new expiry time
	expiresAt := time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)

	// Update the provider with new tokens
	updates := map[string]interface{}{
		"access_token":     tokens.AccessToken,
		"token_expires_at": expiresAt,
	}

	// If a new refresh token was provided, update it too
	if tokens.RefreshToken != "" {
		updates["refresh_token"] = tokens.RefreshToken
	}

	if err := s.db.Model(provider).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update provider tokens: %w", err)
	}

	// Update the in-memory provider object
	provider.AccessToken = tokens.AccessToken
	provider.TokenExpiresAt = &expiresAt
	if tokens.RefreshToken != "" {
		provider.RefreshToken = tokens.RefreshToken
	}

	return nil
}

// EnsureValidToken ensures the provider has a valid access token, refreshing if needed
func (s *OAuthService) EnsureValidToken(provider *models.Provider) error {
	if !provider.IsOAuthProvider() {
		return nil // Not an OAuth provider, nothing to do
	}

	if !provider.NeedsTokenRefresh() {
		return nil // Token is still valid
	}

	return s.RefreshAccessToken(provider)
}

// DisconnectOAuth disconnects OAuth from a provider
func (s *OAuthService) DisconnectOAuth(userID, providerID uint) error {
	provider, err := s.providerService.getProviderByID(userID, providerID)
	if err != nil {
		return err
	}

	updates := map[string]interface{}{
		"access_token":     "",
		"refresh_token":    "",
		"token_expires_at": nil,
		"oauth_connected":  false,
	}

	return s.db.Model(provider).Updates(updates).Error
}

// exchangeCodeForTokens exchanges the authorization code for access and refresh tokens
func (s *OAuthService) exchangeCodeForTokens(code, codeVerifier, redirectURI string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("client_id", AnthropicClientID)
	data.Set("code_verifier", codeVerifier)

	req, err := http.NewRequest("POST", AnthropicTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed: %s", string(body))
	}

	var tokens TokenResponse
	if err := json.Unmarshal(body, &tokens); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokens, nil
}

// cleanupExpiredStates periodically cleans up expired OAuth states
func (s *OAuthService) cleanupExpiredStates() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.stateMu.Lock()
		now := time.Now()
		for state, oauthState := range s.stateStore {
			if now.Sub(oauthState.CreatedAt) > 15*time.Minute {
				delete(s.stateStore, state)
			}
		}
		s.stateMu.Unlock()
	}
}

// generateCodeVerifier generates a cryptographically random code verifier for PKCE
func generateCodeVerifier() (string, error) {
	// Generate 32 random bytes (256 bits)
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// Encode as base64url without padding
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// generateCodeChallenge generates the code challenge from the code verifier using S256
func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// generateState generates a cryptographically random state parameter
func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// GetProviderByID retrieves a provider by ID for a specific user (delegates to providerService)
func (s *OAuthService) GetProviderByID(userID, providerID uint) (*models.Provider, error) {
	return s.providerService.getProviderByID(userID, providerID)
}

// TestOAuthConnection tests if the OAuth connection is working
func (s *OAuthService) TestOAuthConnection(provider *models.Provider) error {
	if !provider.OAuthConnected {
		return fmt.Errorf("OAuth is not connected")
	}

	// Ensure we have a valid token
	if err := s.EnsureValidToken(provider); err != nil {
		return fmt.Errorf("failed to ensure valid token: %w", err)
	}

	// Make a test request to the Anthropic API
	baseURL := provider.GetBaseURL()
	testURL := strings.TrimSuffix(baseURL, "/") + "/v1/messages"

	// Create a minimal test request
	testBody := map[string]interface{}{
		"model":      "claude-3-5-haiku-20241022",
		"max_tokens": 1,
		"messages": []map[string]string{
			{"role": "user", "content": "Hi"},
		},
	}
	bodyBytes, _ := json.Marshal(testBody)

	req, err := http.NewRequest("POST", testURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+provider.AccessToken)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	// 200 means success, 400 means auth worked but request was invalid (acceptable)
	// 401/403 means auth failed
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("authentication failed: OAuth token is invalid or expired")
	}

	if resp.StatusCode >= 500 {
		return fmt.Errorf("provider server error: status %d", resp.StatusCode)
	}

	return nil
}
