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
	db, err := gorm.Open(sqlite.Open(":memory:?_foreign_keys=on"), &gorm.Config{})
	require.NoError(t, err)

	// Manually create the table with the old column, foreign key constraint, and indexes
	err = db.Exec(`
		CREATE TABLE proxy_api_keys (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			user_id INTEGER NOT NULL,
			key_hash TEXT NOT NULL,
			key_prefix TEXT NOT NULL,
			name TEXT,
			is_active BOOLEAN DEFAULT 1,
			last_used_at DATETIME,
			expires_at DATETIME,
			provider_id INTEGER NOT NULL,
			CONSTRAINT fk_proxy_api_keys_provider FOREIGN KEY (provider_id) REFERENCES providers(id) ON DELETE CASCADE
		)
	`).Error
	require.NoError(t, err)

	err = db.Exec("CREATE UNIQUE INDEX idx_proxy_api_keys_key_hash ON proxy_api_keys(key_hash)").Error
	require.NoError(t, err)
	err = db.Exec("CREATE INDEX idx_proxy_api_keys_user_id ON proxy_api_keys(user_id)").Error
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Provider{}, &models.KeyAllowedProvider{})
	require.NoError(t, err)

	// Verify the column exists initially
	assert.True(t, db.Migrator().HasColumn("proxy_api_keys", "provider_id"), "provider_id column should exist initially")

	// Insert test data
	provider := models.Provider{UserID: 1, Name: "Test", ProviderType: "openai", APIKey: "key", IsActive: true}
	err = db.Create(&provider).Error
	require.NoError(t, err)

	err = db.Exec("INSERT INTO proxy_api_keys (user_id, key_hash, key_prefix, provider_id, created_at, updated_at) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))", 1, "hash", "sk-", provider.ID).Error
	require.NoError(t, err)

	// 2. Run our custom migration logic
	err = migrateExistingData(db)
	require.NoError(t, err)

	// 3. Verify the column is gone
	assert.False(t, db.Migrator().HasColumn("proxy_api_keys", "provider_id"), "provider_id column should have been dropped")

	// 4. Verify data was migrated to KeyAllowedProvider
	var count int64
	err = db.Model(&models.KeyAllowedProvider{}).Where("provider_id = ?", provider.ID).Count(&count).Error
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
	assert.NoError(t, err, "should be able to create a new key without provider_id column and with reconstructed indexes")
}
