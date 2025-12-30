package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smoothweb/backend/internal/auth"
	"github.com/smoothweb/backend/internal/services"
)

type SocialLinkHandler struct {
	socialLinkService *services.SocialLinkService
}

func NewSocialLinkHandler(socialLinkService *services.SocialLinkService) *SocialLinkHandler {
	return &SocialLinkHandler{
		socialLinkService: socialLinkService,
	}
}

// GetSocialLinks retrieves authenticated user's social links
// GET /api/v1/social
func (h *SocialLinkHandler) GetSocialLinks(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	links, err := h.socialLinkService.GetSocialLinks(userID, &userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get social links"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    links,
	})
}

// GetUserSocialLinks retrieves a user's public social links
// GET /api/v1/social/:user_id
func (h *SocialLinkHandler) GetUserSocialLinks(c *gin.Context) {
	userIDStr := c.Param("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Get viewer ID from context if authenticated
	var viewerID *uint
	if authUserID := auth.GetUserID(c); authUserID != 0 {
		viewerID = &authUserID
	}

	links, err := h.socialLinkService.GetSocialLinksByUserID(uint(userID), viewerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get social links"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    links,
	})
}

// GetSocialLink retrieves a specific social link
// GET /api/v1/social/link/:id
func (h *SocialLinkHandler) GetSocialLink(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	linkIDStr := c.Param("id")
	if linkIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "link ID is required"})
		return
	}

	linkID, err := strconv.ParseUint(linkIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid link ID"})
		return
	}

	link, err := h.socialLinkService.GetSocialLinkByID(uint(linkID), userID)
	if err != nil {
		if err.Error() == "social link not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get social link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    link,
	})
}

// AddSocialLink adds a new social link
// POST /api/v1/social
func (h *SocialLinkHandler) AddSocialLink(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req services.AddSocialLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	link, err := h.socialLinkService.AddSocialLink(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    link,
		"message": "social link added successfully",
	})
}

// UpdateSocialLink updates an existing social link
// PUT /api/v1/social/:id
func (h *SocialLinkHandler) UpdateSocialLink(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	linkIDStr := c.Param("id")
	if linkIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "link ID is required"})
		return
	}

	linkID, err := strconv.ParseUint(linkIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid link ID"})
		return
	}

	var req services.UpdateSocialLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	link, err := h.socialLinkService.UpdateSocialLink(uint(linkID), userID, &req)
	if err != nil {
		if err.Error() == "social link not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    link,
		"message": "social link updated successfully",
	})
}

// DeleteSocialLink deletes a social link
// DELETE /api/v1/social/:id
func (h *SocialLinkHandler) DeleteSocialLink(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	linkIDStr := c.Param("id")
	if linkIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "link ID is required"})
		return
	}

	linkID, err := strconv.ParseUint(linkIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid link ID"})
		return
	}

	err = h.socialLinkService.DeleteSocialLink(uint(linkID), userID)
	if err != nil {
		if err.Error() == "social link not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "social link deleted successfully",
	})
}

// ReorderSocialLinks reorders social links
// PUT /api/v1/social/reorder
func (h *SocialLinkHandler) ReorderSocialLinks(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req services.ReorderSocialLinksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	links, err := h.socialLinkService.ReorderSocialLinks(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    links,
		"message": "social links reordered successfully",
	})
}
