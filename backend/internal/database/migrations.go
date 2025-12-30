package database

import "github.com/smoothweb/backend/internal/models"

func (d *Database) AutoMigrate() error {
	if err := d.DB.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.PrivacySettings{},
		&models.AppSetting{},
		&models.SocialLink{},
		&models.OAuthAccount{},
		&models.MediaFile{},
		&models.ProfileCompletion{},
		&models.CompletionMilestone{},
		&models.CompletionAchievement{},
	); err != nil {
		return err
	}

	return d.DB.Model(&models.User{}).
		Where("status IS NULL OR status = ''").
		Update("status", "active").Error
}

func (d *Database) Seed() error {
	var userCount int64
	if err := d.DB.Model(&models.User{}).Count(&userCount).Error; err != nil {
		return err
	}

	if userCount == 0 {
		return nil
	}

	return nil
}
