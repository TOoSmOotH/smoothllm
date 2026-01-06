package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/smoothweb/backend/internal/custom/models"
)

// setupProxyTestDB creates an in-memory SQLite database for testing
func setupProxyTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Provider{}, &models.ProxyAPIKey{}, &models.UsageRecord{}, &models.KeyAllowedProvider{})
	require.NoError(t, err)

	return db
}

func createProxyTestServices(t *testing.T, db *gorm.DB) *ProxyService {
	keyService := NewKeyService(db)
	providerService := NewProviderService(db)
	usageService := NewUsageService(db)
	return NewProxyService(keyService, providerService, usageService, nil)
}

func TestProxyService_ParseModelName(t *testing.T) {
	db := setupProxyTestDB(t)
	service := createProxyTestServices(t, db)

	t.Run("parses LiteLLM-style openai/model format", func(t *testing.T) {
		info := service.ParseModelName("openai/gpt-4o", models.ProviderTypeOpenAI)

		assert.Equal(t, "openai", info.ProviderType)
		assert.Equal(t, "gpt-4o", info.ModelName)
		assert.Equal(t, "openai/gpt-4o", info.FullModel)
	})

	t.Run("parses LiteLLM-style anthropic/model format", func(t *testing.T) {
		info := service.ParseModelName("anthropic/claude-sonnet-4", models.ProviderTypeOpenAI)

		assert.Equal(t, "anthropic", info.ProviderType)
		assert.Equal(t, "claude-sonnet-4", info.ModelName)
		assert.Equal(t, "anthropic/claude-sonnet-4", info.FullModel)
	})

	t.Run("handles model without prefix using default provider", func(t *testing.T) {
		info := service.ParseModelName("gpt-4o", models.ProviderTypeOpenAI)

		assert.Equal(t, models.ProviderTypeOpenAI, info.ProviderType)
		assert.Equal(t, "gpt-4o", info.ModelName)
		assert.Equal(t, "gpt-4o", info.FullModel)
	})

	t.Run("handles model with default anthropic provider", func(t *testing.T) {
		info := service.ParseModelName("claude-sonnet-4", models.ProviderTypeAnthropic)

		assert.Equal(t, models.ProviderTypeAnthropic, info.ProviderType)
		assert.Equal(t, "claude-sonnet-4", info.ModelName)
	})

	t.Run("normalizes provider type to lowercase", func(t *testing.T) {
		info := service.ParseModelName("OpenAI/gpt-4o", models.ProviderTypeOpenAI)

		assert.Equal(t, "openai", info.ProviderType)
		assert.Equal(t, "gpt-4o", info.ModelName)
	})

	t.Run("handles local provider prefix", func(t *testing.T) {
		info := service.ParseModelName("local/my-model", models.ProviderTypeLocal)

		assert.Equal(t, "local", info.ProviderType)
		assert.Equal(t, "my-model", info.ModelName)
	})

	t.Run("handles model with slash in name after provider prefix", func(t *testing.T) {
		// e.g., "openai/gpt-4/turbo" should be parsed as provider=openai, model=gpt-4/turbo
		info := service.ParseModelName("openai/gpt-4/turbo", models.ProviderTypeOpenAI)

		assert.Equal(t, "openai", info.ProviderType)
		assert.Equal(t, "gpt-4/turbo", info.ModelName)
	})
}

