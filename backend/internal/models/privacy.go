package models

import "gorm.io/gorm"

type PrivacySettings struct {
	gorm.Model
	UserID uint `gorm:"uniqueIndex;not null" json:"user_id"`

	ProfileVisibility string `gorm:"type:varchar(20);default:'public'" json:"profile_visibility"`

	EmailVisibility      string `gorm:"type:varchar(20);default:'private'" json:"email_visibility"`
	PhoneVisibility      string `gorm:"type:varchar(20);default:'private'" json:"phone_visibility"`
	BirthdayVisibility   string `gorm:"type:varchar(20);default:'connections'" json:"birthday_visibility"`
	LocationVisibility   string `gorm:"type:varchar(20);default:'connections'" json:"location_visibility"`
	LastActiveVisibility string `gorm:"type:varchar(20);default:'connections'" json:"last_active_visibility"`

	ShowEmail           bool `gorm:"default:false" json:"show_email"`
	ShowPhone           bool `gorm:"default:false" json:"show_phone"`
	ShowSocialLinks     bool `gorm:"default:true" json:"show_social_links"`
	AllowDirectMessages bool `gorm:"default:true" json:"allow_direct_messages"`
	ShowOnlineStatus    bool `gorm:"default:true" json:"show_online_status"`

	AppearInSearch   bool `gorm:"default:true" json:"appear_in_search"`
	SuggestToFriends bool `gorm:"default:true" json:"suggest_to_friends"`

	AllowDataSharing bool `gorm:"default:false" json:"allow_data_sharing"`
	AnalyticsEnabled bool `gorm:"default:true" json:"analytics_enabled"`
}

type PrivacyLevel string

const (
	PrivacyLevelPublic      PrivacyLevel = "public"
	PrivacyLevelRegistered  PrivacyLevel = "registered"
	PrivacyLevelConnections PrivacyLevel = "connections"
	PrivacyLevelPrivate     PrivacyLevel = "private"
)
