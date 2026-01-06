package services

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/smoothweb/backend/internal/custom/models"
)

// ProviderService handles provider CRUD operations
type ProviderService struct {
	db           *gorm.DB
	oauthService *OAuthService
}

// NewProviderService creates a new ProviderService instance
func NewProviderService(db *gorm.DB) *ProviderService {
	return &ProviderService{db: db}
}

// SetOAuthService sets the OAuth service (to avoid circular dependency)
func (s *ProviderService) SetOAuthService(oauthService *OAuthService) {
	s.oauthService = oauthService
}

// ProviderResponse represents the provider data returned to clients
// Note: APIKey is never included in responses
type ProviderResponse struct {
	ID                   uint      `json:"id"`
	UserID               uint      `json:"user_id"`
	Name                 string    `json:"name"`
	ProviderType         string    `json:"provider_type"`
	BaseURL              string    `json:"base_url"`
	IsActive             bool      `json:"is_active"`
	DefaultModel         string    `json:"default_model"`
	InputCostPerMillion  float64   `json:"input_cost_per_million"`
	OutputCostPerMillion float64   `json:"output_cost_per_million"`
	OAuthConnected       bool      `json:"oauth_connected"` // Whether OAuth is connected (for anthropic_max)
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// CreateProviderRequest represents the request to create a provider
type CreateProviderRequest struct {
	Name                 string  `json:"name" binding:"required"`
	ProviderType         string  `json:"provider_type" binding:"required"`
	BaseURL              string  `json:"base_url"`
	APIKey               string  `json:"api_key"` // Not required for OAuth providers (validated in validateCreateRequest)
	IsActive             *bool   `json:"is_active"`
	DefaultModel         string  `json:"default_model"`
	InputCostPerMillion  float64 `json:"input_cost_per_million"`
	OutputCostPerMillion float64 `json:"output_cost_per_million"`
}

// UpdateProviderRequest represents the request to update a provider
type UpdateProviderRequest struct {
	Name                 *string  `json:"name,omitempty"`
	ProviderType         *string  `json:"provider_type,omitempty"`
	BaseURL              *string  `json:"base_url,omitempty"`
	APIKey               *string  `json:"api_key,omitempty"`
	IsActive             *bool    `json:"is_active,omitempty"`
	DefaultModel         *string  `json:"default_model,omitempty"`
	InputCostPerMillion  *float64 `json:"input_cost_per_million,omitempty"`
	OutputCostPerMillion *float64 `json:"output_cost_per_million,omitempty"`
}

// ListProviders returns all providers for a user
func (s *ProviderService) ListProviders(userID uint) ([]ProviderResponse, error) {
	var providers []models.Provider
	if err := s.db.Where("user_id = ?", userID).Find(&providers).Error; err != nil {
		return nil, fmt.Errorf("failed to list providers: %w", err)
	}

	responses := make([]ProviderResponse, len(providers))
	for i, provider := range providers {
		responses[i] = s.buildProviderResponse(&provider)
	}

	return responses, nil
}

// GetProvider retrieves a provider by ID for a specific user
func (s *ProviderService) GetProvider(userID, providerID uint) (*ProviderResponse, error) {
	provider, err := s.getProviderByID(userID, providerID)
	if err != nil {
		return nil, err
	}

	response := s.buildProviderResponse(provider)
	return &response, nil
}

// CreateProvider creates a new provider for a user
func (s *ProviderService) CreateProvider(userID uint, req *CreateProviderRequest) (*ProviderResponse, error) {
	// Validate input
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Set default IsActive value if not provided
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	provider := models.Provider{
		UserID:               userID,
		Name:                 req.Name,
		ProviderType:         req.ProviderType,
		BaseURL:              req.BaseURL,
		APIKey:               req.APIKey,
		IsActive:             isActive,
		DefaultModel:         req.DefaultModel,
		InputCostPerMillion:  req.InputCostPerMillion,
		OutputCostPerMillion: req.OutputCostPerMillion,
	}

	// For anthropic_max, the API key is actually a refresh token
	if req.ProviderType == models.ProviderTypeAnthropicMax && req.APIKey != "" {
		provider.RefreshToken = req.APIKey
		provider.APIKey = "" // Don't store refresh token as API key

		// Try to get an access token using the refresh token
		if s.oauthService != nil {
			// Create provider first so we can refresh the token
			if err := s.db.Create(&provider).Error; err != nil {
				return nil, fmt.Errorf("failed to create provider: %w", err)
			}

			// Refresh the token to validate it and get an access token
			if err := s.oauthService.RefreshAccessToken(&provider); err != nil {
				// Delete the provider if token refresh fails
				s.db.Delete(&provider)
				return nil, fmt.Errorf("invalid refresh token: %w", err)
			}

			// Mark as connected
			provider.OAuthConnected = true
			s.db.Save(&provider)

			response := s.buildProviderResponse(&provider)
			return &response, nil
		}
	}

	if err := s.db.Create(&provider).Error; err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	response := s.buildProviderResponse(&provider)
	return &response, nil
}

// UpdateProvider updates an existing provider
func (s *ProviderService) UpdateProvider(userID, providerID uint, req *UpdateProviderRequest) (*ProviderResponse, error) {
	provider, err := s.getProviderByID(userID, providerID)
	if err != nil {
		return nil, err
	}

	// Validate update request
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.ProviderType != nil {
		updates["provider_type"] = *req.ProviderType
	}
	if req.BaseURL != nil {
		updates["base_url"] = *req.BaseURL
	}
	if req.APIKey != nil {
		// For Claude Max providers, api_key is actually the refresh token
		providerType := provider.ProviderType
		if req.ProviderType != nil {
			providerType = *req.ProviderType
		}
		if providerType == models.ProviderTypeAnthropicMax {
			updates["refresh_token"] = *req.APIKey
			// Clear the api_key field
			updates["api_key"] = ""
		} else {
			updates["api_key"] = *req.APIKey
		}
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.DefaultModel != nil {
		updates["default_model"] = *req.DefaultModel
	}
	if req.InputCostPerMillion != nil {
		updates["input_cost_per_million"] = *req.InputCostPerMillion
	}
	if req.OutputCostPerMillion != nil {
		updates["output_cost_per_million"] = *req.OutputCostPerMillion
	}

	if len(updates) > 0 {
		if err := s.db.Model(provider).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update provider: %w", err)
		}
	}

	// Refresh provider data
	if err := s.db.First(provider, providerID).Error; err != nil {
		return nil, fmt.Errorf("failed to refresh provider: %w", err)
	}

	// For Claude Max providers with a new refresh token, validate it
	if req.APIKey != nil && provider.ProviderType == models.ProviderTypeAnthropicMax && s.oauthService != nil {
		if err := s.oauthService.RefreshAccessToken(provider); err != nil {
			// Revert the update if token refresh fails
			updates["refresh_token"] = ""
			updates["oauth_connected"] = false
			s.db.Model(provider).Updates(updates)
			return nil, fmt.Errorf("invalid refresh token: %w", err)
		}
		// Mark as connected
		provider.OAuthConnected = true
		s.db.Model(provider).Update("oauth_connected", true)
	}

	response := s.buildProviderResponse(provider)
	return &response, nil
}

// DeleteProvider deletes a provider for a user
func (s *ProviderService) DeleteProvider(userID, providerID uint) error {
	provider, err := s.getProviderByID(userID, providerID)
	if err != nil {
		return err
	}

	if err := s.db.Delete(provider).Error; err != nil {
		return fmt.Errorf("failed to delete provider: %w", err)
	}

	return nil
}

// TestConnection tests the connection to a provider
func (s *ProviderService) TestConnection(userID, providerID uint) error {
	provider, err := s.getProviderByID(userID, providerID)
	if err != nil {
		return err
	}

	return s.testProviderConnection(provider)
}

// TestConnectionWithRequest tests connection with provided credentials (before saving)
func (s *ProviderService) TestConnectionWithRequest(req *CreateProviderRequest) error {
	provider := &models.Provider{
		ProviderType: req.ProviderType,
		BaseURL:      req.BaseURL,
		APIKey:       req.APIKey,
	}

	return s.testProviderConnection(provider)
}

// GetProviderByIDInternal retrieves a provider by ID (for internal use by other services)
// Returns the full provider model including the API key
func (s *ProviderService) GetProviderByIDInternal(providerID uint) (*models.Provider, error) {
	var provider models.Provider
	if err := s.db.First(&provider, providerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("provider not found")
		}
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	return &provider, nil
}

// getProviderByID retrieves a provider ensuring it belongs to the user
func (s *ProviderService) getProviderByID(userID, providerID uint) (*models.Provider, error) {
	var provider models.Provider
	if err := s.db.Where("id = ? AND user_id = ?", providerID, userID).First(&provider).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("provider not found")
		}
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	return &provider, nil
}

