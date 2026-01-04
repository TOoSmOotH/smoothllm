package services

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/smoothweb/backend/internal/custom/models"
)

// setupKeyTestDB creates an in-memory SQLite database for testing
func setupKeyTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate the required tables
	err = db.AutoMigrate(&models.Provider{}, &models.ProxyAPIKey{})
	require.NoError(t, err)

	return db
}

// createTestProvider creates a test provider in the database
func createTestProvider(t *testing.T, db *gorm.DB, userID uint) *models.Provider {
	provider := &models.Provider{
		UserID:       userID,
		Name:         "Test Provider",
		ProviderType: models.ProviderTypeOpenAI,
		APIKey:       "test-api-key",
		IsActive:     true,
	}
	err := db.Create(provider).Error
	require.NoError(t, err)
	return provider
}

func TestKeyService_GenerateKey(t *testing.T) {
	db := setupKeyTestDB(t)
	service := NewKeyService(db)

	t.Run("generates key with correct prefix", func(t *testing.T) {
		key, err := service.generateKey()
		require.NoError(t, err)

		assert.True(t, strings.HasPrefix(key, models.ProxyAPIKeyPrefix),
			"key should start with %s, got %s", models.ProxyAPIKeyPrefix, key)
	})

	t.Run("generates key with correct length", func(t *testing.T) {
		key, err := service.generateKey()
		require.NoError(t, err)

		// Prefix + 64 hex chars (32 bytes = 64 hex chars)
		expectedLength := len(models.ProxyAPIKeyPrefix) + 64
		assert.Equal(t, expectedLength, len(key),
			"key should be %d chars, got %d", expectedLength, len(key))
	})

	t.Run("generates unique keys", func(t *testing.T) {
		keys := make(map[string]bool)
		for i := 0; i < 100; i++ {
			key, err := service.generateKey()
			require.NoError(t, err)
			assert.False(t, keys[key], "generated duplicate key")
			keys[key] = true
		}
	})

	t.Run("key random part is valid hex", func(t *testing.T) {
		key, err := service.generateKey()
		require.NoError(t, err)

		randomPart := strings.TrimPrefix(key, models.ProxyAPIKeyPrefix)
		_, err = hex.DecodeString(randomPart)
		assert.NoError(t, err, "random part should be valid hex")
	})
}

func TestKeyService_HashKey(t *testing.T) {
	db := setupKeyTestDB(t)
	service := NewKeyService(db)

	t.Run("produces consistent hash for same input", func(t *testing.T) {
		key := "sk-smoothllm-abcdef123456"
		hash1 := service.HashKey(key)
		hash2 := service.HashKey(key)

		assert.Equal(t, hash1, hash2, "same input should produce same hash")
	})

	t.Run("produces different hashes for different inputs", func(t *testing.T) {
		key1 := "sk-smoothllm-abcdef123456"
		key2 := "sk-smoothllm-xyz789abc123"
		hash1 := service.HashKey(key1)
		hash2 := service.HashKey(key2)

		assert.NotEqual(t, hash1, hash2, "different inputs should produce different hashes")
	})

	t.Run("produces valid SHA256 hash", func(t *testing.T) {
		key := "sk-smoothllm-test123"
		hash := service.HashKey(key)

		// SHA256 produces 32 bytes = 64 hex chars
		assert.Equal(t, 64, len(hash), "hash should be 64 hex chars")

		// Verify it's valid hex
		_, err := hex.DecodeString(hash)
		assert.NoError(t, err, "hash should be valid hex")

		// Verify it matches expected SHA256
		expected := sha256.Sum256([]byte(key))
		expectedHex := hex.EncodeToString(expected[:])
		assert.Equal(t, expectedHex, hash)
	})
}

func TestKeyService_ExtractKeyPrefix(t *testing.T) {
	db := setupKeyTestDB(t)
	service := NewKeyService(db)

	t.Run("extracts prefix correctly for standard key", func(t *testing.T) {
		// Generate a real key to test
		fullKey, err := service.generateKey()
		require.NoError(t, err)

		prefix := service.extractKeyPrefix(fullKey)

		// Should start with the standard prefix
		assert.True(t, strings.HasPrefix(prefix, models.ProxyAPIKeyPrefix))

		// Should end with "..."
		assert.True(t, strings.Contains(prefix, "..."))

		// Prefix should be shorter than full key
		assert.Less(t, len(prefix), len(fullKey))
	})

	t.Run("preserves first 6 and last 4 characters of random part", func(t *testing.T) {
		fullKey := models.ProxyAPIKeyPrefix + "abcdef1234567890ghijklmn"
		prefix := service.extractKeyPrefix(fullKey)

		expected := models.ProxyAPIKeyPrefix + "abcdef...klmn"
		assert.Equal(t, expected, prefix)
	})
}