func TestProxyService_TransformToAnthropic(t *testing.T) {
	db := setupProxyTestDB(t)
	service := createProxyTestServices(t, db)

	t.Run("transforms basic OpenAI request to Anthropic format", func(t *testing.T) {
		openAIReq := &OpenAIChatRequest{
			Model: "claude-sonnet-4",
			Messages: []OpenAIMessage{
				{Role: "user", Content: "Hello"},
			},
		}

		body, err := service.transformToAnthropic(openAIReq, "claude-sonnet-4")
		require.NoError(t, err)

		var anthropicReq AnthropicRequest
		err = json.Unmarshal(body, &anthropicReq)
		require.NoError(t, err)

		assert.Equal(t, "claude-sonnet-4", anthropicReq.Model)
		assert.Equal(t, DefaultMaxTokens, anthropicReq.MaxTokens)
		assert.Len(t, anthropicReq.Messages, 1)
		assert.Equal(t, "user", anthropicReq.Messages[0].Role)
		assert.Equal(t, "Hello", anthropicReq.Messages[0].Content)
	})

	t.Run("extracts system message to separate field", func(t *testing.T) {
		openAIReq := &OpenAIChatRequest{
			Model: "claude-sonnet-4",
			Messages: []OpenAIMessage{
				{Role: "system", Content: "You are a helpful assistant."},
				{Role: "user", Content: "Hello"},
			},
		}

		body, err := service.transformToAnthropic(openAIReq, "claude-sonnet-4")
		require.NoError(t, err)

		var anthropicReq AnthropicRequest
		err = json.Unmarshal(body, &anthropicReq)
		require.NoError(t, err)

		assert.Equal(t, "You are a helpful assistant.", anthropicReq.System)
		assert.Len(t, anthropicReq.Messages, 1) // Only user message
		assert.Equal(t, "user", anthropicReq.Messages[0].Role)
	})

	t.Run("concatenates multiple system messages", func(t *testing.T) {
		openAIReq := &OpenAIChatRequest{
			Model: "claude-sonnet-4",
			Messages: []OpenAIMessage{
				{Role: "system", Content: "First instruction."},
				{Role: "system", Content: "Second instruction."},
				{Role: "user", Content: "Hello"},
			},
		}

		body, err := service.transformToAnthropic(openAIReq, "claude-sonnet-4")
		require.NoError(t, err)

		var anthropicReq AnthropicRequest
		err = json.Unmarshal(body, &anthropicReq)
		require.NoError(t, err)

		assert.Contains(t, anthropicReq.System, "First instruction.")
		assert.Contains(t, anthropicReq.System, "Second instruction.")
	})

	t.Run("preserves max_tokens when provided", func(t *testing.T) {
		maxTokens := 1000
		openAIReq := &OpenAIChatRequest{
			Model: "claude-sonnet-4",
			Messages: []OpenAIMessage{
				{Role: "user", Content: "Hello"},
			},
			MaxTokens: &maxTokens,
		}

		body, err := service.transformToAnthropic(openAIReq, "claude-sonnet-4")
		require.NoError(t, err)

		var anthropicReq AnthropicRequest
		err = json.Unmarshal(body, &anthropicReq)
		require.NoError(t, err)

		assert.Equal(t, 1000, anthropicReq.MaxTokens)
	})

	t.Run("preserves temperature", func(t *testing.T) {
		temp := 0.7
		openAIReq := &OpenAIChatRequest{
			Model: "claude-sonnet-4",
			Messages: []OpenAIMessage{
				{Role: "user", Content: "Hello"},
			},
			Temperature: &temp,
		}

		body, err := service.transformToAnthropic(openAIReq, "claude-sonnet-4")
		require.NoError(t, err)

		var anthropicReq AnthropicRequest
		err = json.Unmarshal(body, &anthropicReq)
		require.NoError(t, err)

		assert.NotNil(t, anthropicReq.Temperature)
		assert.Equal(t, 0.7, *anthropicReq.Temperature)
	})

	t.Run("converts stop string to stop_sequences", func(t *testing.T) {
		openAIReq := &OpenAIChatRequest{
			Model: "claude-sonnet-4",
			Messages: []OpenAIMessage{
				{Role: "user", Content: "Hello"},
			},
			Stop: "END",
		}

		body, err := service.transformToAnthropic(openAIReq, "claude-sonnet-4")
		require.NoError(t, err)

		var anthropicReq AnthropicRequest
		err = json.Unmarshal(body, &anthropicReq)
		require.NoError(t, err)

		assert.Contains(t, anthropicReq.StopSequences, "END")
	})

	t.Run("converts stop array to stop_sequences", func(t *testing.T) {
		openAIReq := &OpenAIChatRequest{
			Model: "claude-sonnet-4",
			Messages: []OpenAIMessage{
				{Role: "user", Content: "Hello"},
			},
			Stop: []interface{}{"END", "STOP"},
		}

		body, err := service.transformToAnthropic(openAIReq, "claude-sonnet-4")
		require.NoError(t, err)

		var anthropicReq AnthropicRequest
		err = json.Unmarshal(body, &anthropicReq)
		require.NoError(t, err)

		assert.Contains(t, anthropicReq.StopSequences, "END")
		assert.Contains(t, anthropicReq.StopSequences, "STOP")
	})

	t.Run("preserves stream flag", func(t *testing.T) {
		stream := true
		openAIReq := &OpenAIChatRequest{
			Model: "claude-sonnet-4",
			Messages: []OpenAIMessage{
				{Role: "user", Content: "Hello"},
			},
			Stream: &stream,
		}

		body, err := service.transformToAnthropic(openAIReq, "claude-sonnet-4")
		require.NoError(t, err)

		var anthropicReq AnthropicRequest
		err = json.Unmarshal(body, &anthropicReq)
		require.NoError(t, err)

		assert.NotNil(t, anthropicReq.Stream)
		assert.True(t, *anthropicReq.Stream)
	})

	t.Run("maps function role to user", func(t *testing.T) {
		openAIReq := &OpenAIChatRequest{
			Model: "claude-sonnet-4",
			Messages: []OpenAIMessage{
				{Role: "user", Content: "What's the weather?"},
				{Role: "function", Content: "The weather is sunny."},
			},
		}

		body, err := service.transformToAnthropic(openAIReq, "claude-sonnet-4")
		require.NoError(t, err)

		var anthropicReq AnthropicRequest
		err = json.Unmarshal(body, &anthropicReq)
		require.NoError(t, err)

		assert.Len(t, anthropicReq.Messages, 2)
		assert.Equal(t, "user", anthropicReq.Messages[1].Role)
	})

	t.Run("fails for empty messages", func(t *testing.T) {
		openAIReq := &OpenAIChatRequest{
			Model:    "claude-sonnet-4",
			Messages: []OpenAIMessage{},
		}

		_, err := service.transformToAnthropic(openAIReq, "claude-sonnet-4")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least one")
	})

	t.Run("fails for only system messages", func(t *testing.T) {
		openAIReq := &OpenAIChatRequest{
			Model: "claude-sonnet-4",
			Messages: []OpenAIMessage{
				{Role: "system", Content: "You are helpful."},
			},
		}

		_, err := service.transformToAnthropic(openAIReq, "claude-sonnet-4")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least one")
	})

	t.Run("handles array content in OpenAI messages", func(t *testing.T) {
		openAIReq := &OpenAIChatRequest{
			Model: "claude-sonnet-4",
			Messages: []OpenAIMessage{
				{
					Role: "user",
					Content: []interface{}{
						map[string]interface{}{"type": "text", "text": "Hello"},
						map[string]interface{}{"type": "text", "text": "world!"},
					},
				},
			},
		}

		body, err := service.transformToAnthropic(openAIReq, "claude-sonnet-4")
		require.NoError(t, err)

		var anthropicReq AnthropicRequest
		err = json.Unmarshal(body, &anthropicReq)
		require.NoError(t, err)

		assert.Len(t, anthropicReq.Messages, 1)
		assert.Equal(t, "user", anthropicReq.Messages[0].Role)

		// Verification of Content is tricky because it's interface{}, but we can check if it serialized correctly
		contentBytes, _ := json.Marshal(anthropicReq.Messages[0].Content)
		assert.Contains(t, string(contentBytes), "Hello")
		assert.Contains(t, string(contentBytes), "world!")
	})
}

