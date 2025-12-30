package models

import (
	"time"
	"gorm.io/gorm"
	"gorm.io/datatypes"
)

type SocialLink struct {
	gorm.Model
	UserID     uint   `gorm:"index;not null" json:"user_id"`
	Platform   string `gorm:"type:varchar(50);not null" json:"platform"`
	Username   string `gorm:"type:varchar(100)" json:"username"`
	URL        string `gorm:"type:varchar(500);not null" json:"url"`
	IsVerified bool   `gorm:"default:false" json:"is_verified"`
	IsPrimary  bool   `gorm:"default:false" json:"is_primary"`

	OAuthProvider string     `gorm:"type:varchar(50)" json:"oauth_provider,omitempty"`
	OAuthID       string     `gorm:"type:varchar(255)" json:"oauth_id,omitempty"`
	LinkedAt      *time.Time `json:"linked_at,omitempty"`

	DisplayOrder int  `gorm:"default:0" json:"display_order"`
	IsPublic     bool `gorm:"default:true" json:"is_public"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

type OAuthAccount struct {
	gorm.Model
	UserID         uint       `gorm:"index;not null" json:"user_id"`
	Provider       string     `gorm:"type:varchar(50);not null" json:"provider"`
	ProviderID     string     `gorm:"type:varchar(255);not null" json:"provider_id"`
	AccessToken    string     `gorm:"type:text" json:"-"`
	RefreshToken   string     `gorm:"type:text" json:"-"`
	TokenExpiresAt *time.Time `json:"token_expires_at,omitempty"`

	ProviderProfile datatypes.JSON `gorm:"type:json" json:"provider_profile,omitempty"`

	SocialLinkID *uint `gorm:"index" json:"social_link_id,omitempty"`

	User       User        `gorm:"foreignKey:UserID" json:"-"`
	SocialLink *SocialLink `gorm:"foreignKey:SocialLinkID" json:"social_link,omitempty"`
}
