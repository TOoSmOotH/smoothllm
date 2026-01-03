package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/smoothweb/backend/internal/models"
)

type ProfileService struct {
	db *gorm.DB
}

func NewProfileService(db *gorm.DB) *ProfileService {
	return &ProfileService{db: db}
}

// ProfileResponse represents the profile data returned to clients
type ProfileResponse struct {
	ID           uint                   `json:"id"`
	UserID       uint                   `json:"user_id"`
	Username     string                 `json:"username"`
	FirstName    *string                `json:"first_name,omitempty"`
	LastName     *string                `json:"last_name,omitempty"`
	DisplayName  *string                `json:"display_name,omitempty"`
	Avatar       *models.MediaFile      `json:"avatar,omitempty"`
	CoverPhoto   *models.MediaFile      `json:"cover_photo,omitempty"`
	Phone        *string                `json:"phone,omitempty"`
	Website      *string                `json:"website,omitempty"`
	Bio          *string                `json:"bio,omitempty"`
	Location     *string                `json:"location,omitempty"`
	City         *string                `json:"city,omitempty"`
	State        *string                `json:"state,omitempty"`
	Country      *string                `json:"country,omitempty"`
	Timezone     *string                `json:"timezone,omitempty"`
	Birthday     *time.Time             `json:"birthday,omitempty"`
	Gender       *string                `json:"gender,omitempty"`
	Pronouns     *string                `json:"pronouns,omitempty"`
	Language     *string                `json:"language,omitempty"`
	JobTitle     *string                `json:"job_title,omitempty"`
	Company      *string                `json:"company,omitempty"`
	Industry     *string                `json:"industry,omitempty"`
	LinkedInURL  *string                `json:"linkedin_url,omitempty"`
	PortfolioURL *string                `json:"portfolio_url,omitempty"`
	Interests    []string               `json:"interests,omitempty"`
	Skills       []string               `json:"skills,omitempty"`
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// UpdateProfileRequest represents the request to update a profile
type UpdateProfileRequest struct {
	FirstName    *string                `json:"first_name,omitempty"`
	LastName     *string                `json:"last_name,omitempty"`
	DisplayName  *string                `json:"display_name,omitempty"`
	Phone        *string                `json:"phone,omitempty"`
	Website      *string                `json:"website,omitempty"`
	Bio          *string                `json:"bio,omitempty"`
	Location     *string                `json:"location,omitempty"`
	City         *string                `json:"city,omitempty"`
	State        *string                `json:"state,omitempty"`
	Country      *string                `json:"country,omitempty"`
	Timezone     *string                `json:"timezone,omitempty"`
	Birthday     *time.Time             `json:"birthday,omitempty"`
	Gender       *string                `json:"gender,omitempty"`
	Pronouns     *string                `json:"pronouns,omitempty"`
	Language     *string                `json:"language,omitempty"`
	JobTitle     *string                `json:"job_title,omitempty"`
	Company      *string                `json:"company,omitempty"`
	Industry     *string                `json:"industry,omitempty"`
	LinkedInURL  *string                `json:"linkedin_url,omitempty"`
	PortfolioURL *string                `json:"portfolio_url,omitempty"`
	Interests    []string               `json:"interests,omitempty"`
	Skills       []string               `json:"skills,omitempty"`
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
}

// GetProfileByUsername retrieves a user's profile by username with privacy filtering
func (s *ProfileService) GetProfileByUsername(username string, viewerID *uint) (*ProfileResponse, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return s.getProfileByUserID(user.ID, viewerID)
}

// GetProfileByUserID retrieves a user's profile by user ID with privacy filtering
func (s *ProfileService) GetProfileByUserID(userID uint, viewerID *uint) (*ProfileResponse, error) {
	return s.getProfileByUserID(userID, viewerID)
}

// getProfileByUserID internal method to get profile with privacy filtering
func (s *ProfileService) getProfileByUserID(userID uint, viewerID *uint) (*ProfileResponse, error) {
	var profile models.UserProfile
	if err := s.db.Preload("Avatar").Preload("CoverPhoto").Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create a default profile if it doesn't exist
			if err := s.createDefaultProfile(userID); err != nil {
				return nil, fmt.Errorf("failed to create default profile: %w", err)
			}
			// Try again after creating
			if err := s.db.Preload("Avatar").Preload("CoverPhoto").Where("user_id = ?", userID).First(&profile).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// Get user info for username
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	// Get privacy settings
	var privacySettings models.PrivacySettings
	if err := s.db.Where("user_id = ?", userID).First(&privacySettings).Error; err != nil {
		// Create default privacy settings if they don't exist
		privacySettings = models.PrivacySettings{
			UserID:             userID,
			ProfileVisibility:  string(models.PrivacyLevelPublic),
			EmailVisibility:    string(models.PrivacyLevelPrivate),
			PhoneVisibility:    string(models.PrivacyLevelPrivate),
			BirthdayVisibility: string(models.PrivacyLevelConnections),
			LocationVisibility: string(models.PrivacyLevelConnections),
		}
		s.db.Create(&privacySettings)
	}

	// Apply privacy filtering
	response := s.filterProfileByPrivacy(&profile, &user, &privacySettings, viewerID)

	return response, nil
}

// UpdateProfile updates the authenticated user's profile
func (s *ProfileService) UpdateProfile(userID uint, req *UpdateProfileRequest) (*ProfileResponse, error) {
	// Validate input
	if err := s.validateProfileUpdate(req); err != nil {
		return nil, err
	}

	var profile models.UserProfile
	if err := s.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create profile if it doesn't exist
			profile = models.UserProfile{UserID: userID}
			if err := s.db.Create(&profile).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// Update fields
	updates := make(map[string]interface{})
	if req.FirstName != nil {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = req.LastName
	}
	if req.DisplayName != nil {
		updates["display_name"] = req.DisplayName
	}
	if req.Phone != nil {
		updates["phone"] = req.Phone
	}
	if req.Website != nil {
		updates["website"] = req.Website
	}
	if req.Bio != nil {
		updates["bio"] = req.Bio
	}
	if req.Location != nil {
		updates["location"] = req.Location
	}
	if req.City != nil {
		updates["city"] = req.City
	}
	if req.State != nil {
		updates["state"] = req.State
	}
	if req.Country != nil {
		updates["country"] = req.Country
	}
	if req.Timezone != nil {
		updates["timezone"] = req.Timezone
	}
	if req.Birthday != nil {
		updates["birthday"] = req.Birthday
	}
	if req.Gender != nil {
		updates["gender"] = req.Gender
	}
	if req.Pronouns != nil {
		updates["pronouns"] = req.Pronouns
	}
	if req.Language != nil {
		updates["language"] = req.Language
	}
	if req.JobTitle != nil {
		updates["job_title"] = req.JobTitle
	}
	if req.Company != nil {
		updates["company"] = req.Company
	}
	if req.Industry != nil {
		updates["industry"] = req.Industry
	}
	if req.LinkedInURL != nil {
		updates["linked_in_url"] = req.LinkedInURL
	}
	if req.PortfolioURL != nil {
		updates["portfolio_url"] = req.PortfolioURL
	}
	if req.Interests != nil {
		interestsJSON, _ := json.Marshal(req.Interests)
		updates["interests"] = datatypes.JSON(interestsJSON)
	}
	if req.Skills != nil {
		skillsJSON, _ := json.Marshal(req.Skills)
		updates["skills"] = datatypes.JSON(skillsJSON)
	}
	if req.CustomFields != nil {
		customFieldsJSON, _ := json.Marshal(req.CustomFields)
		updates["custom_fields"] = datatypes.JSON(customFieldsJSON)
	}

	if err := s.db.Model(&profile).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Refresh profile data
	if err := s.db.Preload("Avatar").Preload("CoverPhoto").Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}

	// Get user info
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	// Get privacy settings
	var privacySettings models.PrivacySettings
	s.db.Where("user_id = ?", userID).First(&privacySettings)

	// Update profile completion score
	go s.updateProfileCompletion(userID)

	// Return full profile (no privacy filtering for owner)
	response := s.buildProfileResponse(&profile, &user)

	return response, nil
}