func TestProxyService_ExtractUsageFromResponse(t *testing.T) {
	db := setupProxyTestDB(t)
	service := createProxyTestServices(t, db)

	t.Run("extracts OpenAI usage format", func(t *testing.T) {
		response := `{
			"id": "chatcmpl-123",
			"choices": [{"message": {"content": "Hello!"}}],
			"usage": {
				"prompt_tokens": 10,
				"completion_tokens": 5,
				"total_tokens": 15
			}
		}`

		result := &ProxyResult{StatusCode: 200}
		service.extractUsageFromResponse([]byte(response), models.ProviderTypeOpenAI, result)

		assert.Equal(t, 10, result.InputTokens)
		assert.Equal(t, 5, result.OutputTokens)
		assert.Equal(t, 15, result.TotalTokens)
	})

	t.Run("extracts Anthropic usage format", func(t *testing.T) {
		response := `{
			"id": "msg_123",
			"content": [{"text": "Hello!"}],
			"usage": {
				"input_tokens": 20,
				"output_tokens": 10
			}
		}`

		result := &ProxyResult{StatusCode: 200}
		service.extractUsageFromResponse([]byte(response), models.ProviderTypeAnthropic, result)

		assert.Equal(t, 20, result.InputTokens)
		assert.Equal(t, 10, result.OutputTokens)
		assert.Equal(t, 30, result.TotalTokens) // Calculated
	})

	t.Run("skips extraction for error responses", func(t *testing.T) {
		response := `{"error": "something went wrong"}`

		result := &ProxyResult{StatusCode: 500}
		service.extractUsageFromResponse([]byte(response), models.ProviderTypeOpenAI, result)

		assert.Equal(t, 0, result.InputTokens)
		assert.Equal(t, 0, result.OutputTokens)
	})

	t.Run("handles invalid JSON gracefully", func(t *testing.T) {
		response := `not valid json`

		result := &ProxyResult{StatusCode: 200}
		// Should not panic
		service.extractUsageFromResponse([]byte(response), models.ProviderTypeOpenAI, result)

		assert.Equal(t, 0, result.InputTokens)
	})

	t.Run("handles missing usage field gracefully", func(t *testing.T) {
		response := `{"id": "chatcmpl-123", "choices": []}`

		result := &ProxyResult{StatusCode: 200}
		service.extractUsageFromResponse([]byte(response), models.ProviderTypeOpenAI, result)

		assert.Equal(t, 0, result.InputTokens)
	})
}