// buildProviderResponse creates a ProviderResponse from a Provider model
// Note: APIKey is never included in the response
func (s *ProviderService) buildProviderResponse(provider *models.Provider) ProviderResponse {
	return ProviderResponse{
		ID:                   provider.ID,
		UserID:               provider.UserID,
		Name:                 provider.Name,
		ProviderType:         provider.ProviderType,
		BaseURL:              provider.GetBaseURL(),
		IsActive:             provider.IsActive,
		DefaultModel:         provider.DefaultModel,
		InputCostPerMillion:  provider.InputCostPerMillion,
		OutputCostPerMillion: provider.OutputCostPerMillion,
		OAuthConnected:       provider.OAuthConnected,
		CreatedAt:            provider.CreatedAt,
		UpdatedAt:            provider.UpdatedAt,
	}
}

// validateCreateRequest validates the create provider request
func (s *ProviderService) validateCreateRequest(req *CreateProviderRequest) error {
	// Validate name
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if len(req.Name) > 100 {
		return fmt.Errorf("name must be 100 characters or less")
	}

	// Validate provider type
	if err := s.validateProviderType(req.ProviderType); err != nil {
		return err
	}

	// Validate base URL if provided
	if req.BaseURL != "" {
		if err := s.validateBaseURL(req.BaseURL); err != nil {
			return err
		}
	}

	// Validate API key / refresh token
	if strings.TrimSpace(req.APIKey) == "" {
		if req.ProviderType == models.ProviderTypeAnthropicMax {
			return fmt.Errorf("refresh_token is required for Claude Max providers")
		}
		return fmt.Errorf("api_key is required")
	}

	// Validate cost values
	if req.InputCostPerMillion < 0 {
		return fmt.Errorf("input_cost_per_million cannot be negative")
	}
	if req.OutputCostPerMillion < 0 {
		return fmt.Errorf("output_cost_per_million cannot be negative")
	}

	return nil
}