func TestKeyService_CreateKey(t *testing.T) {
	t.Run("creates key successfully", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)
		provider := createTestProvider(t, db, 1)

		req := &CreateKeyRequest{
			ProviderID: provider.ID,
			Name:       "Test Key",
		}

		result, err := service.CreateKey(1, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.Key)
		assert.True(t, strings.HasPrefix(result.Key, models.ProxyAPIKeyPrefix))
		assert.Equal(t, "Test Key", result.Name)
		assert.True(t, result.IsActive)
		assert.Equal(t, provider.ID, result.ProviderID)
	})

	t.Run("returns full key only on creation", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)
		provider := createTestProvider(t, db, 1)

		req := &CreateKeyRequest{
			ProviderID: provider.ID,
			Name:       "Test Key",
		}

		createResult, err := service.CreateKey(1, req)
		require.NoError(t, err)
		assert.NotEmpty(t, createResult.Key)

		// Get the key again - should not return full key
		getResult, err := service.GetKey(1, createResult.ID)
		require.NoError(t, err)
		assert.NotNil(t, getResult)
		// The KeyResponse doesn't have the full key field
	})

	t.Run("fails for non-existent provider", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)

		req := &CreateKeyRequest{
			ProviderID: 999,
			Name:       "Test Key",
		}

		result, err := service.CreateKey(1, req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "provider not found")
	})

	t.Run("fails for inactive provider", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)

		// Create provider first as active, then deactivate it
		// (GORM doesn't properly set IsActive:false on create because false is zero value)
		provider := &models.Provider{
			UserID:       1,
			Name:         "Inactive Provider",
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "test-api-key",
			IsActive:     true,
		}
		err := db.Create(provider).Error
		require.NoError(t, err)

		// Now deactivate it
		err = db.Model(provider).Update("is_active", false).Error
		require.NoError(t, err)

		req := &CreateKeyRequest{
			ProviderID: provider.ID,
			Name:       "Test Key",
		}

		result, err := service.CreateKey(1, req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not active")
	})

	t.Run("fails for provider belonging to different user", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)
		provider := createTestProvider(t, db, 1) // User 1's provider

		req := &CreateKeyRequest{
			ProviderID: provider.ID,
			Name:       "Test Key",
		}

		result, err := service.CreateKey(2, req) // User 2 trying to use User 1's provider
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "provider not found")
	})

	t.Run("fails for name longer than 100 characters", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)
		provider := createTestProvider(t, db, 1)

		req := &CreateKeyRequest{
			ProviderID: provider.ID,
			Name:       strings.Repeat("a", 101),
		}

		result, err := service.CreateKey(1, req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "100 characters")
	})

	t.Run("fails for past expiration date", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)
		provider := createTestProvider(t, db, 1)

		pastDate := time.Now().Add(-24 * time.Hour)
		req := &CreateKeyRequest{
			ProviderID: provider.ID,
			Name:       "Test Key",
			ExpiresAt:  &pastDate,
		}

		result, err := service.CreateKey(1, req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "future")
	})

	t.Run("creates key with expiration date", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)
		provider := createTestProvider(t, db, 1)

		futureDate := time.Now().Add(30 * 24 * time.Hour)
		req := &CreateKeyRequest{
			ProviderID: provider.ID,
			Name:       "Expiring Key",
			ExpiresAt:  &futureDate,
		}

		result, err := service.CreateKey(1, req)
		require.NoError(t, err)
		assert.NotNil(t, result.ExpiresAt)
	})
}

