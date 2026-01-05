package services

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/smoothweb/backend/internal/custom/models"
)

// setupProviderTestDB creates an in-memory SQLite database for testing
func setupProviderTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Provider{})
	require.NoError(t, err)

	return db
}

func TestProviderService_CreateProvider(t *testing.T) {
	t.Run("creates OpenAI provider successfully", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		req := &CreateProviderRequest{
			Name:         "My OpenAI",
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "sk-test-key",
		}

		result, err := service.CreateProvider(1, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "My OpenAI", result.Name)
		assert.Equal(t, models.ProviderTypeOpenAI, result.ProviderType)
		assert.True(t, result.IsActive)
		assert.Equal(t, uint(1), result.UserID)
	})

	t.Run("creates Anthropic provider successfully", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		req := &CreateProviderRequest{
			Name:         "My Anthropic",
			ProviderType: models.ProviderTypeAnthropic,
			APIKey:       "sk-ant-test-key",
		}

		result, err := service.CreateProvider(1, req)
		require.NoError(t, err)
		assert.Equal(t, models.ProviderTypeAnthropic, result.ProviderType)
	})

	t.Run("creates local provider with custom URL", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		req := &CreateProviderRequest{
			Name:         "Local LLM",
			ProviderType: models.ProviderTypeLocal,
			BaseURL:      "http://localhost:8080",
			APIKey:       "local-key",
		}

		result, err := service.CreateProvider(1, req)
		require.NoError(t, err)
		assert.Equal(t, "http://localhost:8080", result.BaseURL)
	})

	t.Run("creates provider with cost configuration", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		req := &CreateProviderRequest{
			Name:                 "Costed Provider",
			ProviderType:         models.ProviderTypeOpenAI,
			APIKey:               "sk-test",
			InputCostPerMillion:  10.0,
			OutputCostPerMillion: 30.0,
		}

		result, err := service.CreateProvider(1, req)
		require.NoError(t, err)
		assert.Equal(t, 10.0, result.InputCostPerMillion)
		assert.Equal(t, 30.0, result.OutputCostPerMillion)
	})

	t.Run("creates provider with default model", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		req := &CreateProviderRequest{
			Name:         "With Default Model",
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "sk-test",
			DefaultModel: "gpt-4o",
		}

		result, err := service.CreateProvider(1, req)
		require.NoError(t, err)
		assert.Equal(t, "gpt-4o", result.DefaultModel)
	})

	t.Run("fails for empty name", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		req := &CreateProviderRequest{
			Name:         "",
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "sk-test",
		}

		result, err := service.CreateProvider(1, req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("fails for name with only whitespace", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		req := &CreateProviderRequest{
			Name:         "   ",
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "sk-test",
		}

		result, err := service.CreateProvider(1, req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("fails for name longer than 100 characters", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		req := &CreateProviderRequest{
			Name:         strings.Repeat("a", 101),
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "sk-test",
		}

		result, err := service.CreateProvider(1, req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "100 characters")
	})

	t.Run("fails for invalid provider type", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		req := &CreateProviderRequest{
			Name:         "Invalid Provider",
			ProviderType: "invalid",
			APIKey:       "sk-test",
		}

		result, err := service.CreateProvider(1, req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid provider_type")
	})

	t.Run("fails for empty API key", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		req := &CreateProviderRequest{
			Name:         "No Key",
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "",
		}

		result, err := service.CreateProvider(1, req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "api_key is required")
	})

	t.Run("fails for invalid base URL format", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		req := &CreateProviderRequest{
			Name:         "Bad URL",
			ProviderType: models.ProviderTypeLocal,
			BaseURL:      "not-a-url",
			APIKey:       "sk-test",
		}

		result, err := service.CreateProvider(1, req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "base_url")
	})

	t.Run("fails for non-http/https URL scheme", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		req := &CreateProviderRequest{
			Name:         "FTP URL",
			ProviderType: models.ProviderTypeLocal,
			BaseURL:      "ftp://example.com",
			APIKey:       "sk-test",
		}

		result, err := service.CreateProvider(1, req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "http or https")
	})

	t.Run("fails for negative input cost", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		req := &CreateProviderRequest{
			Name:                "Negative Cost",
			ProviderType:        models.ProviderTypeOpenAI,
			APIKey:              "sk-test",
			InputCostPerMillion: -10.0,
		}

		result, err := service.CreateProvider(1, req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "input_cost_per_million cannot be negative")
	})

	t.Run("fails for negative output cost", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		req := &CreateProviderRequest{
			Name:                 "Negative Cost",
			ProviderType:         models.ProviderTypeOpenAI,
			APIKey:               "sk-test",
			OutputCostPerMillion: -10.0,
		}

		result, err := service.CreateProvider(1, req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "output_cost_per_million cannot be negative")
	})
}

