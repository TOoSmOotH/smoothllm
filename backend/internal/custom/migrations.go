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
		&models.UsageRecord{},
	); err != nil {
		return err
	}

	return nil
}
