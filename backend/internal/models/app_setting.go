package models

import "time"

type AppSetting struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"uniqueIndex;size:64" json:"key"`
	Value     string    `gorm:"size:255" json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const (
	AppSettingThemeKey               = "theme"
	AppSettingRegistrationEnabledKey = "registration_enabled"
	AppSettingAutoApproveNewUsersKey = "auto_approve_new_users"
)
