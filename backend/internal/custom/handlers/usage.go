package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smoothweb/backend/internal/auth"
	"github.com/smoothweb/backend/internal/custom/services"
)

// UsageHandler handles usage statistics endpoints
type UsageHandler struct {
	usageService *services.UsageService
}

// NewUsageHandler creates a new UsageHandler instance
func NewUsageHandler(usageService *services.UsageService) *UsageHandler {
	return &UsageHandler{
		usageService: usageService,
	}
}

// parseQueryParams extracts common query parameters for filtering usage data
func (h *UsageHandler) parseQueryParams(c *gin.Context) *services.UsageQueryParams {
	params := &services.UsageQueryParams{}

	// Parse start_date
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if t, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			params.StartDate = &t
		} else if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			params.StartDate = &t
		}
	}

	// Parse end_date
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if t, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			params.EndDate = &t
		} else if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// Set to end of day
			t = t.Add(24*time.Hour - time.Second)
			params.EndDate = &t
		}
	}

	// Parse provider_id
	if providerIDStr := c.Query("provider_id"); providerIDStr != "" {
		if id, err := strconv.ParseUint(providerIDStr, 10, 32); err == nil {
			providerID := uint(id)
			params.ProviderID = &providerID
		}
	}

	// Parse key_id
	if keyIDStr := c.Query("key_id"); keyIDStr != "" {
		if id, err := strconv.ParseUint(keyIDStr, 10, 32); err == nil {
			keyID := uint(id)
			params.KeyID = &keyID
		}
	}

	// Parse model
	if model := c.Query("model"); model != "" {
		params.Model = &model
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			params.Limit = limit
		}
	}

	// Parse offset
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			params.Offset = offset
		}
	}

	return params
}

// GetUsageSummary handles GET /usage - returns overall usage summary for the authenticated user
func (h *UsageHandler) GetUsageSummary(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	params := h.parseQueryParams(c)

	summary, err := h.usageService.GetUsageSummary(userID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetDailyUsage handles GET /usage/daily - returns usage data broken down by day
func (h *UsageHandler) GetDailyUsage(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	params := h.parseQueryParams(c)

	daily, err := h.usageService.GetDailyUsage(userID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, daily)
}

// GetUsageByKey handles GET /usage/by-key - returns usage data grouped by proxy API key
func (h *UsageHandler) GetUsageByKey(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	params := h.parseQueryParams(c)

	byKey, err := h.usageService.GetUsageByKey(userID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, byKey)
}

// GetUsageByProvider handles GET /usage/by-provider - returns usage data grouped by provider
func (h *UsageHandler) GetUsageByProvider(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	params := h.parseQueryParams(c)

	byProvider, err := h.usageService.GetUsageByProvider(userID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, byProvider)
}

// GetUsageByModel handles GET /usage/by-model - returns usage data grouped by model
func (h *UsageHandler) GetUsageByModel(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	params := h.parseQueryParams(c)

	byModel, err := h.usageService.GetUsageByModel(userID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, byModel)
}

// GetRecentUsage handles GET /usage/recent - returns recent usage records with pagination
func (h *UsageHandler) GetRecentUsage(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	params := h.parseQueryParams(c)

	// Get usage records
	records, err := h.usageService.GetRecentUsage(userID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get total count for pagination
	total, err := h.usageService.GetUsageCount(userID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"records": records,
		"total":   total,
		"limit":   params.Limit,
		"offset":  params.Offset,
	})
}