// CreateProfile creates a new profile for a user
func (s *ProfileService) CreateProfile(userID uint) (*ProfileResponse, error) {
	// Check if profile already exists
	var existingProfile models.UserProfile
	if err := s.db.Where("user_id = ?", userID).First(&existingProfile).Error; err == nil {
		return nil, fmt.Errorf("profile already exists")
	}

	// Create default profile
	profile := models.UserProfile{UserID: userID}
	if err := s.db.Create(&profile).Error; err != nil {
		return nil, err
	}

	// Create default privacy settings
	privacySettings := models.PrivacySettings{
		UserID:             userID,
		ProfileVisibility:  string(models.PrivacyLevelPublic),
		EmailVisibility:    string(models.PrivacyLevelPrivate),
		PhoneVisibility:    string(models.PrivacyLevelPrivate),
		BirthdayVisibility: string(models.PrivacyLevelConnections),
		LocationVisibility: string(models.PrivacyLevelConnections),
	}
	s.db.Create(&privacySettings)

	// Get user info
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	// Update profile completion score
	go s.updateProfileCompletion(userID)

	// Return profile
	response := s.buildProfileResponse(&profile, &user)

	return response, nil
}

// CheckUsernameAvailability checks if a username is available
func (s *ProfileService) CheckUsernameAvailability(username string) (bool, error) {
	var count int64
	if err := s.db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}

