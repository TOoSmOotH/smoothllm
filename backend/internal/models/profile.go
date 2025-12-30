package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UserProfile struct {
	gorm.Model
	UserID uint `gorm:"uniqueIndex;not null" json:"user_id"`

	FirstName   *string `gorm:"type:varchar(100)" json:"first_name,omitempty"`
	LastName    *string `gorm:"type:varchar(100)" json:"last_name,omitempty"`
	DisplayName *string `gorm:"type:varchar(100)" json:"display_name,omitempty"`

	AvatarID     *uint `gorm:"index" json:"avatar_id,omitempty"`
	CoverPhotoID *uint `gorm:"index" json:"cover_photo_id,omitempty"`

	Phone   *string `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Website *string `gorm:"type:varchar(255)" json:"website,omitempty"`

	Bio      *string `gorm:"type:text" json:"bio,omitempty"`
	Location *string `gorm:"type:varchar(255)" json:"location,omitempty"`
	City     *string `gorm:"type:varchar(100)" json:"city,omitempty"`
	State    *string `gorm:"type:varchar(100)" json:"state,omitempty"`
	Country  *string `gorm:"type:varchar(100)" json:"country,omitempty"`
	Timezone *string `gorm:"type:varchar(50)" json:"timezone,omitempty"`

	Birthday *time.Time `json:"birthday,omitempty"`
	Gender   *string    `gorm:"type:varchar(50)" json:"gender,omitempty"`
	Pronouns *string    `gorm:"type:varchar(50)" json:"pronouns,omitempty"`
	Language *string    `gorm:"type:varchar(10);default:'en'" json:"language,omitempty"`

	JobTitle     *string `gorm:"type:varchar(100)" json:"job_title,omitempty"`
	Company      *string `gorm:"type:varchar(255)" json:"company,omitempty"`
	Industry     *string `gorm:"type:varchar(100)" json:"industry,omitempty"`
	LinkedInURL  *string `gorm:"type:varchar(255)" json:"linkedin_url,omitempty"`
	PortfolioURL *string `gorm:"type:varchar(255)" json:"portfolio_url,omitempty"`

	Interests datatypes.JSON `gorm:"type:json" json:"interests,omitempty"`
	Skills    datatypes.JSON `gorm:"type:json" json:"skills,omitempty"`

	CustomFields datatypes.JSON `gorm:"type:json" json:"custom_fields,omitempty"`

	User       User       `gorm:"foreignKey:UserID" json:"-"`
	Avatar     *MediaFile `gorm:"foreignKey:AvatarID" json:"avatar,omitempty"`
	CoverPhoto *MediaFile `gorm:"foreignKey:CoverPhotoID" json:"cover_photo,omitempty"`
}
