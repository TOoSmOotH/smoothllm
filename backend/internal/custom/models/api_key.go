package models

import (
	"time"

	"gorm.io/gorm"
)

// ProxyAPIKey represents a user-generated proxy API key that routes to a configured provider
type ProxyAPIKey struct {
	gorm.Model

	UserID    uint   `gorm:"not null;index" json:"user_id"`
	KeyHash   string `gorm:"type:varchar(255);uniqueIndex;not null" json:"-"` // Never expose hash
	KeyPrefix string `gorm:"type:varchar(20);not null" json:"key_prefix"`     // sk-smoothllm-xxx (visible part)
	Name      string `gorm:"type:varchar(100)" json:"name"`
	IsActive  bool   `gorm:"default:true" json:"is_active"`

	LastUsedAt *time.Time `json:"last_used_at"`
	ExpiresAt  *time.Time `json:"expires_at"`

	// Relationships
	AllowedProviders []KeyAllowedProvider `gorm:"foreignKey:ProxyAPIKeyID;constraint:OnDelete:CASCADE" json:"allowed_providers"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ProxyAPIKeyPrefix is the prefix for generated proxy API keys
const ProxyAPIKeyPrefix = "sk-smoothllm-"

// IsExpired checks if the API key has expired
func (k *ProxyAPIKey) IsExpired() bool {
	if k.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*k.ExpiresAt)
}

// IsValid checks if the API key is valid (active and not expired)
func (k *ProxyAPIKey) IsValid() bool {
	return k.IsActive && !k.IsExpired()
}

// UpdateLastUsed updates the LastUsedAt timestamp to the current time
func (k *ProxyAPIKey) UpdateLastUsed() {
	now := time.Now()
	k.LastUsedAt = &now
}