// validateUpdateRequest validates the update provider request
func (s *ProviderService) validateUpdateRequest(req *UpdateProviderRequest) error {
	// Validate name if provided
	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			return fmt.Errorf("name cannot be empty")
		}
		if len(*req.Name) > 100 {
			return fmt.Errorf("name must be 100 characters or less")
		}
	}

	// Validate provider type if provided
	if req.ProviderType != nil {
		if err := s.validateProviderType(*req.ProviderType); err != nil {
			return err
		}
	}

	// Validate base URL if provided
	if req.BaseURL != nil && *req.BaseURL != "" {
		if err := s.validateBaseURL(*req.BaseURL); err != nil {
			return err
		}
	}

	// Validate API key if provided
	if req.APIKey != nil && strings.TrimSpace(*req.APIKey) == "" {
		return fmt.Errorf("api_key cannot be empty")
	}

	// Validate cost values if provided
	if req.InputCostPerMillion != nil && *req.InputCostPerMillion < 0 {
		return fmt.Errorf("input_cost_per_million cannot be negative")
	}
	if req.OutputCostPerMillion != nil && *req.OutputCostPerMillion < 0 {
		return fmt.Errorf("output_cost_per_million cannot be negative")
	}

	return nil
}

// validateProviderType validates the provider type
func (s *ProviderService) validateProviderType(providerType string) error {
	validTypes := []string{
		models.ProviderTypeOpenAI,
		models.ProviderTypeAnthropic,
		models.ProviderTypeAnthropicMax,
		models.ProviderTypeVLLM,
		models.ProviderTypeLocal,
		models.ProviderTypeZai,
	}
	for _, vt := range validTypes {
		if providerType == vt {
			return nil
		}
	}
	return fmt.Errorf("invalid provider_type: must be one of %v", validTypes)
}

