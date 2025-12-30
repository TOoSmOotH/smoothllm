package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smoothweb/backend/internal/models"
	"gorm.io/gorm"
)

const (
	defaultTheme               = "warm-editorial"
	defaultRegistrationEnabled = true
	defaultAutoApproveNewUsers = true
)

var allowedThemes = map[string]bool{
	"warm-editorial":         true,
	"warm-editorial-dark":    true,
	"solarpunk":              true,
	"gothic-punk":            true,
	"cyberpunk":              true,
	"cassette-futurism":      true,
	"cassette-futurism-dark": true,
	"vaporware":              true,
}

type SettingsHandler struct {
	db *gorm.DB
}

func NewSettingsHandler(db *gorm.DB) *SettingsHandler {
	return &SettingsHandler{db: db}
}

// GetTheme returns the global theme for the app.
// GET /api/v1/settings/theme
func (h *SettingsHandler) GetTheme(c *gin.Context) {
	theme := defaultTheme
	var setting models.AppSetting
	if err := h.db.Where("key = ?", models.AppSettingThemeKey).First(&setting).Error; err == nil {
		if setting.Value != "" {
			theme = setting.Value
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"theme": theme,
		},
	})
}

// UpdateTheme updates the global theme for the app (admin only).
// PUT /api/v1/admin/settings/theme
func (h *SettingsHandler) UpdateTheme(c *gin.Context) {
	var req struct {
		Theme string `json:"theme" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !allowedThemes[req.Theme] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid theme"})
		return
	}

	var setting models.AppSetting
	err := h.db.Where("key = ?", models.AppSettingThemeKey).First(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			setting = models.AppSetting{
				Key:   models.AppSettingThemeKey,
				Value: req.Theme,
			}
			if err := h.db.Create(&setting).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save theme"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load settings"})
			return
		}
	} else {
		setting.Value = req.Theme
		if err := h.db.Save(&setting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update theme"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"theme": req.Theme,
		},
	})
}

// GetRegistrationSettings returns registration and approval settings (admin only).
// GET /api/v1/admin/settings/registration
func (h *SettingsHandler) GetRegistrationSettings(c *gin.Context) {
	registrationEnabled := h.getBoolSetting(models.AppSettingRegistrationEnabledKey, defaultRegistrationEnabled)
	autoApprove := h.getBoolSetting(models.AppSettingAutoApproveNewUsersKey, defaultAutoApproveNewUsers)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"registration_enabled":   registrationEnabled,
			"auto_approve_new_users": autoApprove,
		},
	})
}

// UpdateRegistrationSettings updates registration and approval settings (admin only).
// PUT /api/v1/admin/settings/registration
func (h *SettingsHandler) UpdateRegistrationSettings(c *gin.Context) {
	var req struct {
		RegistrationEnabled bool `json:"registration_enabled"`
		AutoApproveNewUsers bool `json:"auto_approve_new_users"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.setBoolSetting(models.AppSettingRegistrationEnabledKey, req.RegistrationEnabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update registration settings"})
		return
	}

	if err := h.setBoolSetting(models.AppSettingAutoApproveNewUsersKey, req.AutoApproveNewUsers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update registration settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"registration_enabled":   req.RegistrationEnabled,
			"auto_approve_new_users": req.AutoApproveNewUsers,
		},
	})
}

func (h *SettingsHandler) getBoolSetting(key string, defaultValue bool) bool {
	var setting models.AppSetting
	if err := h.db.Where("key = ?", key).First(&setting).Error; err == nil {
		if setting.Value == "true" {
			return true
		}
		if setting.Value == "false" {
			return false
		}
	}
	return defaultValue
}

func (h *SettingsHandler) setBoolSetting(key string, value bool) error {
	valueString := "false"
	if value {
		valueString = "true"
	}

	var setting models.AppSetting
	err := h.db.Where("key = ?", key).First(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			setting = models.AppSetting{
				Key:   key,
				Value: valueString,
			}
			return h.db.Create(&setting).Error
		}
		return err
	}

	setting.Value = valueString
	return h.db.Save(&setting).Error
}
