package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/smoothweb/backend/internal/custom/models"
)

const (
	// DefaultUserAgent is used when the client doesn't provide a User-Agent
	DefaultUserAgent = "SmoothLLM-Proxy/1.0"

	// DefaultMaxTokens is the default max_tokens value for Anthropic if not provided
	DefaultMaxTokens = 4096

	// AnthropicVersion is the API version header required by Anthropic
	AnthropicVersion = "2023-06-01"
)

// ProxyService handles LLM request proxying with model routing and request transformation
type ProxyService struct {
	keyService      *KeyService
	providerService *ProviderService
	usageService    *UsageService
}

// NewProxyService creates a new ProxyService instance
func NewProxyService(keyService *KeyService, providerService *ProviderService, usageService *UsageService) *ProxyService {
	return &ProxyService{
		keyService:      keyService,
		providerService: providerService,
		usageService:    usageService,
	}
}

// OpenAIChatRequest represents an OpenAI-compatible chat completion request
type OpenAIChatRequest struct {
	Model            string                 `json:"model"`
	Messages         []OpenAIMessage        `json:"messages"`
	MaxTokens        *int                   `json:"max_tokens,omitempty"`
	Temperature      *float64               `json:"temperature,omitempty"`
	TopP             *float64               `json:"top_p,omitempty"`
	N                *int                   `json:"n,omitempty"`
	Stream           *bool                  `json:"stream,omitempty"`
	Stop             interface{}            `json:"stop,omitempty"`
	PresencePenalty  *float64               `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64               `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]float64     `json:"logit_bias,omitempty"`
	User             string                 `json:"user,omitempty"`
	Extra            map[string]interface{} `json:"-"` // Catch any additional fields
}

// OpenAIMessage represents a message in the OpenAI format
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

// AnthropicRequest represents an Anthropic API request
type AnthropicRequest struct {
	Model       string             `json:"model"`
	Messages    []AnthropicMessage `json:"messages"`
	MaxTokens   int                `json:"max_tokens"`
	System      string             `json:"system,omitempty"`
	Temperature *float64           `json:"temperature,omitempty"`
	TopP        *float64           `json:"top_p,omitempty"`
	TopK        *int               `json:"top_k,omitempty"`
	Stream      *bool              `json:"stream,omitempty"`
	StopSequences []string         `json:"stop_sequences,omitempty"`
	Metadata    map[string]string  `json:"metadata,omitempty"`
}

// AnthropicMessage represents a message in the Anthropic format
type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ModelInfo contains parsed model routing information
type ModelInfo struct {
	ProviderType string // openai, anthropic, local
	ModelName    string // The actual model name to forward to provider
	FullModel    string // Original model string from request
}

// ProxyResult contains the result of a proxy operation for usage tracking
type ProxyResult struct {
	StatusCode      int
	InputTokens     int
	OutputTokens    int
	TotalTokens     int
	RequestDuration time.Duration
	ErrorMessage    string
	Model           string
}

// ValidateAndGetProvider validates the API key and returns the associated provider
func (s *ProxyService) ValidateAndGetProvider(apiKey string) (*models.ProxyAPIKey, *models.Provider, error) {
	// Validate the proxy API key
	proxyKey, err := s.keyService.ValidateKey(apiKey)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid API key: %w", err)
	}

	// Get the associated provider
	provider, err := s.providerService.GetProviderByIDInternal(proxyKey.ProviderID)
	if err != nil {
		return nil, nil, fmt.Errorf("provider not found: %w", err)
	}

	// Check if provider is active
	if !provider.IsActive {
		return nil, nil, fmt.Errorf("provider is not active")
	}

	return proxyKey, provider, nil
}

// ParseModelName parses a LiteLLM-style model name (provider/model) into components
func (s *ProxyService) ParseModelName(model string, defaultProviderType string) *ModelInfo {
	info := &ModelInfo{
		FullModel: model,
	}

	// Check if model has provider prefix (e.g., "openai/gpt-4o", "anthropic/claude-sonnet-4")
	if parts := strings.SplitN(model, "/", 2); len(parts) == 2 {
		info.ProviderType = strings.ToLower(parts[0])
		info.ModelName = parts[1]
	} else {
		// No prefix, use default provider type and full model name
		info.ProviderType = defaultProviderType
		info.ModelName = model
	}

	return info
}

