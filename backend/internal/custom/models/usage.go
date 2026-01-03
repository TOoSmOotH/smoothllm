package models

import (
	"gorm.io/gorm"
)

// UsageRecord tracks individual LLM API request usage for billing and analytics
type UsageRecord struct {
	gorm.Model

	UserID     uint   `gorm:"not null;index" json:"user_id"`
	ProxyKeyID uint   `gorm:"not null;index" json:"proxy_key_id"`
	ProviderID uint   `gorm:"not null;index" json:"provider_id"`
	Model      string `gorm:"type:varchar(100);index" json:"model"`

	InputTokens  int `gorm:"default:0" json:"input_tokens"`
	OutputTokens int `gorm:"default:0" json:"output_tokens"`
	TotalTokens  int `gorm:"default:0" json:"total_tokens"`

	Cost            float64 `gorm:"default:0" json:"cost"`
	RequestDuration int     `gorm:"default:0" json:"request_duration"` // milliseconds
	StatusCode      int     `gorm:"default:0" json:"status_code"`
	ErrorMessage    string  `gorm:"type:text" json:"error_message,omitempty"`

	// Relationships
	ProxyKey *ProxyAPIKey `gorm:"foreignKey:ProxyKeyID;constraint:OnDelete:CASCADE" json:"proxy_key,omitempty"`
	Provider *Provider    `gorm:"foreignKey:ProviderID;constraint:OnDelete:CASCADE" json:"provider,omitempty"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// CalculateCost computes the cost based on token usage and provider rates
func (u *UsageRecord) CalculateCost(inputCostPer1K, outputCostPer1K float64) float64 {
	inputCost := float64(u.InputTokens) / 1000.0 * inputCostPer1K
	outputCost := float64(u.OutputTokens) / 1000.0 * outputCostPer1K
	return inputCost + outputCost
}

// IsError checks if the request resulted in an error
func (u *UsageRecord) IsError() bool {
	return u.StatusCode >= 400 || u.ErrorMessage != ""
}

// IsSuccess checks if the request was successful
func (u *UsageRecord) IsSuccess() bool {
	return u.StatusCode >= 200 && u.StatusCode < 300 && u.ErrorMessage == ""
}