func TestProxyService_GetProxyKeyFromRequest(t *testing.T) {
	// Note: This test would require mocking gin.Context
	// For now, we test the logic paths documented in the function
}

func TestProxyService_HandleProviderError(t *testing.T) {
	db := setupProxyTestDB(t)
	service := createProxyTestServices(t, db)

	t.Run("returns authentication error for 401", func(t *testing.T) {
		code, body := service.HandleProviderError(401, "Unauthorized")

		assert.Equal(t, 401, code)
		errMap := body["error"].(map[string]interface{})
		assert.Equal(t, "authentication_error", errMap["type"])
		assert.Equal(t, "provider_auth_error", errMap["code"])
	})

	t.Run("returns rate limit error for 429", func(t *testing.T) {
		code, body := service.HandleProviderError(429, "Rate limited")

		assert.Equal(t, 429, code)
		errMap := body["error"].(map[string]interface{})
		assert.Equal(t, "rate_limit_error", errMap["type"])
		assert.Equal(t, "provider_rate_limit", errMap["code"])
	})

	t.Run("returns server error for 502", func(t *testing.T) {
		code, body := service.HandleProviderError(502, "Bad Gateway")

		assert.Equal(t, 502, code)
		errMap := body["error"].(map[string]interface{})
		assert.Equal(t, "server_error", errMap["type"])
		assert.Equal(t, "provider_unavailable", errMap["code"])
	})

	t.Run("returns server error for 503", func(t *testing.T) {
		code, body := service.HandleProviderError(503, "Service Unavailable")

		assert.Equal(t, 502, code) // Maps to 502
		errMap := body["error"].(map[string]interface{})
		assert.Equal(t, "provider_unavailable", errMap["code"])
	})

	t.Run("returns server error for 504", func(t *testing.T) {
		code, body := service.HandleProviderError(504, "Gateway Timeout")

		assert.Equal(t, 502, code) // Maps to 502
		errMap := body["error"].(map[string]interface{})
		assert.Equal(t, "provider_unavailable", errMap["code"])
	})

	t.Run("returns generic error for other status codes", func(t *testing.T) {
		code, body := service.HandleProviderError(400, "Bad Request")

		assert.Equal(t, 400, code)
		errMap := body["error"].(map[string]interface{})
		assert.Equal(t, "api_error", errMap["type"])
		assert.Equal(t, "proxy_error", errMap["code"])
		assert.Equal(t, "Bad Request", errMap["message"])
	})
}