// ProxyRequest handles the proxying of an LLM request to the appropriate provider
func (s *ProxyService) ProxyRequest(c *gin.Context, proxyKey *models.ProxyAPIKey, provider *models.Provider) (*ProxyResult, error) {
	startTime := time.Now()
	result := &ProxyResult{}

	// Read the request body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		result.StatusCode = http.StatusBadRequest
		result.ErrorMessage = "failed to read request body"
		return result, fmt.Errorf("failed to read request body: %w", err)
	}

	// Parse the OpenAI-format request
	var chatReq OpenAIChatRequest
	if err := json.Unmarshal(bodyBytes, &chatReq); err != nil {
		result.StatusCode = http.StatusBadRequest
		result.ErrorMessage = "invalid request body"
		return result, fmt.Errorf("failed to parse request body: %w", err)
	}

	// Store the model for usage tracking
	result.Model = chatReq.Model

	// Parse the model name
	modelInfo := s.ParseModelName(chatReq.Model, provider.ProviderType)

	// Determine the target URL and transform request if needed
	var targetURL string
	var requestBody []byte

	switch provider.ProviderType {
	case models.ProviderTypeAnthropic:
		targetURL = strings.TrimSuffix(provider.GetBaseURL(), "/") + "/v1/messages"
		requestBody, err = s.transformToAnthropic(&chatReq, modelInfo.ModelName)
		if err != nil {
			result.StatusCode = http.StatusBadRequest
			result.ErrorMessage = fmt.Sprintf("failed to transform request: %v", err)
			return result, err
		}
	default:
		// OpenAI and local providers use the same format
		targetURL = strings.TrimSuffix(provider.GetBaseURL(), "/") + "/v1/chat/completions"
		// Update the model name in the request if it was prefixed
		chatReq.Model = modelInfo.ModelName
		requestBody, err = json.Marshal(chatReq)
		if err != nil {
			result.StatusCode = http.StatusInternalServerError
			result.ErrorMessage = "failed to marshal request"
			return result, fmt.Errorf("failed to marshal request: %w", err)
		}
	}

	// Create the proxy request
	proxyReq, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewReader(requestBody))
	if err != nil {
		result.StatusCode = http.StatusInternalServerError
		result.ErrorMessage = "failed to create proxy request"
		return result, fmt.Errorf("failed to create proxy request: %w", err)
	}

	// Copy relevant headers, preserving User-Agent
	s.copyHeaders(c.Request, proxyReq, provider)

	// Execute the proxy request
	client := &http.Client{
		Timeout: 5 * time.Minute, // Long timeout for LLM responses
	}

	resp, err := client.Do(proxyReq)
	if err != nil {
		result.StatusCode = http.StatusBadGateway
		result.ErrorMessage = fmt.Sprintf("proxy request failed: %v", err)
		result.RequestDuration = time.Since(startTime)
		return result, fmt.Errorf("proxy request failed: %w", err)
	}
	defer resp.Body.Close()

	// Record timing
	result.RequestDuration = time.Since(startTime)
	result.StatusCode = resp.StatusCode

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		result.ErrorMessage = "failed to read response"
		return result, fmt.Errorf("failed to read response body: %w", err)
	}

	// Extract usage information from response if available
	s.extractUsageFromResponse(respBody, provider.ProviderType, result)

	// Record usage asynchronously (non-blocking)
	s.recordUsage(proxyKey, provider, result)

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Write the response
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)

	return result, nil
}

// ProxyWithReverseProxy uses httputil.ReverseProxy for streaming support
func (s *ProxyService) ProxyWithReverseProxy(c *gin.Context, provider *models.Provider, modelInfo *ModelInfo) error {
	baseURL := provider.GetBaseURL()
	if baseURL == "" {
		return fmt.Errorf("no base URL configured for provider")
	}

	// Determine target path based on provider type
	var targetPath string
	switch provider.ProviderType {
	case models.ProviderTypeAnthropic:
		targetPath = "/v1/messages"
	default:
		targetPath = "/v1/chat/completions"
	}

	targetURL, err := url.Parse(baseURL + targetPath)
	if err != nil {
		return fmt.Errorf("failed to parse target URL: %w", err)
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(targetURL)
			r.Out.Host = targetURL.Host

			// Preserve client User-Agent (use slice assignment per Go docs)
			if userAgent := r.In.Header.Get("User-Agent"); userAgent != "" {
				r.Out.Header["User-Agent"] = []string{userAgent}
			} else {
				r.Out.Header["User-Agent"] = []string{DefaultUserAgent}
			}

			// Set provider-specific auth headers
			switch provider.ProviderType {
			case models.ProviderTypeAnthropic:
				r.Out.Header.Set("x-api-key", provider.APIKey)
				r.Out.Header.Set("anthropic-version", AnthropicVersion)
				// Remove Authorization header if present
				r.Out.Header.Del("Authorization")
			default:
				r.Out.Header.Set("Authorization", "Bearer "+provider.APIKey)
			}

			// Ensure Content-Type is set
			if r.Out.Header.Get("Content-Type") == "" {
				r.Out.Header.Set("Content-Type", "application/json")
			}
		},
	}

	proxy.ServeHTTP(c.Writer, c.Request)
	return nil
}

