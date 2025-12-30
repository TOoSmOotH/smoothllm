package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Email        string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	Username     string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Role         string `gorm:"type:varchar(20);default:'user'" json:"role"`
	Status       string `gorm:"type:varchar(20);default:'active'" json:"status"`

	Profile         *UserProfile     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"profile,omitempty"`
	PrivacySettings *PrivacySettings `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"privacy_settings,omitempty"`
	SocialLinks     []SocialLink     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"social_links,omitempty"`
	OAuthAccounts   []OAuthAccount   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"oauth_accounts,omitempty"`

	CompletionScore int            `gorm:"default:0" json:"completion_score"`
	CompletedFields datatypes.JSON `gorm:"type:json" json:"completed_fields"`

	LastActiveAt    *time.Time `json:"last_active_at"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