// validateBaseURL validates a base URL
func (s *ProviderService) validateBaseURL(baseURL string) error {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("invalid base_url: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("base_url must include scheme and host")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("base_url scheme must be http or https")
	}
	return nil
}

// testProviderConnection tests connectivity to a provider
func (s *ProviderService) testProviderConnection(provider *models.Provider) error {
	// For OAuth providers, skip this test - they use OAuth test instead
	if provider.ProviderType == models.ProviderTypeAnthropicMax {
		if !provider.OAuthConnected {
			return fmt.Errorf("OAuth not connected - please connect via OAuth first")
		}
		return nil // OAuth connection test is handled separately
	}

	baseURL := provider.GetBaseURL()
	if baseURL == "" {
		return fmt.Errorf("no base URL configured for provider")
	}

	// Build test endpoint based on provider type
	var testURL string
	switch provider.ProviderType {
	case models.ProviderTypeOpenAI:
		testURL = strings.TrimSuffix(baseURL, "/") + "/v1/models"
	case models.ProviderTypeAnthropic:
		// Anthropic doesn't have a simple models endpoint, we'll check the base URL is reachable
		testURL = strings.TrimSuffix(baseURL, "/") + "/v1/messages"
	case models.ProviderTypeVLLM:
		// vLLM uses OpenAI-compatible API
		testURL = strings.TrimSuffix(baseURL, "/") + "/v1/models"
	case models.ProviderTypeLocal:
		testURL = strings.TrimSuffix(baseURL, "/") + "/v1/models"
	case models.ProviderTypeZai:
		testURL = strings.TrimSuffix(baseURL, "/") + "/models"
	default:
		testURL = strings.TrimSuffix(baseURL, "/") + "/v1/models"
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", testURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set appropriate auth headers based on provider type
	switch provider.ProviderType {
	case models.ProviderTypeAnthropic:
		req.Header.Set("x-api-key", provider.APIKey)
		req.Header.Set("anthropic-version", "2023-06-01")
		req.Header.Set("Content-Type", "application/json")
		// For Anthropic, we need to make a POST request to test, but for simplicity
		// we'll just verify the API key format and base URL reachability
		req.Method = "POST"
		// We can't actually test without sending a valid request body,
		// so we'll accept 400 (bad request) as a success indicator that auth worked
	default:
		req.Header.Set("Authorization", "Bearer "+provider.APIKey)
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	switch provider.ProviderType {
	case models.ProviderTypeAnthropic:
		// For Anthropic, 400 means auth worked but request was invalid (expected)
		// 401/403 means auth failed
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			return fmt.Errorf("authentication failed: invalid API key")
		}
		// 400 or 200 range are acceptable
		if resp.StatusCode >= 500 {
			return fmt.Errorf("provider server error: status %d", resp.StatusCode)
		}
	default:
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			return fmt.Errorf("authentication failed: invalid API key")
		}
		if resp.StatusCode >= 400 {
			return fmt.Errorf("connection test failed: status %d", resp.StatusCode)
		}
	}

	return nil
}