// DeleteProfile soft deletes the user's profile
func (s *ProfileService) DeleteProfile(userID uint) error {
	var profile models.UserProfile
	if err := s.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("profile not found")
		}
		return err
	}

	if err := s.db.Delete(&profile).Error; err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	return nil
}

// createDefaultProfile creates a default profile for a user
func (s *ProfileService) createDefaultProfile(userID uint) error {
	profile := models.UserProfile{UserID: userID}
	return s.db.Create(&profile).Error
}

// filterProfileByPrivacy applies privacy settings to profile data
func (s *ProfileService) filterProfileByPrivacy(profile *models.UserProfile, user *models.User, privacy *models.PrivacySettings, viewerID *uint) *ProfileResponse {
	response := s.buildProfileResponse(profile, user)

	// If viewer is the owner, return full profile
	if viewerID != nil && *viewerID == profile.UserID {
		return response
	}

	// Check profile visibility
	if !s.canViewField(privacy.ProfileVisibility, viewerID, profile.UserID) {
		// Return minimal public info
		return &ProfileResponse{
			ID:       profile.ID,
			UserID:   profile.UserID,
			Username: user.Username,
		}
	}

	// Apply field-level privacy
	if !s.canViewField(privacy.PhoneVisibility, viewerID, profile.UserID) {
		response.Phone = nil
	}
	if !s.canViewField(privacy.BirthdayVisibility, viewerID, profile.UserID) {
		response.Birthday = nil
	}
	if !s.canViewField(privacy.LocationVisibility, viewerID, profile.UserID) {
		response.Location = nil
		response.City = nil
		response.State = nil
		response.Country = nil
	}

	return response
}

// canViewField determines if a field can be viewed based on privacy settings
func (s *ProfileService) canViewField(visibility string, viewerID *uint, profileUserID uint) bool {
	switch visibility {
	case string(models.PrivacyLevelPublic):
		return true
	case string(models.PrivacyLevelRegistered):
		return viewerID != nil
	case string(models.PrivacyLevelConnections):
		// For now, only allow self to see connections-level fields
		// In future, implement actual connections logic
		return viewerID != nil && *viewerID == profileUserID
	case string(models.PrivacyLevelPrivate):
		return viewerID != nil && *viewerID == profileUserID
	default:
		return false
	}
}

// buildProfileResponse creates a ProfileResponse from models
func (s *ProfileService) buildProfileResponse(profile *models.UserProfile, user *models.User) *ProfileResponse {
	response := &ProfileResponse{
		ID:           profile.ID,
		UserID:       profile.UserID,
		Username:     user.Username,
		FirstName:    profile.FirstName,
		LastName:     profile.LastName,
		DisplayName:  profile.DisplayName,
		Avatar:       profile.Avatar,
		CoverPhoto:   profile.CoverPhoto,
		Phone:        profile.Phone,
		Website:      profile.Website,
		Bio:          profile.Bio,
		Location:     profile.Location,
		City:         profile.City,
		State:        profile.State,
		Country:      profile.Country,
		Timezone:     profile.Timezone,
		Birthday:     profile.Birthday,
		Gender:       profile.Gender,
		Pronouns:     profile.Pronouns,
		Language:     profile.Language,
		JobTitle:     profile.JobTitle,
		Company:      profile.Company,
		Industry:     profile.Industry,
		LinkedInURL:  profile.LinkedInURL,
		PortfolioURL: profile.PortfolioURL,
		CreatedAt:    profile.CreatedAt,
		UpdatedAt:    profile.UpdatedAt,
	}

	// Parse JSON fields
	if profile.Interests != nil {
		var interests []string
		json.Unmarshal(profile.Interests, &interests)
		response.Interests = interests
	}
	if profile.Skills != nil {
		var skills []string
		json.Unmarshal(profile.Skills, &skills)
		response.Skills = skills
	}
	if profile.CustomFields != nil {
		var customFields map[string]interface{}
		json.Unmarshal(profile.CustomFields, &customFields)
		response.CustomFields = customFields
	}

	return response
}

