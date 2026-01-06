package models

import (
	"time"

	"gorm.io/gorm"
)

// Provider represents an LLM provider configuration (OpenAI, Anthropic, local endpoints)
type Provider struct {
	gorm.Model

	UserID       uint   `gorm:"not null;index" json:"user_id"`
	Name         string `gorm:"type:varchar(100);not null" json:"name"`
	ProviderType string `gorm:"type:varchar(50);not null" json:"provider_type"` // openai, anthropic, anthropic_max, local
	BaseURL      string `gorm:"type:varchar(500)" json:"base_url"`
	APIKey       string `gorm:"type:varchar(500)" json:"-"` // Never expose in API responses (not required for OAuth providers)

	// OAuth fields for Claude Max subscription
	RefreshToken   string     `gorm:"type:varchar(500)" json:"-"`          // OAuth refresh token (never expose)
	AccessToken    string     `gorm:"type:varchar(500)" json:"-"`          // OAuth access token (never expose)
	TokenExpiresAt *time.Time `gorm:"column:token_expires_at" json:"-"`    // When access token expires
	OAuthConnected bool       `gorm:"default:false" json:"oauth_connected"` // Whether OAuth is connected

	IsActive             bool    `gorm:"default:true" json:"is_active"`
	DefaultModel         string  `gorm:"type:varchar(100)" json:"default_model"`
	InputCostPerMillion  float64 `gorm:"column:input_cost_per_million;default:0" json:"input_cost_per_million"`
	OutputCostPerMillion float64 `gorm:"column:output_cost_per_million;default:0" json:"output_cost_per_million"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ProviderType constants for supported LLM providers
const (
	ProviderTypeOpenAI       = "openai"
	ProviderTypeAnthropic    = "anthropic"
	ProviderTypeAnthropicMax = "anthropic_max" // Claude Max subscription via OAuth
	ProviderTypeVLLM         = "vllm"
	ProviderTypeLocal        = "local"
	ProviderTypeZai          = "zai"
)

// DefaultBaseURLs for known providers
var DefaultBaseURLs = map[string]string{
	ProviderTypeOpenAI:       "https://api.openai.com",
	ProviderTypeAnthropic:    "https://api.anthropic.com",
	ProviderTypeAnthropicMax: "https://api.anthropic.com",
	ProviderTypeZai:          "https://api.z.ai/api/paas/v4/",
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

// IsOAuthProvider returns true if this provider uses OAuth authentication
func (p *Provider) IsOAuthProvider() bool {
	return p.ProviderType == ProviderTypeAnthropicMax
}

// IsTokenExpired returns true if the OAuth access token has expired or will expire soon
func (p *Provider) IsTokenExpired() bool {
	if p.TokenExpiresAt == nil {
		return true
	}
	// Consider token expired if it expires within 5 minutes
	return time.Now().Add(5 * time.Minute).After(*p.TokenExpiresAt)
}

// NeedsTokenRefresh returns true if the provider needs a token refresh
func (p *Provider) NeedsTokenRefresh() bool {
	return p.IsOAuthProvider() && p.IsTokenExpired() && p.RefreshToken != ""
}