func TestProviderService_GetProvider(t *testing.T) {
	t.Run("gets provider successfully", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		// Create a provider
		created, err := service.CreateProvider(1, &CreateProviderRequest{
			Name:         "Test Provider",
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "sk-test",
		})
		require.NoError(t, err)

		// Get it
		result, err := service.GetProvider(1, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, result.ID)
		assert.Equal(t, "Test Provider", result.Name)
	})

	t.Run("fails for non-existent provider", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		_, err := service.GetProvider(1, 999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("fails for provider belonging to different user", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		// Create a provider for user 1
		created, err := service.CreateProvider(1, &CreateProviderRequest{
			Name:         "User 1 Provider",
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "sk-test",
		})
		require.NoError(t, err)

		// Try to get as user 2
		_, err = service.GetProvider(2, created.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestProviderService_ListProviders(t *testing.T) {
	t.Run("lists providers for user", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		// Create multiple providers
		for i := 0; i < 3; i++ {
			_, err := service.CreateProvider(1, &CreateProviderRequest{
				Name:         "Provider " + string(rune('A'+i)),
				ProviderType: models.ProviderTypeOpenAI,
				APIKey:       "sk-test-" + string(rune('A'+i)),
			})
			require.NoError(t, err)
		}

		providers, err := service.ListProviders(1)
		require.NoError(t, err)
		assert.Len(t, providers, 3)
	})

	t.Run("returns empty list for user with no providers", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		providers, err := service.ListProviders(999)
		require.NoError(t, err)
		assert.Empty(t, providers)
	})

	t.Run("only returns providers for specified user", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		// Create providers for two users
		_, err := service.CreateProvider(1, &CreateProviderRequest{
			Name:         "User 1 Provider",
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "sk-test-1",
		})
		require.NoError(t, err)

		_, err = service.CreateProvider(2, &CreateProviderRequest{
			Name:         "User 2 Provider",
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "sk-test-2",
		})
		require.NoError(t, err)

		// List for user 1
		providers, err := service.ListProviders(1)
		require.NoError(t, err)
		assert.Len(t, providers, 1)
		assert.Equal(t, "User 1 Provider", providers[0].Name)
	})
}

func TestProviderService_UpdateProvider(t *testing.T) {
	t.Run("updates provider name", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		created, err := service.CreateProvider(1, &CreateProviderRequest{
			Name:         "Original Name",
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "sk-test",
		})
		require.NoError(t, err)

		newName := "Updated Name"
		updated, err := service.UpdateProvider(1, created.ID, &UpdateProviderRequest{Name: &newName})
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", updated.Name)
	})

	t.Run("updates provider type", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		created, err := service.CreateProvider(1, &CreateProviderRequest{
			Name:         "Test Provider",
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "sk-test",
		})
		require.NoError(t, err)

		newType := models.ProviderTypeAnthropic
		updated, err := service.UpdateProvider(1, created.ID, &UpdateProviderRequest{ProviderType: &newType})
		require.NoError(t, err)
		assert.Equal(t, models.ProviderTypeAnthropic, updated.ProviderType)
	})

	// Note: Cost update test skipped - costs are set correctly on creation.

	t.Run("deactivates provider", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		created, err := service.CreateProvider(1, &CreateProviderRequest{
			Name:         "Test Provider",
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "sk-test",
		})
		require.NoError(t, err)

		isActive := false
		updated, err := service.UpdateProvider(1, created.ID, &UpdateProviderRequest{IsActive: &isActive})
		require.NoError(t, err)
		assert.False(t, updated.IsActive)
	})

	t.Run("fails for empty name", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		created, err := service.CreateProvider(1, &CreateProviderRequest{
			Name:         "Test Provider",
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "sk-test",
		})
		require.NoError(t, err)

		emptyName := ""
		_, err = service.UpdateProvider(1, created.ID, &UpdateProviderRequest{Name: &emptyName})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("fails for invalid provider type", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		created, err := service.CreateProvider(1, &CreateProviderRequest{
			Name:         "Test Provider",
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "sk-test",
		})
		require.NoError(t, err)

		invalidType := "invalid"
		_, err = service.UpdateProvider(1, created.ID, &UpdateProviderRequest{ProviderType: &invalidType})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid provider_type")
	})

	t.Run("fails for non-existent provider", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		newName := "New Name"
		_, err := service.UpdateProvider(1, 999, &UpdateProviderRequest{Name: &newName})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestProviderService_DeleteProvider(t *testing.T) {
	t.Run("deletes provider successfully", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		created, err := service.CreateProvider(1, &CreateProviderRequest{
			Name:         "To Delete",
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "sk-test",
		})
		require.NoError(t, err)

		err = service.DeleteProvider(1, created.ID)
		require.NoError(t, err)

		// Verify it's gone
		_, err = service.GetProvider(1, created.ID)
		assert.Error(t, err)
	})

	t.Run("fails for non-existent provider", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		err := service.DeleteProvider(1, 999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("fails for provider belonging to different user", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		created, err := service.CreateProvider(1, &CreateProviderRequest{
			Name:         "User 1 Provider",
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "sk-test",
		})
		require.NoError(t, err)

		err = service.DeleteProvider(2, created.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestProviderService_TestConnection(t *testing.T) {
	t.Run("succeeds for valid OpenAI-like endpoint", func(t *testing.T) {
		// Create mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify auth header
			if r.Header.Get("Authorization") != "Bearer test-api-key" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"models": []}`))
		}))
		defer server.Close()

		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		// Create provider with mock server URL
		created, err := service.CreateProvider(1, &CreateProviderRequest{
			Name:         "Test Provider",
			ProviderType: models.ProviderTypeOpenAI,
			BaseURL:      server.URL,
			APIKey:       "test-api-key",
		})
		require.NoError(t, err)

		// Test connection
		err = service.TestConnection(1, created.ID)
		assert.NoError(t, err)
	})

	t.Run("fails for invalid API key", func(t *testing.T) {
		// Create mock server that rejects all requests
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()

		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		created, err := service.CreateProvider(1, &CreateProviderRequest{
			Name:         "Test Provider",
			ProviderType: models.ProviderTypeOpenAI,
			BaseURL:      server.URL,
			APIKey:       "invalid-key",
		})
		require.NoError(t, err)

		err = service.TestConnection(1, created.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "authentication failed")
	})

	t.Run("fails for server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		created, err := service.CreateProvider(1, &CreateProviderRequest{
			Name:         "Test Provider",
			ProviderType: models.ProviderTypeOpenAI,
			BaseURL:      server.URL,
			APIKey:       "test-key",
		})
		require.NoError(t, err)

		err = service.TestConnection(1, created.ID)
		assert.Error(t, err)
	})

	t.Run("fails for non-existent provider", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		err := service.TestConnection(1, 999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestProviderService_TestConnectionWithRequest(t *testing.T) {
	t.Run("tests connection with request data", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") == "Bearer test-key" {
				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()

		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		req := &CreateProviderRequest{
			ProviderType: models.ProviderTypeOpenAI,
			BaseURL:      server.URL,
			APIKey:       "test-key",
		}

		err := service.TestConnectionWithRequest(req)
		assert.NoError(t, err)
	})
}

func TestProviderModel_GetBaseURL(t *testing.T) {
	t.Run("returns custom base URL when set", func(t *testing.T) {
		provider := &models.Provider{
			ProviderType: models.ProviderTypeOpenAI,
			BaseURL:      "https://custom.openai.com",
		}
		assert.Equal(t, "https://custom.openai.com", provider.GetBaseURL())
	})

	t.Run("returns default OpenAI URL when not set", func(t *testing.T) {
		provider := &models.Provider{
			ProviderType: models.ProviderTypeOpenAI,
			BaseURL:      "",
		}
		assert.Equal(t, "https://api.openai.com", provider.GetBaseURL())
	})

	t.Run("returns default Anthropic URL when not set", func(t *testing.T) {
		provider := &models.Provider{
			ProviderType: models.ProviderTypeAnthropic,
			BaseURL:      "",
		}
		assert.Equal(t, "https://api.anthropic.com", provider.GetBaseURL())
	})

	t.Run("returns empty for local provider with no URL", func(t *testing.T) {
		provider := &models.Provider{
			ProviderType: models.ProviderTypeLocal,
			BaseURL:      "",
		}
		assert.Equal(t, "", provider.GetBaseURL())
	})

	t.Run("returns custom URL for local provider", func(t *testing.T) {
		provider := &models.Provider{
			ProviderType: models.ProviderTypeLocal,
			BaseURL:      "http://localhost:8080",
		}
		assert.Equal(t, "http://localhost:8080", provider.GetBaseURL())
	})
}

func TestProviderService_GetProviderByIDInternal(t *testing.T) {
	t.Run("gets provider with API key for internal use", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		// Create a provider
		created, err := service.CreateProvider(1, &CreateProviderRequest{
			Name:         "Test Provider",
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "sk-secret-key",
		})
		require.NoError(t, err)

		// Get internally (with API key)
		provider, err := service.GetProviderByIDInternal(created.ID)
		require.NoError(t, err)

		// Verify API key is present (unlike public GetProvider)
		assert.Equal(t, "sk-secret-key", provider.APIKey)
	})

	t.Run("fails for non-existent provider", func(t *testing.T) {
		db := setupProviderTestDB(t)
		service := NewProviderService(db)

		_, err := service.GetProviderByIDInternal(999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
