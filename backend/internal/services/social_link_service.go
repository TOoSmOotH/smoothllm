package services

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/smoothweb/backend/internal/models"
	"gorm.io/gorm"
)

type SocialLinkService struct {
	db *gorm.DB
}

func NewSocialLinkService(db *gorm.DB) *SocialLinkService {
	return &SocialLinkService{db: db}
}

// AddSocialLinkRequest represents request to add a social link
type AddSocialLinkRequest struct {
	Platform  string `json:"platform" binding:"required"`
	Username  string `json:"username,omitempty"`
	URL       string `json:"url" binding:"required"`
	IsPrimary bool   `json:"is_primary,omitempty"`
	IsPublic  bool   `json:"is_public,omitempty"`
}

// UpdateSocialLinkRequest represents request to update a social link
type UpdateSocialLinkRequest struct {
	Platform  *string `json:"platform,omitempty"`
	Username  *string `json:"username,omitempty"`
	URL       *string `json:"url,omitempty"`
	IsPrimary *bool   `json:"is_primary,omitempty"`
	IsPublic  *bool   `json:"is_public,omitempty"`
}

// ReorderSocialLinksRequest represents request to reorder social links
type ReorderSocialLinksRequest struct {
	Links []LinkOrder `json:"links" binding:"required"`
}

type LinkOrder struct {
	ID    uint `json:"id" binding:"required"`
	Order int  `json:"order" binding:"required"`
}

