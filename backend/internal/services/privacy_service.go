package services

import (
	"errors"
	"fmt"

	"github.com/smoothweb/backend/internal/models"
	"gorm.io/gorm"
)

// PrivacyService handles privacy settings operations
type PrivacyService struct {
	db *gorm.DB
}

// PrivacySettingsUpdate represents privacy settings update request
type PrivacySettingsUpdate struct {
	ProfileVisibility    *string `json:"profile_visibility"`
	EmailVisibility      *string `json:"email_visibility"`
	PhoneVisibility      *string `json:"phone_visibility"`
	BirthdayVisibility   *string `json:"birthday_visibility"`
	LocationVisibility   *string `json:"location_visibility"`
	LastActiveVisibility *string `json:"last_active_visibility"`

	ShowEmail           *bool `json:"show_email"`
	ShowPhone           *bool `json:"show_phone"`
	ShowSocialLinks     *bool `json:"show_social_links"`
	AllowDirectMessages *bool `json:"allow_direct_messages"`
	ShowOnlineStatus    *bool `json:"show_online_status"`

	AppearInSearch   *bool `json:"appear_in_search"`
	SuggestToFriends *bool `json:"suggest_to_friends"`

	AllowDataSharing *bool `json:"allow_data_sharing"`
	AnalyticsEnabled *bool `json:"analytics_enabled"`
}

// NewPrivacyService creates a new privacy service instance
func NewPrivacyService(db *gorm.DB) *PrivacyService {
	return &PrivacyService{db: db}
}

// GetPrivacySettings retrieves privacy settings for a user
func (s *PrivacyService) GetPrivacySettings(userID uint) (*models.PrivacySettings, error) {
	var settings models.PrivacySettings
	err := s.db.Where("user_id = ?", userID).First(&settings).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create default privacy settings for the user
			defaultSettings := s.createDefaultSettings(userID)
			if err := s.db.Create(&defaultSettings).Error; err != nil {
				return nil, fmt.Errorf("failed to create default privacy settings: %w", err)
			}
			return &defaultSettings, nil
		}
		return nil, fmt.Errorf("failed to get privacy settings: %w", err)
	}

	return &settings, nil
}

// UpdatePrivacySettings updates privacy settings for a user
func (s *PrivacyService) UpdatePrivacySettings(userID uint, update *PrivacySettingsUpdate) (*models.PrivacySettings, error) {
	// Get existing settings or create defaults
	var existing models.PrivacySettings
	err := s.db.Where("user_id = ?", userID).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create default settings for the user
			existing = s.createDefaultSettings(userID)
			if err := s.db.Create(&existing).Error; err != nil {
				return nil, fmt.Errorf("failed to create privacy settings: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get privacy settings: %w", err)
		}
	}

	// Apply partial updates
	if update.ProfileVisibility != nil {
		existing.ProfileVisibility = *update.ProfileVisibility
	}
	if update.EmailVisibility != nil {
		existing.EmailVisibility = *update.EmailVisibility
	}
	if update.PhoneVisibility != nil {
		existing.PhoneVisibility = *update.PhoneVisibility
	}
	if update.BirthdayVisibility != nil {
		existing.BirthdayVisibility = *update.BirthdayVisibility
	}
	if update.LocationVisibility != nil {
		existing.LocationVisibility = *update.LocationVisibility
	}
	if update.LastActiveVisibility != nil {
		existing.LastActiveVisibility = *update.LastActiveVisibility
	}
	if update.ShowEmail != nil {
		existing.ShowEmail = *update.ShowEmail
	}
	if update.ShowPhone != nil {
		existing.ShowPhone = *update.ShowPhone
	}
	if update.ShowSocialLinks != nil {
		existing.ShowSocialLinks = *update.ShowSocialLinks
	}
	if update.AllowDirectMessages != nil {
		existing.AllowDirectMessages = *update.AllowDirectMessages
	}
	if update.ShowOnlineStatus != nil {
		existing.ShowOnlineStatus = *update.ShowOnlineStatus
	}
	if update.AppearInSearch != nil {
		existing.AppearInSearch = *update.AppearInSearch
	}
	if update.SuggestToFriends != nil {
		existing.SuggestToFriends = *update.SuggestToFriends
	}
	if update.AllowDataSharing != nil {
		existing.AllowDataSharing = *update.AllowDataSharing
	}
	if update.AnalyticsEnabled != nil {
		existing.AnalyticsEnabled = *update.AnalyticsEnabled
	}

	// Save the updated settings
	if err := s.db.Save(&existing).Error; err != nil {
		return nil, fmt.Errorf("failed to update privacy settings: %w", err)
	}

	return &existing, nil
}