// validateProfileUpdate validates the profile update request
func (s *ProfileService) validateProfileUpdate(req *UpdateProfileRequest) error {
	// Validate website URL
	if req.Website != nil && *req.Website != "" {
		if !s.isValidURL(*req.Website) {
			return fmt.Errorf("invalid website URL")
		}
	}

	// Validate LinkedIn URL
	if req.LinkedInURL != nil && *req.LinkedInURL != "" {
		if !strings.Contains(*req.LinkedInURL, "linkedin.com") {
			return fmt.Errorf("invalid LinkedIn URL")
		}
		if !s.isValidURL(*req.LinkedInURL) {
			return fmt.Errorf("invalid LinkedIn URL format")
		}
	}

	// Validate portfolio URL
	if req.PortfolioURL != nil && *req.PortfolioURL != "" {
		if !s.isValidURL(*req.PortfolioURL) {
			return fmt.Errorf("invalid portfolio URL")
		}
	}

	// Validate phone number
	if req.Phone != nil && *req.Phone != "" {
		if !s.isValidPhone(*req.Phone) {
			return fmt.Errorf("invalid phone number format")
		}
	}

	return nil
}

// isValidURL checks if a string is a valid URL
func (s *ProfileService) isValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// isValidPhone checks if a string is a valid phone number
func (s *ProfileService) isValidPhone(phone string) bool {
	// Simple phone validation - allows digits, spaces, +, -, (, )
	phoneRegex := regexp.MustCompile(`^[\d\s\-\+\(\)]+$`)
	return phoneRegex.MatchString(phone) && len(strings.ReplaceAll(phone, " ", "")) >= 10
}

// updateProfileCompletion calculates and updates the profile completion score
func (s *ProfileService) updateProfileCompletion(userID uint) error {
	var profile models.UserProfile
	if err := s.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return err
	}

	// Calculate completion score
	score := s.calculateCompletionScore(&profile)

	// Update or create completion record
	var completion models.ProfileCompletion
	if err := s.db.Where("user_id = ?", userID).First(&completion).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			completion = models.ProfileCompletion{
				UserID:        userID,
				Score:         score,
				LastUpdatedAt: time.Now(),
			}
			return s.db.Create(&completion).Error
		}
		return err
	}

	completion.Score = score
	completion.LastUpdatedAt = time.Now()
	return s.db.Save(&completion).Error
}

// calculateCompletionScore calculates the completion score for a profile
func (s *ProfileService) calculateCompletionScore(profile *models.UserProfile) int {
	score := 0

	// Basic info (40 points)
	if profile.FirstName != nil && *profile.FirstName != "" {
		score += 10
	}
	if profile.LastName != nil && *profile.LastName != "" {
		score += 10
	}
	if profile.DisplayName != nil && *profile.DisplayName != "" {
		score += 10
	}
	if profile.Bio != nil && *profile.Bio != "" {
		score += 10
	}

	// Contact info (20 points)
	if profile.Phone != nil && *profile.Phone != "" {
		score += 5
	}
	if profile.Website != nil && *profile.Website != "" {
		score += 5
	}
	if profile.Location != nil && *profile.Location != "" {
		score += 10
	}

	// Personal info (20 points)
	if profile.Birthday != nil {
		score += 5
	}
	if profile.Gender != nil && *profile.Gender != "" {
		score += 5
	}
	if profile.Pronouns != nil && *profile.Pronouns != "" {
		score += 5
	}
	if profile.Language != nil && *profile.Language != "" {
		score += 5
	}

	// Professional info (20 points)
	if profile.JobTitle != nil && *profile.JobTitle != "" {
		score += 5
	}
	if profile.Company != nil && *profile.Company != "" {
		score += 5
	}
	if profile.LinkedInURL != nil && *profile.LinkedInURL != "" {
		score += 5
	}
	if profile.PortfolioURL != nil && *profile.PortfolioURL != "" {
		score += 5
	}

	return score
}
