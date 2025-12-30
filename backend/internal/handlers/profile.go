package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smoothweb/backend/internal/auth"
	"github.com/smoothweb/backend/internal/services"
)

type ProfileHandler struct {
	profileService *services.ProfileService
}

func NewProfileHandler(profileService *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

// GetProfile retrieves a user's profile by username (public endpoint)
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	// Get viewer ID from context if authenticated
	var viewerID *uint
	if userID := auth.GetUserID(c); userID != 0 {
		viewerID = &userID
	}

	profile, err := h.profileService.GetProfileByUsername(username, viewerID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
	})
}

// GetMyProfile retrieves the authenticated user's own profile
func (h *ProfileHandler) GetMyProfile(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	profile, err := h.profileService.GetProfileByUserID(userID, &userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
	})
}

// UpdateProfile updates the authenticated user's profile
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req services.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.profileService.UpdateProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
		"message": "profile updated successfully",
	})
}

// CreateProfile creates a new profile for the authenticated user
func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	profile, err := h.profileService.CreateProfile(userID)
	if err != nil {
		if err.Error() == "profile already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "profile already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create profile"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    profile,
		"message": "profile created successfully",
	})
}

// CheckUsernameAvailability checks if a username is available
func (h *ProfileHandler) CheckUsernameAvailability(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username query parameter is required"})
		return
	}

	available, err := h.profileService.CheckUsernameAvailability(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check username availability"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"available": available,
	})
}

// GetProfileByID retrieves a user's profile by user ID (public endpoint)
func (h *ProfileHandler) GetProfileByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Get viewer ID from context if authenticated
	var viewerID *uint
	if authUserID := auth.GetUserID(c); authUserID != 0 {
		viewerID = &authUserID
	}

	profile, err := h.profileService.GetProfileByUserID(uint(userID), viewerID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
	})
}

// DeleteProfile soft deletes the authenticated user's profile
func (h *ProfileHandler) DeleteProfile(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.profileService.DeleteProfile(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "profile deleted successfully",
	})
}