// SocialLinkResponse represents social link data returned to clients
type SocialLinkResponse struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	Platform     string `json:"platform"`
	Username     string `json:"username"`
	URL          string `json:"url"`
	IsVerified   bool   `json:"is_verified"`
	IsPrimary    bool   `json:"is_primary"`
	IsPublic     bool   `json:"is_public"`
	DisplayOrder int    `json:"display_order"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// GetSocialLinks retrieves all social links for a user
func (s *SocialLinkService) GetSocialLinks(userID uint, viewerID *uint) ([]SocialLinkResponse, error) {
	var links []models.SocialLink
	query := s.db.Where("user_id = ?", userID)

	// Filter by privacy if viewer is not the owner
	if viewerID == nil || *viewerID != userID {
		query = query.Where("is_public = ?", true)
	}

	if err := query.Order("display_order ASC, created_at DESC").Find(&links).Error; err != nil {
		return nil, fmt.Errorf("failed to get social links: %w", err)
	}

	return s.buildSocialLinkResponses(links), nil
}

// GetSocialLinksByUserID retrieves public social links for a specific user
func (s *SocialLinkService) GetSocialLinksByUserID(targetUserID uint, viewerID *uint) ([]SocialLinkResponse, error) {
	var links []models.SocialLink
	query := s.db.Where("user_id = ?", targetUserID)

	// Filter by privacy
	if viewerID == nil || *viewerID != targetUserID {
		query = query.Where("is_public = ?", true)
	}

	if err := query.Order("display_order ASC, created_at DESC").Find(&links).Error; err != nil {
		return nil, fmt.Errorf("failed to get social links: %w", err)
	}

	return s.buildSocialLinkResponses(links), nil
}

// GetSocialLinkByID retrieves a specific social link
func (s *SocialLinkService) GetSocialLinkByID(linkID uint, userID uint) (*SocialLinkResponse, error) {
	var link models.SocialLink
	if err := s.db.Where("id = ? AND user_id = ?", linkID, userID).First(&link).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("social link not found")
		}
		return nil, fmt.Errorf("failed to get social link: %w", err)
	}

	response := s.buildSocialLinkResponse(&link)
	return &response, nil
}

// AddSocialLink adds a new social link for a user
func (s *SocialLinkService) AddSocialLink(userID uint, req *AddSocialLinkRequest) (*SocialLinkResponse, error) {
	// Validate request
	if err := s.validateSocialLinkRequest(req); err != nil {
		return nil, err
	}

	// Check if link with same platform and URL already exists
	var existingCount int64
	if err := s.db.Model(&models.SocialLink{}).
		Where("user_id = ? AND platform = ? AND url = ?", userID, req.Platform, req.URL).
		Count(&existingCount).Error; err != nil {
		return nil, err
	}
	if existingCount > 0 {
		return nil, fmt.Errorf("social link already exists")
	}

	// If setting as primary, remove primary flag from other links
	if req.IsPrimary {
		s.db.Model(&models.SocialLink{}).
			Where("user_id = ? AND is_primary = ?", userID, true).
			Update("is_primary", false)
	}

	// Determine display order (max + 1)
	var maxOrder int
	s.db.Model(&models.SocialLink{}).
		Where("user_id = ?", userID).
		Select("COALESCE(MAX(display_order), 0)").
		Scan(&maxOrder)

	link := models.SocialLink{
		UserID:       userID,
		Platform:     req.Platform,
		Username:     req.Username,
		URL:          req.URL,
		IsPrimary:    req.IsPrimary,
		IsPublic:     req.IsPublic,
		DisplayOrder: maxOrder + 1,
	}

	if err := s.db.Create(&link).Error; err != nil {
		return nil, fmt.Errorf("failed to create social link: %w", err)
	}

	response := s.buildSocialLinkResponse(&link)
	return &response, nil
}

// UpdateSocialLink updates an existing social link
func (s *SocialLinkService) UpdateSocialLink(linkID uint, userID uint, req *UpdateSocialLinkRequest) (*SocialLinkResponse, error) {
	var link models.SocialLink
	if err := s.db.Where("id = ? AND user_id = ?", linkID, userID).First(&link).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("social link not found")
		}
		return nil, fmt.Errorf("failed to find social link: %w", err)
	}

	// Validate URL if provided
	if req.URL != nil {
		if !s.isValidURL(*req.URL) {
			return nil, fmt.Errorf("invalid URL format")
		}
	}

	// If setting as primary, remove primary flag from other links
	if req.IsPrimary != nil && *req.IsPrimary {
		s.db.Model(&models.SocialLink{}).
			Where("user_id = ? AND id != ? AND is_primary = ?", userID, linkID, true).
			Update("is_primary", false)
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.Platform != nil {
		updates["platform"] = req.Platform
	}
	if req.Username != nil {
		updates["username"] = req.Username
	}
	if req.URL != nil {
		updates["url"] = req.URL
	}
	if req.IsPrimary != nil {
		updates["is_primary"] = req.IsPrimary
	}
	if req.IsPublic != nil {
		updates["is_public"] = req.IsPublic
	}

	if err := s.db.Model(&link).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update social link: %w", err)
	}

	// Refresh data
	s.db.First(&link, linkID)
	response := s.buildSocialLinkResponse(&link)
	return &response, nil
}

// DeleteSocialLink deletes a social link
func (s *SocialLinkService) DeleteSocialLink(linkID uint, userID uint) error {
	var link models.SocialLink
	if err := s.db.Where("id = ? AND user_id = ?", linkID, userID).First(&link).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("social link not found")
		}
		return fmt.Errorf("failed to find social link: %w", err)
	}

	if err := s.db.Delete(&link).Error; err != nil {
		return fmt.Errorf("failed to delete social link: %w", err)
	}

	return nil
}

// ReorderSocialLinks reorders social links for a user
func (s *SocialLinkService) ReorderSocialLinks(userID uint, req *ReorderSocialLinksRequest) ([]SocialLinkResponse, error) {
	// Validate all links belong to user
	linkIDs := make([]uint, len(req.Links))
	for i, linkOrder := range req.Links {
		linkIDs[i] = linkOrder.ID

		var link models.SocialLink
		if err := s.db.Where("id = ? AND user_id = ?", linkOrder.ID, userID).First(&link).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("social link %d not found or not owned by user", linkOrder.ID)
			}
			return nil, fmt.Errorf("failed to find social link: %w", err)
		}
	}

	// Update display order for all links
	for _, linkOrder := range req.Links {
		if err := s.db.Model(&models.SocialLink{}).
			Where("id = ? AND user_id = ?", linkOrder.ID, userID).
			Update("display_order", linkOrder.Order).Error; err != nil {
			return nil, fmt.Errorf("failed to reorder social link: %w", err)
		}
	}

	// Return updated links
	var links []models.SocialLink
	if err := s.db.Where("user_id = ?", userID).Order("display_order ASC").Find(&links).Error; err != nil {
		return nil, fmt.Errorf("failed to get updated social links: %w", err)
	}

	return s.buildSocialLinkResponses(links), nil
}

// validateSocialLinkRequest validates social link request
func (s *SocialLinkService) validateSocialLinkRequest(req *AddSocialLinkRequest) error {
	// Validate platform
	if req.Platform == "" {
		return fmt.Errorf("platform is required")
	}
	if !s.isValidPlatform(req.Platform) {
		return fmt.Errorf("invalid platform: %s", req.Platform)
	}

	// Validate URL
	if req.URL == "" {
		return fmt.Errorf("URL is required")
	}
	if !s.isValidURL(req.URL) {
		return fmt.Errorf("invalid URL format")
	}

	return nil
}

// isValidPlatform checks if platform is supported
func (s *SocialLinkService) isValidPlatform(platform string) bool {
	supportedPlatforms := []string{
		"github", "twitter", "linkedin", "instagram",
		"facebook", "youtube", "tiktok", "discord",
		"telegram", "reddit", "stackoverflow", "medium",
		"dribbble", "behance", "portfolio", "website",
	}

	platformLower := strings.ToLower(platform)
	for _, p := range supportedPlatforms {
		if p == platformLower {
			return true
		}
	}
	return false
}

// isValidURL checks if URL is valid
func (s *SocialLinkService) isValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// buildSocialLinkResponse builds response from model
func (s *SocialLinkService) buildSocialLinkResponse(link *models.SocialLink) SocialLinkResponse {
	return SocialLinkResponse{
		ID:           link.ID,
		UserID:       link.UserID,
		Platform:     link.Platform,
		Username:     link.Username,
		URL:          link.URL,
		IsVerified:   link.IsVerified,
		IsPrimary:    link.IsPrimary,
		IsPublic:     link.IsPublic,
		DisplayOrder: link.DisplayOrder,
		CreatedAt:    link.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    link.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// buildSocialLinkResponses builds responses from models
func (s *SocialLinkService) buildSocialLinkResponses(links []models.SocialLink) []SocialLinkResponse {
	responses := make([]SocialLinkResponse, len(links))
	for i, link := range links {
		responses[i] = s.buildSocialLinkResponse(&link)
	}
	return responses
}
