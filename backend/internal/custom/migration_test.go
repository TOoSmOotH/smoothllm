package custom

import (
	"testing"

	"github.com/smoothweb/backend/internal/custom/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMigration_DropProviderID(t *testing.T) {
	// 1. Setup a database with the OLD schema using GORM
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Define the old model structure locally for setup
	type OldProxyAPIKey struct {
		gorm.Model
		UserID     uint   `gorm:"not null"`
		KeyHash    string `gorm:"not null"`
		KeyPrefix  string `gorm:"not null"`
		ProviderID uint   `gorm:"not null"` // The offending column
	}

	// Create tables using the old model and regular models
	err = db.Table("proxy_api_keys").AutoMigrate(&OldProxyAPIKey{})
	require.NoError(t, err)
	err = db.AutoMigrate(&models.Provider{}, &models.KeyAllowedProvider{})
	require.NoError(t, err)

	// Verify the column exists initially
	assert.True(t, db.Migrator().HasColumn("proxy_api_keys", "provider_id"), "provider_id column should exist initially")

	// Insert test data
	provider := models.Provider{UserID: 1, Name: "Test", ProviderType: "openai", APIKey: "key", IsActive: true}
	err = db.Create(&provider).Error
	require.NoError(t, err)

	oldKey := OldProxyAPIKey{UserID: 1, KeyHash: "hash", KeyPrefix: "sk-", ProviderID: provider.ID}
	err = db.Table("proxy_api_keys").Create(&oldKey).Error
	require.NoError(t, err)

	// 2. Run our custom AutoMigrate which calls migrateExistingData
	err = AutoMigrate(db)
	require.NoError(t, err)

	// 3. Verify the column is gone
	assert.False(t, db.Migrator().HasColumn("proxy_api_keys", "provider_id"), "provider_id column should have been dropped")

	// 4. Verify data was migrated to KeyAllowedProvider
	var count int64
	err = db.Model(&models.KeyAllowedProvider{}).Where("proxy_api_key_id = ? AND provider_id = ?", oldKey.ID, provider.ID).Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(1), count, "existing key should have been migrated to KeyAllowedProvider")

	// 5. Verify we can create a NEW key using the current GORM model
	newKey := models.ProxyAPIKey{
		UserID:    1,
		KeyHash:   "new-hash",
		KeyPrefix: "sk-",
		Name:      "New Key",
	}
	err = db.Create(&newKey).Error
	assert.NoError(t, err, "should be able to create a new key without provider_id column")
}