func TestProxyService_ListModels(t *testing.T) {
	db := setupProxyTestDB(t)
	service := createProxyTestServices(t, db)

	t.Run("returns OpenAI models for OpenAI provider", func(t *testing.T) {
		provider := &models.Provider{
			ProviderType: models.ProviderTypeOpenAI,
		}

		result, err := service.ListModels(provider)
		require.NoError(t, err)

		// Convert to JSON to check contents
		jsonBytes, err := json.Marshal(result)
		require.NoError(t, err)
		jsonStr := string(jsonBytes)

		assert.Contains(t, jsonStr, `"object":"list"`)
		assert.Contains(t, jsonStr, "openai/gpt-4o")
		assert.Contains(t, jsonStr, "openai/gpt-4o-mini")
	})

	t.Run("returns Anthropic models for Anthropic provider", func(t *testing.T) {
		provider := &models.Provider{
			ProviderType: models.ProviderTypeAnthropic,
		}

		result, err := service.ListModels(provider)
		require.NoError(t, err)

		jsonBytes, err := json.Marshal(result)
		require.NoError(t, err)
		jsonStr := string(jsonBytes)

		assert.Contains(t, jsonStr, "anthropic/claude-sonnet-4-20250514")
	})

	t.Run("includes default model if configured", func(t *testing.T) {
		provider := &models.Provider{
			ProviderType: models.ProviderTypeLocal,
			DefaultModel: "my-custom-model",
		}

		result, err := service.ListModels(provider)
		require.NoError(t, err)

		jsonBytes, _ := json.Marshal(result)
		assert.Contains(t, string(jsonBytes), "local/my-custom-model")
	})
}

