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

		// 2. Drop the column after migration to prevent NOT NULL constraint failures on new inserts
		if err := db.Migrator().DropColumn(&models.ProxyAPIKey{}, "provider_id"); err != nil {
			return err
		}
	}

	return nil
}