// ApplyPrivacyPreset applies a preset to user's privacy settings
func (s *PrivacyService) ApplyPrivacyPreset(userID uint, preset string) (*models.PrivacySettings, error) {
	var settings models.PrivacySettings

	// Get existing settings or create defaults
	err := s.db.Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			settings = s.createDefaultSettings(userID)
		} else {
			return nil, fmt.Errorf("failed to get privacy settings: %w", err)
		}
	}

	// Apply preset values
	switch preset {
	case "public":
		settings.ProfileVisibility = string(models.PrivacyLevelPublic)
		settings.EmailVisibility = string(models.PrivacyLevelPublic)
		settings.PhoneVisibility = string(models.PrivacyLevelPublic)
		settings.BirthdayVisibility = string(models.PrivacyLevelPublic)
		settings.LocationVisibility = string(models.PrivacyLevelPublic)
		settings.LastActiveVisibility = string(models.PrivacyLevelPublic)
		settings.ShowEmail = true
		settings.ShowPhone = true
		settings.ShowSocialLinks = true
		settings.AllowDirectMessages = true
		settings.ShowOnlineStatus = true
		settings.AppearInSearch = true
		settings.SuggestToFriends = true
		settings.AllowDataSharing = true
		settings.AnalyticsEnabled = true

	case "connections":
		settings.ProfileVisibility = string(models.PrivacyLevelConnections)
		settings.EmailVisibility = string(models.PrivacyLevelPrivate)
		settings.PhoneVisibility = string(models.PrivacyLevelPrivate)
		settings.BirthdayVisibility = string(models.PrivacyLevelConnections)
		settings.LocationVisibility = string(models.PrivacyLevelConnections)
		settings.LastActiveVisibility = string(models.PrivacyLevelConnections)
		settings.ShowEmail = false
		settings.ShowPhone = false
		settings.ShowSocialLinks = true
		settings.AllowDirectMessages = true
		settings.ShowOnlineStatus = true
		settings.AppearInSearch = true
		settings.SuggestToFriends = true
		settings.AllowDataSharing = true
		settings.AnalyticsEnabled = true

	case "private":
		settings.ProfileVisibility = string(models.PrivacyLevelPrivate)
		settings.EmailVisibility = string(models.PrivacyLevelPrivate)
		settings.PhoneVisibility = string(models.PrivacyLevelPrivate)
		settings.BirthdayVisibility = string(models.PrivacyLevelPrivate)
		settings.LocationVisibility = string(models.PrivacyLevelPrivate)
		settings.LastActiveVisibility = string(models.PrivacyLevelPrivate)
		settings.ShowEmail = false
		settings.ShowPhone = false
		settings.ShowSocialLinks = false
		settings.AllowDirectMessages = false
		settings.ShowOnlineStatus = false
		settings.AppearInSearch = false
		settings.SuggestToFriends = false
		settings.AllowDataSharing = false
		settings.AnalyticsEnabled = false

	case "professional":
		settings.ProfileVisibility = string(models.PrivacyLevelConnections)
		settings.EmailVisibility = string(models.PrivacyLevelPrivate)
		settings.PhoneVisibility = string(models.PrivacyLevelPrivate)
		settings.BirthdayVisibility = string(models.PrivacyLevelPrivate)
		settings.LocationVisibility = string(models.PrivacyLevelPrivate)
		settings.LastActiveVisibility = string(models.PrivacyLevelConnections)
		settings.ShowEmail = false
		settings.ShowPhone = false
		settings.ShowSocialLinks = true
		settings.AllowDirectMessages = true
		settings.ShowOnlineStatus = true
		settings.AppearInSearch = true
		settings.SuggestToFriends = true
		settings.AllowDataSharing = true
		settings.AnalyticsEnabled = true

	default:
		return nil, errors.New("invalid privacy preset")
	}

	// Save the settings
	if settings.ID == 0 {
		if err := s.db.Create(&settings).Error; err != nil {
			return nil, fmt.Errorf("failed to create privacy settings: %w", err)
		}
	} else {
		if err := s.db.Save(&settings).Error; err != nil {
			return nil, fmt.Errorf("failed to update privacy settings: %w", err)
		}
	}

	return &settings, nil
}