func TestProxyService_ValidateAndGetProvider(t *testing.T) {
	t.Run("returns key and provider for valid API key", func(t *testing.T) {
		db := setupProxyTestDB(t)
		keyService := NewKeyService(db)
		providerService := NewProviderService(db)
		usageService := NewUsageService(db)
		service := NewProxyService(keyService, providerService, usageService, nil)

		// Create provider
		provider := &models.Provider{
			UserID:       1,
			Name:         "Test Provider",
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "test-api-key",
			IsActive:     true,
		}
		err := db.Create(provider).Error
		require.NoError(t, err)

		// Create key
		createReq := &CreateKeyRequest{
			AllowedProviders: []ProviderSelection{{ProviderID: provider.ID}},
			Name:             "Test Key",
		}
		created, err := keyService.CreateKey(1, createReq)
		require.NoError(t, err)

		// Validate
		key, prov, err := service.ValidateAndGetProvider(created.Key)
		require.NoError(t, err)
		assert.NotNil(t, key)
		assert.NotNil(t, prov)
		assert.Equal(t, provider.ID, prov.ID)
	})

	t.Run("fails for invalid API key", func(t *testing.T) {
		db := setupProxyTestDB(t)
		service := createProxyTestServices(t, db)

		_, _, err := service.ValidateAndGetProvider("sk-smoothllm-invalid")
		assert.Error(t, err)
	})

	t.Run("fails for inactive provider", func(t *testing.T) {
		db := setupProxyTestDB(t)
		keyService := NewKeyService(db)
		providerService := NewProviderService(db)
		usageService := NewUsageService(db)
		service := NewProxyService(keyService, providerService, usageService, nil)

		// Create active provider first
		provider := &models.Provider{
			UserID:       1,
			Name:         "Test Provider",
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "test-api-key",
			IsActive:     true,
		}
		err := db.Create(provider).Error
		require.NoError(t, err)

		// Create key
		createReq := &CreateKeyRequest{
			AllowedProviders: []ProviderSelection{{ProviderID: provider.ID}},
			Name:             "Test Key",
		}
		created, err := keyService.CreateKey(1, createReq)
		require.NoError(t, err)

		// Deactivate provider using map to ensure zero-value bool is included
		err = db.Model(&models.Provider{}).Where("id = ?", provider.ID).Updates(map[string]interface{}{"is_active": false}).Error
		require.NoError(t, err)

		// Validate should fail
		_, _, err = service.ValidateAndGetProvider(created.Key)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no active providers")
	})
}

func TestOpenAIChatRequest(t *testing.T) {
	t.Run("deserializes correctly", func(t *testing.T) {
		jsonStr := `{
			"model": "gpt-4o",
			"messages": [
				{"role": "system", "content": "You are helpful."},
				{"role": "user", "content": "Hello"}
			],
			"max_tokens": 100,
			"temperature": 0.7,
			"stream": true
		}`

		var req OpenAIChatRequest
		err := json.Unmarshal([]byte(jsonStr), &req)
		require.NoError(t, err)

		assert.Equal(t, "gpt-4o", req.Model)
		assert.Len(t, req.Messages, 2)
		assert.Equal(t, "system", req.Messages[0].Role)
		assert.NotNil(t, req.MaxTokens)
		assert.Equal(t, 100, *req.MaxTokens)
		assert.NotNil(t, req.Temperature)
		assert.Equal(t, 0.7, *req.Temperature)
		assert.NotNil(t, req.Stream)
		assert.True(t, *req.Stream)
	})
}

func TestAnthropicRequest(t *testing.T) {
	t.Run("serializes correctly", func(t *testing.T) {
		temp := 0.8
		stream := true
		req := AnthropicRequest{
			Model:         "claude-sonnet-4",
			MaxTokens:     4096,
			System:        "You are helpful.",
			Temperature:   &temp,
			Stream:        &stream,
			StopSequences: []string{"END"},
			Messages: []AnthropicMessage{
				{Role: "user", Content: "Hello"},
			},
		}

		jsonBytes, err := json.Marshal(req)
		require.NoError(t, err)

		jsonStr := string(jsonBytes)
		assert.Contains(t, jsonStr, `"model":"claude-sonnet-4"`)
		assert.Contains(t, jsonStr, `"max_tokens":4096`)
		assert.Contains(t, jsonStr, `"system":"You are helpful."`)
		assert.Contains(t, jsonStr, `"temperature":0.8`)
		assert.Contains(t, jsonStr, `"stream":true`)
		assert.Contains(t, jsonStr, `"stop_sequences":["END"]`)
	})
}

func TestProxyResult(t *testing.T) {
	t.Run("tracks request metrics", func(t *testing.T) {
		result := &ProxyResult{
			StatusCode:   200,
			InputTokens:  100,
			OutputTokens: 50,
			TotalTokens:  150,
			Model:        "gpt-4o",
		}

		assert.Equal(t, 200, result.StatusCode)
		assert.Equal(t, 100, result.InputTokens)
		assert.Equal(t, 50, result.OutputTokens)
		assert.Equal(t, 150, result.TotalTokens)
	})
}
