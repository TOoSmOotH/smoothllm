package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smoothweb/backend/internal/auth"
	"github.com/smoothweb/backend/internal/custom/services"
)

// ProviderHandler handles provider API endpoints
type ProviderHandler struct {
	providerService *services.ProviderService
}

// NewProviderHandler creates a new ProviderHandler instance
func NewProviderHandler(providerService *services.ProviderService) *ProviderHandler {
	return &ProviderHandler{
		providerService: providerService,
	}
}

// ListProviders handles GET /providers - lists all providers for the authenticated user
func (h *ProviderHandler) ListProviders(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	providers, err := h.providerService.ListProviders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, providers)
}

// GetProvider handles GET /providers/:id - retrieves a specific provider
func (h *ProviderHandler) GetProvider(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	providerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider id"})
		return
	}

	provider, err := h.providerService.GetProvider(userID, uint(providerID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, provider)
}

// CreateProvider handles POST /providers - creates a new provider
func (h *ProviderHandler) CreateProvider(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req services.CreateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	provider, err := h.providerService.CreateProvider(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, provider)
}

// UpdateProvider handles PUT /providers/:id - updates an existing provider
func (h *ProviderHandler) UpdateProvider(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	providerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider id"})
		return
	}

	var req services.UpdateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	provider, err := h.providerService.UpdateProvider(userID, uint(providerID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, provider)
}

// DeleteProvider handles DELETE /providers/:id - deletes a provider
func (h *ProviderHandler) DeleteProvider(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	providerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider id"})
		return
	}

	if err := h.providerService.DeleteProvider(userID, uint(providerID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// TestConnection handles POST /providers/:id/test - tests provider connection
func (h *ProviderHandler) TestConnection(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	providerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider id"})
		return
	}

	if err := h.providerService.TestConnection(userID, uint(providerID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "connection successful"})
}

// TestConnectionWithCredentials handles POST /providers/test - tests connection with provided credentials (before saving)
func (h *ProviderHandler) TestConnectionWithCredentials(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req services.CreateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.providerService.TestConnectionWithRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "connection successful"})
}

// FetchModels handles GET /providers/:id/available-models - fetches models from provider
func (h *ProviderHandler) FetchModels(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	providerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider id"})
		return
	}

	models, err := h.providerService.FetchAvailableModels(userID, uint(providerID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"models": models})
}

// FetchModelsWithCredentials handles POST /providers/fetch-models - fetches models with provided credentials (before saving)
func (h *ProviderHandler) FetchModelsWithCredentials(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req services.CreateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models, err := h.providerService.FetchAvailableModelsWithRequest(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"models": models})
}
