package models

import "gorm.io/gorm"

type MediaFile struct {
	gorm.Model
	UserID uint `gorm:"index;not null" json:"user_id"`

	FileName      string `gorm:"type:varchar(255);not null" json:"file_name"`
	OriginalName  string `gorm:"type:varchar(255);not null" json:"original_name"`
	FilePath      string `gorm:"type:varchar(500);not null" json:"file_path"`
	ThumbnailPath string `gorm:"type:varchar(500)" json:"thumbnail_path,omitempty"`

	MimeType string `gorm:"type:varchar(100);not null" json:"mime_type"`
	FileSize int64  `gorm:"not null" json:"file_size"`
	Width    *int   `json:"width,omitempty"`
	Height   *int   `json:"height,omitempty"`

	FileType MediaType `gorm:"type:varchar(20);not null" json:"file_type"`
	IsPublic bool      `gorm:"default:true" json:"is_public"`

	StorageProvider string `gorm:"type:varchar(50);default:'local'" json:"storage_provider"`
	StoragePath     string `gorm:"type:varchar(500)" json:"storage_path,omitempty"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

type MediaType string

const (
	MediaTypeAvatar     MediaType = "avatar"
	MediaTypeCoverPhoto MediaType = "cover_photo"
	MediaTypeDocument   MediaType = "document"
	MediaTypeGallery    MediaType = "gallery"
)
