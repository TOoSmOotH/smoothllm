package custom

import (
	"github.com/smoothweb/backend/internal/custom/models"
	"gorm.io/gorm"
)

// AutoMigrate runs database migrations for custom models
func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.Provider{},
		&models.ProxyAPIKey{},
		&models.KeyAllowedProvider{},
		&models.UsageRecord{},
	); err != nil {
		return err
	}

	// Run data migration for existing keys if needed
	if err := migrateExistingData(db); err != nil {
		return err
	}

	return nil
}

// migrateExistingData handles data migration from single-provider keys to multi-provider keys
func migrateExistingData(db *gorm.DB) error {
	// 1. Migrate ProxyAPIKey.ProviderID to KeyAllowedProvider
	if db.Migrator().HasColumn(&models.ProxyAPIKey{}, "provider_id") {
		var keys []struct {
			ID         uint
			ProviderID uint
		}

		// Find keys that still use the old provider_id column
		if err := db.Table("proxy_api_keys").
			Select("id, provider_id").
			Where("provider_id IS NOT NULL AND provider_id > 0").
			Scan(&keys).Error; err == nil {

			for _, k := range keys {
				// Check if allowed provider record already exists
				var count int64
				db.Model(&models.KeyAllowedProvider{}).
					Where("proxy_api_key_id = ? AND provider_id = ?", k.ID, k.ProviderID).
					Count(&count)

				if count == 0 {
					ap := models.KeyAllowedProvider{
						ProxyAPIKeyID: k.ID,
						ProviderID:    k.ProviderID,
					}
					db.Create(&ap)
				}
			}
		}

		// 2. Drop the column and its foreign key constraint
		// GORM's DropColumn is unreliable with foreign keys in SQLite.
		// We use a manual reconstruction approach.
		err := db.Transaction(func(tx *gorm.DB) error {
			// Rename old table
			if err := tx.Migrator().RenameTable("proxy_api_keys", "proxy_api_keys_old"); err != nil {
				return err
			}

			// Create new table via AutoMigrate
			if err := tx.AutoMigrate(&models.ProxyAPIKey{}); err != nil {
				return err
			}

			// Copy data from old to new (omitting provider_id)
			// We list all columns explicitly to be safe
			if err := tx.Exec(`
				INSERT INTO proxy_api_keys (id, created_at, updated_at, deleted_at, user_id, key_hash, key_prefix, name, is_active, last_used_at, expires_at)
				SELECT id, created_at, updated_at, deleted_at, user_id, key_hash, key_prefix, name, is_active, last_used_at, expires_at
				FROM proxy_api_keys_old
			`).Error; err != nil {
				return err
			}

			// Drop old table
			if err := tx.Migrator().DropTable("proxy_api_keys_old"); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}
