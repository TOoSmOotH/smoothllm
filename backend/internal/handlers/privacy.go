package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smoothweb/backend/internal/services"
)

// PrivacyHandler handles privacy settings HTTP requests
type PrivacyHandler struct {
	privacyService *services.PrivacyService
}

// NewPrivacyHandler creates a new privacy handler instance
func NewPrivacyHandler(privacyService *services.PrivacyService) *PrivacyHandler {
	return &PrivacyHandler{privacyService: privacyService}
}

// GetPrivacySettings returns the current user's privacy settings
// GET /api/v1/privacy
func (h *PrivacyHandler) GetPrivacySettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	settings, err := h.privacyService.GetPrivacySettings(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get privacy settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": settings})
}

// UpdatePrivacySettings updates the current user's privacy settings
// PUT /api/v1/privacy
func (h *PrivacyHandler) UpdatePrivacySettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var settings services.PrivacySettingsUpdate
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updated, err := h.privacyService.UpdatePrivacySettings(userID.(uint), &settings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update privacy settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updated})
}

// ApplyPrivacyPreset applies a preset to the user's privacy settings
// POST /api/v1/privacy/preset/:name
func (h *PrivacyHandler) ApplyPrivacyPreset(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	preset := c.Param("name")
	if preset == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Preset name is required"})
		return
	}

	settings, err := h.privacyService.ApplyPrivacyPreset(userID.(uint), preset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid privacy preset"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": settings})
}

// GetPrivacyPresets returns all available privacy presets
// GET /api/v1/privacy/presets
func (h *PrivacyHandler) GetPrivacyPresets(c *gin.Context) {
	presets := h.privacyService.GetPrivacyPresets()
	c.JSON(http.StatusOK, gin.H{"data": presets})
}

// GetPrivacySettingsByAdmin returns privacy settings for a specific user (admin only)
// GET /api/v1/admin/users/:userId/privacy
func (h *PrivacyHandler) GetPrivacySettingsByAdmin(c *gin.Context) {
	userIdStr := c.Param("userId")
	if userIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	userId, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	settings, err := h.privacyService.GetPrivacySettings(uint(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get privacy settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": settings})
}

// UpdatePrivacySettingsByAdmin updates privacy settings for a specific user (admin only)
// PUT /api/v1/admin/users/:userId/privacy
func (h *PrivacyHandler) UpdatePrivacySettingsByAdmin(c *gin.Context) {
	userIdStr := c.Param("userId")
	if userIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	userId, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var settings services.PrivacySettingsUpdate
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updated, err := h.privacyService.UpdatePrivacySettings(uint(userId), &settings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update privacy settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updated})
}