// GetPrivacyPresets returns all available privacy presets
func (s *PrivacyService) GetPrivacyPresets() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "public",
			"label":       "Public",
			"description": "All information is visible to everyone. Analytics enabled.",
			"settings": map[string]interface{}{
				"profile_visibility":     string(models.PrivacyLevelPublic),
				"email_visibility":       string(models.PrivacyLevelPublic),
				"phone_visibility":       string(models.PrivacyLevelPublic),
				"birthday_visibility":    string(models.PrivacyLevelPublic),
				"location_visibility":    string(models.PrivacyLevelPublic),
				"last_active_visibility": string(models.PrivacyLevelPublic),
				"show_email":             true,
				"show_phone":             true,
				"show_social_links":      true,
				"allow_direct_messages":  true,
				"show_online_status":     true,
				"appear_in_search":       true,
				"suggest_to_friends":     true,
				"allow_data_sharing":     true,
				"analytics_enabled":      true,
			},
		},
		{
			"name":        "connections",
			"label":       "Connections",
			"description": "Profile is visible to connections only. Limited analytics.",
			"settings": map[string]interface{}{
				"profile_visibility":     string(models.PrivacyLevelConnections),
				"email_visibility":       string(models.PrivacyLevelPrivate),
				"phone_visibility":       string(models.PrivacyLevelPrivate),
				"birthday_visibility":    string(models.PrivacyLevelConnections),
				"location_visibility":    string(models.PrivacyLevelConnections),
				"last_active_visibility": string(models.PrivacyLevelConnections),
				"show_email":             false,
				"show_phone":             false,
				"show_social_links":      true,
				"allow_direct_messages":  true,
				"show_online_status":     true,
				"appear_in_search":       true,
				"suggest_to_friends":     true,
				"allow_data_sharing":     true,
				"analytics_enabled":      true,
			},
		},
		{
			"name":        "private",
			"label":       "Private",
			"description": "Profile is private. Only connections see you. All other settings disabled.",
			"settings": map[string]interface{}{
				"profile_visibility":     string(models.PrivacyLevelPrivate),
				"email_visibility":       string(models.PrivacyLevelPrivate),
				"phone_visibility":       string(models.PrivacyLevelPrivate),
				"birthday_visibility":    string(models.PrivacyLevelPrivate),
				"location_visibility":    string(models.PrivacyLevelPrivate),
				"last_active_visibility": string(models.PrivacyLevelPrivate),
				"show_email":             false,
				"show_phone":             false,
				"show_social_links":      false,
				"allow_direct_messages":  false,
				"show_online_status":     false,
				"appear_in_search":       false,
				"suggest_to_friends":     false,
				"allow_data_sharing":     false,
				"analytics_enabled":      false,
			},
		},
		{
			"name":        "professional",
			"label":       "Professional",
			"description": "Email and phone hidden. Only connections can see profile.",
			"settings": map[string]interface{}{
				"profile_visibility":     string(models.PrivacyLevelConnections),
				"email_visibility":       string(models.PrivacyLevelPrivate),
				"phone_visibility":       string(models.PrivacyLevelPrivate),
				"birthday_visibility":    string(models.PrivacyLevelPrivate),
				"location_visibility":    string(models.PrivacyLevelPrivate),
				"last_active_visibility": string(models.PrivacyLevelConnections),
				"show_email":             false,
				"show_phone":             false,
				"show_social_links":      true,
				"allow_direct_messages":  true,
				"show_online_status":     true,
				"appear_in_search":       true,
				"suggest_to_friends":     true,
				"allow_data_sharing":     true,
				"analytics_enabled":      true,
			},
		},
	}
}

// createDefaultSettings creates default privacy settings for a user
func (s *PrivacyService) createDefaultSettings(userID uint) models.PrivacySettings {
	return models.PrivacySettings{
		UserID:               userID,
		ProfileVisibility:    string(models.PrivacyLevelPublic),
		EmailVisibility:      string(models.PrivacyLevelPrivate),
		PhoneVisibility:      string(models.PrivacyLevelPrivate),
		BirthdayVisibility:   string(models.PrivacyLevelConnections),
		LocationVisibility:   string(models.PrivacyLevelConnections),
		LastActiveVisibility: string(models.PrivacyLevelConnections),
		ShowEmail:            false,
		ShowPhone:            false,
		ShowSocialLinks:      true,
		AllowDirectMessages:  true,
		ShowOnlineStatus:     true,
		AppearInSearch:       true,
		SuggestToFriends:     true,
		AllowDataSharing:     false,
		AnalyticsEnabled:     true,
	}
}

// PrivacySettingsUpdateToModel converts update request to model
func PrivacySettingsUpdateToModel(update *PrivacySettingsUpdate) *models.PrivacySettings {
	settings := &models.PrivacySettings{}

	if update.ProfileVisibility != nil {
		settings.ProfileVisibility = *update.ProfileVisibility
	}
	if update.EmailVisibility != nil {
		settings.EmailVisibility = *update.EmailVisibility
	}
	if update.PhoneVisibility != nil {
		settings.PhoneVisibility = *update.PhoneVisibility
	}
	if update.BirthdayVisibility != nil {
		settings.BirthdayVisibility = *update.BirthdayVisibility
	}
	if update.LocationVisibility != nil {
		settings.LocationVisibility = *update.LocationVisibility
	}
	if update.LastActiveVisibility != nil {
		settings.LastActiveVisibility = *update.LastActiveVisibility
	}
	if update.ShowEmail != nil {
		settings.ShowEmail = *update.ShowEmail
	}
	if update.ShowPhone != nil {
		settings.ShowPhone = *update.ShowPhone
	}
	if update.ShowSocialLinks != nil {
		settings.ShowSocialLinks = *update.ShowSocialLinks
	}
	if update.AllowDirectMessages != nil {
		settings.AllowDirectMessages = *update.AllowDirectMessages
	}
	if update.ShowOnlineStatus != nil {
		settings.ShowOnlineStatus = *update.ShowOnlineStatus
	}
	if update.AppearInSearch != nil {
		settings.AppearInSearch = *update.AppearInSearch
	}
	if update.SuggestToFriends != nil {
		settings.SuggestToFriends = *update.SuggestToFriends
	}
	if update.AllowDataSharing != nil {
		settings.AllowDataSharing = *update.AllowDataSharing
	}
	if update.AnalyticsEnabled != nil {
		settings.AnalyticsEnabled = *update.AnalyticsEnabled
	}

	return settings
}
