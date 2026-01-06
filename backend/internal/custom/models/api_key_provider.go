package models

import (
	"gorm.io/gorm"
)

// KeyAllowedProvider links a ProxyAPIKey to a Provider and specifies which models are allowed
type KeyAllowedProvider struct {
	gorm.Model

	ProxyAPIKeyID uint `gorm:"not null;index" json:"proxy_api_key_id"`
	ProviderID    uint `gorm:"not null;index" json:"provider_id"`

	// Models is a list of allowed model names for this key+provider combination
	// If empty, all models of the provider are allowed
	Models []string `gorm:"serializer:json" json:"models"`

	// Relationships
	Provider *Provider `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
}