func TestKeyService_ValidateKey(t *testing.T) {
	t.Run("validates correct key", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)
		provider := createTestProvider(t, db, 1)

		// Create a key
		createReq := &CreateKeyRequest{
			ProviderID: provider.ID,
			Name:       "Test Key",
		}
		created, err := service.CreateKey(1, createReq)
		require.NoError(t, err)

		// Validate it
		key, err := service.ValidateKey(created.Key)
		require.NoError(t, err)
		assert.NotNil(t, key)
		assert.Equal(t, created.ID, key.ID)
	})

	t.Run("fails for invalid key format", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)

		_, err := service.ValidateKey("invalid-key-format")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid key format")
	})

	t.Run("fails for non-existent key", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)

		_, err := service.ValidateKey(models.ProxyAPIKeyPrefix + "nonexistent0000000000000000000000000000000000000000000000000000")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid API key")
	})

	t.Run("fails for inactive key", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)
		provider := createTestProvider(t, db, 1)

		// Create a key
		createReq := &CreateKeyRequest{
			ProviderID: provider.ID,
			Name:       "Test Key",
		}
		created, err := service.CreateKey(1, createReq)
		require.NoError(t, err)

		// Deactivate it
		isActive := false
		_, err = service.UpdateKey(1, created.ID, &UpdateKeyRequest{IsActive: &isActive})
		require.NoError(t, err)

		// Try to validate
		_, err = service.ValidateKey(created.Key)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "inactive")
	})

	t.Run("fails for expired key", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)
		provider := createTestProvider(t, db, 1)

		// Create a key with a future expiration
		futureDate := time.Now().Add(1 * time.Hour)
		createReq := &CreateKeyRequest{
			ProviderID: provider.ID,
			Name:       "Expiring Key",
			ExpiresAt:  &futureDate,
		}
		created, err := service.CreateKey(1, createReq)
		require.NoError(t, err)

		// Manually expire it by updating the database directly
		pastDate := time.Now().Add(-1 * time.Hour)
		err = db.Model(&models.ProxyAPIKey{}).Where("id = ?", created.ID).Update("expires_at", pastDate).Error
		require.NoError(t, err)

		// Try to validate
		_, err = service.ValidateKey(created.Key)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expired")
	})

	t.Run("updates last used timestamp", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)
		provider := createTestProvider(t, db, 1)

		// Create a key
		createReq := &CreateKeyRequest{
			ProviderID: provider.ID,
			Name:       "Test Key",
		}
		created, err := service.CreateKey(1, createReq)
		require.NoError(t, err)

		// First validation - last used should be set
		beforeValidation := time.Now()
		key, err := service.ValidateKey(created.Key)
		require.NoError(t, err)
		assert.NotNil(t, key.LastUsedAt)
		assert.True(t, key.LastUsedAt.After(beforeValidation) || key.LastUsedAt.Equal(beforeValidation))
	})
}

func TestKeyService_ListKeys(t *testing.T) {
	t.Run("lists keys for user", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)
		provider := createTestProvider(t, db, 1)

		// Create multiple keys
		for i := 0; i < 3; i++ {
			req := &CreateKeyRequest{
				ProviderID: provider.ID,
				Name:       "Key " + string(rune('A'+i)),
			}
			_, err := service.CreateKey(1, req)
			require.NoError(t, err)
		}

		// List keys
		keys, err := service.ListKeys(1)
		require.NoError(t, err)
		assert.Len(t, keys, 3)
	})

	t.Run("returns empty list for user with no keys", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)

		keys, err := service.ListKeys(999)
		require.NoError(t, err)
		assert.Empty(t, keys)
	})

	t.Run("only returns keys for specified user", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)

		// Create providers for two users
		provider1 := createTestProvider(t, db, 1)
		provider2 := &models.Provider{
			UserID:       2,
			Name:         "User 2 Provider",
			ProviderType: models.ProviderTypeOpenAI,
			APIKey:       "test-api-key-2",
			IsActive:     true,
		}
		err := db.Create(provider2).Error
		require.NoError(t, err)

		// Create keys for both users
		_, err = service.CreateKey(1, &CreateKeyRequest{ProviderID: provider1.ID, Name: "User 1 Key"})
		require.NoError(t, err)
		_, err = service.CreateKey(2, &CreateKeyRequest{ProviderID: provider2.ID, Name: "User 2 Key"})
		require.NoError(t, err)

		// List keys for user 1
		keys, err := service.ListKeys(1)
		require.NoError(t, err)
		assert.Len(t, keys, 1)
		assert.Equal(t, "User 1 Key", keys[0].Name)
	})
}