// transformToAnthropic transforms an OpenAI-format request to Anthropic format
func (s *ProxyService) transformToAnthropic(req *OpenAIChatRequest, modelName string) ([]byte, error) {
	anthropicReq := AnthropicRequest{
		Model: modelName,
	}

	// Handle max_tokens - required for Anthropic
	if req.MaxTokens != nil {
		anthropicReq.MaxTokens = *req.MaxTokens
	} else {
		anthropicReq.MaxTokens = DefaultMaxTokens
	}

	// Copy optional parameters
	if req.Temperature != nil {
		anthropicReq.Temperature = req.Temperature
	}
	if req.TopP != nil {
		anthropicReq.TopP = req.TopP
	}
	if req.Stream != nil {
		anthropicReq.Stream = req.Stream
	}

	// Handle stop sequences
	if req.Stop != nil {
		switch v := req.Stop.(type) {
		case string:
			anthropicReq.StopSequences = []string{v}
		case []interface{}:
			for _, item := range v {
				if s, ok := item.(string); ok {
					anthropicReq.StopSequences = append(anthropicReq.StopSequences, s)
				}
			}
		}
	}

	// Transform messages - extract system message to separate field
	for _, msg := range req.Messages {
		switch msg.Role {
		case "system":
			// Anthropic uses a separate system field, not a system message
			if anthropicReq.System == "" {
				anthropicReq.System = msg.Content
			} else {
				anthropicReq.System += "\n\n" + msg.Content
			}
		case "user", "assistant":
			anthropicReq.Messages = append(anthropicReq.Messages, AnthropicMessage{
				Role:    msg.Role,
				Content: msg.Content,
			})
		default:
			// Map other roles to user (e.g., "function" results)
			anthropicReq.Messages = append(anthropicReq.Messages, AnthropicMessage{
				Role:    "user",
				Content: msg.Content,
			})
		}
	}

	// Ensure we have at least one message
	if len(anthropicReq.Messages) == 0 {
		return nil, fmt.Errorf("at least one user or assistant message is required")
	}

	return json.Marshal(anthropicReq)
}

// copyHeaders copies relevant headers from the original request to the proxy request
func (s *ProxyService) copyHeaders(original *http.Request, proxy *http.Request, provider *models.Provider) {
	// Preserve User-Agent
	userAgent := original.Header.Get("User-Agent")
	if userAgent == "" {
		userAgent = DefaultUserAgent
	}
	proxy.Header.Set("User-Agent", userAgent)

	// Set Content-Type
	contentType := original.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}
	proxy.Header.Set("Content-Type", contentType)

	// Set Accept header
	if accept := original.Header.Get("Accept"); accept != "" {
		proxy.Header.Set("Accept", accept)
	}

	// Set provider-specific auth headers
	switch provider.ProviderType {
	case models.ProviderTypeAnthropic:
		proxy.Header.Set("x-api-key", provider.APIKey)
		proxy.Header.Set("anthropic-version", AnthropicVersion)
	default:
		proxy.Header.Set("Authorization", "Bearer "+provider.APIKey)
	}
}

// extractUsageFromResponse extracts token usage information from the LLM response
func (s *ProxyService) extractUsageFromResponse(body []byte, providerType string, result *ProxyResult) {
	// Only try to extract usage from successful responses
	if result.StatusCode < 200 || result.StatusCode >= 300 {
		return
	}

	switch providerType {
	case models.ProviderTypeAnthropic:
		s.extractAnthropicUsage(body, result)
	default:
		s.extractOpenAIUsage(body, result)
	}
}

