package custom

import (
	"github.com/smoothweb/backend/internal/custom/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.Provider{},
		&models.ProxyAPIKey{},
		&models.KeyAllowedProvider{},
		&models.UsageRecord{},
	); err != nil {
		return err
	}

	// 1. Handle usage_records schema inconsistencies (e.g. model vs model_name)
	if err := migrateUsageRecords(db); err != nil {
		return err
	}

	// 2. Run data migration for existing keys if needed
	if err := migrateExistingData(db); err != nil {
		return err
	}

	return nil
}

// migrateUsageRecords handles column renaming and additions for the usage_records table in SQLite
func migrateUsageRecords(db *gorm.DB) error {
	// If the table doesn't exist yet, return
	if !db.Migrator().HasTable("usage_records") {
		return nil
	}

	// SQLite doesn't support RENAME COLUMN, so if we have model_name, we need to handle it.
	// Check if 'model' column exists. If not, check if 'model_name' exists and rename.
	if !db.Migrator().HasColumn(&models.UsageRecord{}, "model") {
		if db.Migrator().HasColumn("usage_records", "model_name") {
			// Try to rename column (only works in recent SQLite/GORM versions)
			if err := db.Migrator().RenameColumn("usage_records", "model_name", "model"); err != nil {
				// Fallback: This is safer for some SQLite versions
				return db.Transaction(func(tx *gorm.DB) error {
					if err := tx.Migrator().RenameTable("usage_records", "usage_records_old"); err != nil {
						return err
					}
					if err := tx.AutoMigrate(&models.UsageRecord{}); err != nil {
						return err
					}
					// Copy data, mapping model_name to model
					// Also handle status_code, cost, etc. if they were added later
					return tx.Exec(`
						INSERT INTO usage_records (id, created_at, updated_at, deleted_at, user_id, proxy_key_id, provider_id, model, input_tokens, output_tokens, total_tokens, cost, request_duration, status_code, error_message)
						SELECT id, created_at, updated_at, deleted_at, user_id, proxy_key_id, provider_id, model_name, input_tokens, output_tokens, total_tokens, cost, request_duration, status_code, error_message
						FROM usage_records_old
					`).Error
				})
			}
		}
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

			// Drop indexes on old table to prevent conflict with new table indexes in SQLite
			// SQLite index names are database-wide, and GORM AutoMigrate will try to create them on the new table.
			_ = tx.Migrator().DropIndex("proxy_api_keys_old", "idx_proxy_api_keys_key_hash")
			_ = tx.Migrator().DropIndex("proxy_api_keys_old", "idx_proxy_api_keys_user_id")

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
