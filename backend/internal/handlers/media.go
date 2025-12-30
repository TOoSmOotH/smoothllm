package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smoothweb/backend/internal/auth"
	"github.com/smoothweb/backend/internal/models"
	"github.com/smoothweb/backend/internal/services"
	"gorm.io/gorm"
)

type MediaHandler struct {
	mediaService *services.MediaService
	db           *gorm.DB
}

func NewMediaHandler(mediaService *services.MediaService, db *gorm.DB) *MediaHandler {
	return &MediaHandler{
		mediaService: mediaService,
		db:           db,
	}
}

// UploadAvatar handles avatar upload with cropping support
func (h *MediaHandler) UploadAvatar(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get file from form
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}
	defer file.Close()

	// Create upload request
	req := &services.MediaUploadRequest{
		UserID: userID,
		File:   file,
		Header: header,
	}

	// Upload avatar
	response, err := h.mediaService.UploadAvatar(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "avatar uploaded successfully",
	})
}

// UploadCoverPhoto handles cover photo upload
func (h *MediaHandler) UploadCoverPhoto(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get file from form
	file, header, err := c.Request.FormFile("cover")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}
	defer file.Close()

	// Create upload request
	req := &services.MediaUploadRequest{
		UserID: userID,
		File:   file,
		Header: header,
	}

	// Upload cover photo
	response, err := h.mediaService.UploadCoverPhoto(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "cover photo uploaded successfully",
	})
}

// CropMedia crops an existing media file
func (h *MediaHandler) CropMedia(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse media ID from URL
	mediaIDStr := c.Param("id")
	mediaID, err := strconv.ParseUint(mediaIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid media ID"})
		return
	}

	// Parse crop parameters from request body
	var req services.CropRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set user ID and media ID from URL
	req.UserID = userID
	req.MediaID = uint(mediaID)

	// Crop the media
	response, err := h.mediaService.CropMedia(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "media cropped successfully",
	})
}

// DeleteMedia deletes a media file
func (h *MediaHandler) DeleteMedia(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse media ID from URL
	mediaIDStr := c.Param("id")
	mediaID, err := strconv.ParseUint(mediaIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid media ID"})
		return
	}

	// Delete the media
	err = h.mediaService.DeleteMedia(uint(mediaID), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "media deleted successfully",
	})
}

// ServeMedia serves uploaded files
func (h *MediaHandler) ServeMedia(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename is required"})
		return
	}

	// Get media by filename
	media, err := h.mediaService.GetMediaByFilename(filename)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}

	// Check if user has access to this media
	userID := auth.GetUserID(c)
	if !media.IsPublic && (userID == 0 || media.UserID != userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Serve the file
	c.File(media.FilePath)
}

// GetMedia retrieves media information
func (h *MediaHandler) GetMedia(c *gin.Context) {
	// Parse media ID from URL
	mediaIDStr := c.Param("id")
	mediaID, err := strconv.ParseUint(mediaIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid media ID"})
		return
	}

	// Get user ID from context if authenticated
	var userID *uint
	if authUserID := auth.GetUserID(c); authUserID != 0 {
		userID = &authUserID
	}

	// Get media
	media, err := h.mediaService.GetMedia(uint(mediaID), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    media,
	})
}

// GetUserMedia retrieves all media files for a user
func (h *MediaHandler) GetUserMedia(c *gin.Context) {
	// Parse user ID from URL
	userIDStr := c.Param("userId")
	targetUserID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Get current user ID
	currentUserID := auth.GetUserID(c)

	// Users can only view their own media
	if currentUserID == 0 || uint(targetUserID) != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Parse query parameters
	fileType := c.Query("type")
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Get user's media files
	var media []models.MediaFile
	query := h.db.Where("user_id = ?", targetUserID)
	
	if fileType != "" {
		query = query.Where("file_type = ?", fileType)
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&media).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve media"})
		return
	}

	// Get total count
	var total int64
	h.db.Model(&models.MediaFile{}).Where("user_id = ?", targetUserID).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"media": media,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		},
	})
}