// extractOpenAIUsage extracts usage from an OpenAI-format response
func (s *ProxyService) extractOpenAIUsage(body []byte, result *ProxyResult) {
	var resp struct {
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(body, &resp); err == nil {
		result.InputTokens = resp.Usage.PromptTokens
		result.OutputTokens = resp.Usage.CompletionTokens
		result.TotalTokens = resp.Usage.TotalTokens
	}
}

// extractAnthropicUsage extracts usage from an Anthropic response
func (s *ProxyService) extractAnthropicUsage(body []byte, result *ProxyResult) {
	var resp struct {
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(body, &resp); err == nil {
		result.InputTokens = resp.Usage.InputTokens
		result.OutputTokens = resp.Usage.OutputTokens
		result.TotalTokens = resp.Usage.InputTokens + resp.Usage.OutputTokens
	}
}

// ListModels returns a list of available models from the proxy key's provider
func (s *ProxyService) ListModels(provider *models.Provider) (interface{}, error) {
	// Build list of available models based on provider type
	type Model struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int64  `json:"created"`
		OwnedBy string `json:"owned_by"`
	}

	type ModelsResponse struct {
		Object string  `json:"object"`
		Data   []Model `json:"data"`
	}

	// Return a standardized list based on provider type
	// If the provider has a default model, include it in the list
	modelList := []Model{}
	now := time.Now().Unix()

	switch provider.ProviderType {
	case models.ProviderTypeOpenAI:
		modelList = append(modelList,
			Model{ID: "openai/gpt-4o", Object: "model", Created: now, OwnedBy: "openai"},
			Model{ID: "openai/gpt-4o-mini", Object: "model", Created: now, OwnedBy: "openai"},
			Model{ID: "openai/gpt-4-turbo", Object: "model", Created: now, OwnedBy: "openai"},
			Model{ID: "openai/gpt-3.5-turbo", Object: "model", Created: now, OwnedBy: "openai"},
		)
	case models.ProviderTypeAnthropic:
		modelList = append(modelList,
			Model{ID: "anthropic/claude-sonnet-4-20250514", Object: "model", Created: now, OwnedBy: "anthropic"},
			Model{ID: "anthropic/claude-opus-4-20250514", Object: "model", Created: now, OwnedBy: "anthropic"},
			Model{ID: "anthropic/claude-3-5-sonnet-20241022", Object: "model", Created: now, OwnedBy: "anthropic"},
			Model{ID: "anthropic/claude-3-5-haiku-20241022", Object: "model", Created: now, OwnedBy: "anthropic"},
		)
	case models.ProviderTypeLocal:
		// For local providers, just return the default model if configured
		if provider.DefaultModel != "" {
			modelList = append(modelList,
				Model{ID: "local/" + provider.DefaultModel, Object: "model", Created: now, OwnedBy: "local"},
			)
		}
	}

	// Add the provider's default model if configured and not already in list
	if provider.DefaultModel != "" {
		defaultModelID := provider.ProviderType + "/" + provider.DefaultModel
		found := false
		for _, m := range modelList {
			if m.ID == defaultModelID {
				found = true
				break
			}
		}
		if !found {
			modelList = append([]Model{{
				ID:      defaultModelID,
				Object:  "model",
				Created: now,
				OwnedBy: provider.ProviderType,
			}}, modelList...)
		}
	}

	return ModelsResponse{
		Object: "list",
		Data:   modelList,
	}, nil
}

// GetProxyKeyFromRequest extracts the proxy API key from the Authorization header
func (s *ProxyService) GetProxyKeyFromRequest(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	// Support "Bearer sk-smoothllm-xxx" format
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer "), nil
	}

	// Also support just the key directly
	if strings.HasPrefix(authHeader, models.ProxyAPIKeyPrefix) {
		return authHeader, nil
	}

	return "", fmt.Errorf("invalid Authorization header format")
}

// HandleProviderError returns appropriate error responses based on provider errors
func (s *ProxyService) HandleProviderError(statusCode int, errorMessage string) (int, map[string]interface{}) {
	switch statusCode {
	case http.StatusUnauthorized:
		return http.StatusUnauthorized, map[string]interface{}{
			"error": map[string]interface{}{
				"message": "Authentication failed with provider",
				"type":    "authentication_error",
				"code":    "provider_auth_error",
			},
		}
	case http.StatusTooManyRequests:
		return http.StatusTooManyRequests, map[string]interface{}{
			"error": map[string]interface{}{
				"message": "Rate limit exceeded at provider",
				"type":    "rate_limit_error",
				"code":    "provider_rate_limit",
			},
		}
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return http.StatusBadGateway, map[string]interface{}{
			"error": map[string]interface{}{
				"message": "Provider service unavailable",
				"type":    "server_error",
				"code":    "provider_unavailable",
			},
		}
	default:
		return statusCode, map[string]interface{}{
			"error": map[string]interface{}{
				"message": errorMessage,
				"type":    "api_error",
				"code":    "proxy_error",
			},
		}
	}
}

// recordUsage records API usage asynchronously using the UsageService
func (s *ProxyService) recordUsage(proxyKey *models.ProxyAPIKey, provider *models.Provider, result *ProxyResult) {
	if s.usageService == nil {
		return
	}

	req := &RecordUsageRequest{
		UserID:          proxyKey.UserID,
		ProxyKeyID:      proxyKey.ID,
		ProviderID:      provider.ID,
		Model:           result.Model,
		InputTokens:     result.InputTokens,
		OutputTokens:    result.OutputTokens,
		TotalTokens:     result.TotalTokens,
		RequestDuration: int(result.RequestDuration.Milliseconds()),
		StatusCode:      result.StatusCode,
		ErrorMessage:    result.ErrorMessage,
		InputCostPer1K:  provider.InputCostPer1K,
		OutputCostPer1K: provider.OutputCostPer1K,
	}

	s.usageService.RecordUsageAsync(req)
}
