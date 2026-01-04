package models

import (
	"gorm.io/gorm"
)

// Provider represents an LLM provider configuration (OpenAI, Anthropic, local endpoints)
type Provider struct {
	gorm.Model

	UserID       uint   `gorm:"not null;index" json:"user_id"`
	Name         string `gorm:"type:varchar(100);not null" json:"name"`
	ProviderType string `gorm:"type:varchar(50);not null" json:"provider_type"` // openai, anthropic, local
	BaseURL      string `gorm:"type:varchar(500)" json:"base_url"`
	APIKey       string `gorm:"type:varchar(500);not null" json:"-"` // Never expose in API responses

	IsActive        bool    `gorm:"default:true" json:"is_active"`
	DefaultModel    string  `gorm:"type:varchar(100)" json:"default_model"`
	InputCostPer1K  float64 `gorm:"column:input_cost_per_1k;default:0" json:"input_cost_per_1k"`
	OutputCostPer1K float64 `gorm:"column:output_cost_per_1k;default:0" json:"output_cost_per_1k"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ProviderType constants for supported LLM providers
const (
	ProviderTypeOpenAI    = "openai"
	ProviderTypeAnthropic = "anthropic"
	ProviderTypeLocal     = "local"
)

// DefaultBaseURLs for known providers
var DefaultBaseURLs = map[string]string{
	ProviderTypeOpenAI:    "https://api.openai.com",
	ProviderTypeAnthropic: "https://api.anthropic.com",
}

// GetBaseURL returns the provider's base URL, falling back to default if empty
func (p *Provider) GetBaseURL() string {
	if p.BaseURL != "" {
		return p.BaseURL
	}
	if defaultURL, ok := DefaultBaseURLs[p.ProviderType]; ok {
		return defaultURL
	}
	return ""
}