func TestKeyService_UpdateKey(t *testing.T) {
	t.Run("updates key name", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)
		provider := createTestProvider(t, db, 1)

		// Create a key
		created, err := service.CreateKey(1, &CreateKeyRequest{
			ProviderID: provider.ID,
			Name:       "Original Name",
		})
		require.NoError(t, err)

		// Update name
		newName := "Updated Name"
		updated, err := service.UpdateKey(1, created.ID, &UpdateKeyRequest{Name: &newName})
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", updated.Name)
	})

	t.Run("deactivates key", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)
		provider := createTestProvider(t, db, 1)

		// Create a key
		created, err := service.CreateKey(1, &CreateKeyRequest{
			ProviderID: provider.ID,
			Name:       "Active Key",
		})
		require.NoError(t, err)
		assert.True(t, created.IsActive)

		// Deactivate
		isActive := false
		updated, err := service.UpdateKey(1, created.ID, &UpdateKeyRequest{IsActive: &isActive})
		require.NoError(t, err)
		assert.False(t, updated.IsActive)
	})

	t.Run("fails for non-existent key", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)

		newName := "New Name"
		_, err := service.UpdateKey(1, 999, &UpdateKeyRequest{Name: &newName})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("fails for key belonging to different user", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)
		provider := createTestProvider(t, db, 1)

		// Create a key for user 1
		created, err := service.CreateKey(1, &CreateKeyRequest{
			ProviderID: provider.ID,
			Name:       "User 1 Key",
		})
		require.NoError(t, err)

		// Try to update as user 2
		newName := "Hacked"
		_, err = service.UpdateKey(2, created.ID, &UpdateKeyRequest{Name: &newName})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestKeyService_DeleteKey(t *testing.T) {
	t.Run("deletes key successfully", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)
		provider := createTestProvider(t, db, 1)

		// Create a key
		created, err := service.CreateKey(1, &CreateKeyRequest{
			ProviderID: provider.ID,
			Name:       "To Delete",
		})
		require.NoError(t, err)

		// Delete it
		err = service.DeleteKey(1, created.ID)
		require.NoError(t, err)

		// Verify it's gone
		_, err = service.GetKey(1, created.ID)
		assert.Error(t, err)
	})

	t.Run("fails for non-existent key", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)

		err := service.DeleteKey(1, 999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("fails for key belonging to different user", func(t *testing.T) {
		db := setupKeyTestDB(t)
		service := NewKeyService(db)
		provider := createTestProvider(t, db, 1)

		// Create a key for user 1
		created, err := service.CreateKey(1, &CreateKeyRequest{
			ProviderID: provider.ID,
			Name:       "User 1 Key",
		})
		require.NoError(t, err)

		// Try to delete as user 2
		err = service.DeleteKey(2, created.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestProxyAPIKeyModel(t *testing.T) {
	t.Run("IsExpired returns false for nil expiration", func(t *testing.T) {
		key := &models.ProxyAPIKey{
			ExpiresAt: nil,
		}
		assert.False(t, key.IsExpired())
	})

	t.Run("IsExpired returns false for future expiration", func(t *testing.T) {
		future := time.Now().Add(24 * time.Hour)
		key := &models.ProxyAPIKey{
			ExpiresAt: &future,
		}
		assert.False(t, key.IsExpired())
	})

	t.Run("IsExpired returns true for past expiration", func(t *testing.T) {
		past := time.Now().Add(-24 * time.Hour)
		key := &models.ProxyAPIKey{
			ExpiresAt: &past,
		}
		assert.True(t, key.IsExpired())
	})

	t.Run("IsValid returns true for active non-expired key", func(t *testing.T) {
		key := &models.ProxyAPIKey{
			IsActive:  true,
			ExpiresAt: nil,
		}
		assert.True(t, key.IsValid())
	})

	t.Run("IsValid returns false for inactive key", func(t *testing.T) {
		key := &models.ProxyAPIKey{
			IsActive:  false,
			ExpiresAt: nil,
		}
		assert.False(t, key.IsValid())
	})

	t.Run("IsValid returns false for expired key", func(t *testing.T) {
		past := time.Now().Add(-24 * time.Hour)
		key := &models.ProxyAPIKey{
			IsActive:  true,
			ExpiresAt: &past,
		}
		assert.False(t, key.IsValid())
	})

	t.Run("UpdateLastUsed sets timestamp", func(t *testing.T) {
		key := &models.ProxyAPIKey{}
		assert.Nil(t, key.LastUsedAt)

		before := time.Now()
		key.UpdateLastUsed()
		after := time.Now()

		assert.NotNil(t, key.LastUsedAt)
		assert.True(t, key.LastUsedAt.After(before) || key.LastUsedAt.Equal(before))
		assert.True(t, key.LastUsedAt.Before(after) || key.LastUsedAt.Equal(after))
	})
}
