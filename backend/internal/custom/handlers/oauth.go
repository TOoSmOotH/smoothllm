package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smoothweb/backend/internal/auth"
	"github.com/smoothweb/backend/internal/custom/services"
)

// OAuthHandler handles OAuth-related endpoints
type OAuthHandler struct {
	oauthService *services.OAuthService
	frontendURL  string
}

// NewOAuthHandler creates a new OAuthHandler instance
func NewOAuthHandler(oauthService *services.OAuthService, frontendURL string) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
		frontendURL:  frontendURL,
	}
}

// AuthorizeRequest represents the request to start OAuth flow
type AuthorizeRequest struct {
	ProviderID uint `json:"provider_id" binding:"required"`
}

// Authorize handles GET /oauth/anthropic/authorize - starts the OAuth flow
func (h *OAuthHandler) Authorize(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	providerIDStr := c.Query("provider_id")
	if providerIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider_id is required"})
		return
	}

	providerID, err := strconv.ParseUint(providerIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider_id"})
		return
	}

	// Build the redirect URI - this should be the backend callback endpoint
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	redirectURI := scheme + "://" + c.Request.Host + "/api/v1/oauth/anthropic/callback"

	authURL, err := h.oauthService.GenerateAuthURL(userID, uint(providerID), redirectURI)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return the authorization URL for the frontend to redirect to
	c.JSON(http.StatusOK, gin.H{"authorization_url": authURL})
}

// Callback handles GET /oauth/anthropic/callback - handles the OAuth callback
func (h *OAuthHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")
	errorDescription := c.Query("error_description")

	// Handle OAuth errors
	if errorParam != "" {
		// Redirect to frontend with error
		redirectURL := h.frontendURL + "/providers?oauth_error=" + errorParam
		if errorDescription != "" {
			redirectURL += "&error_description=" + errorDescription
		}
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	if code == "" || state == "" {
		c.Redirect(http.StatusFound, h.frontendURL+"/providers?oauth_error=missing_parameters")
		return
	}

	// Build the redirect URI (must match what was sent in authorize)
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	redirectURI := scheme + "://" + c.Request.Host + "/api/v1/oauth/anthropic/callback"

	// Exchange the code for tokens
	provider, err := h.oauthService.ExchangeCode(code, state, redirectURI)
	if err != nil {
		c.Redirect(http.StatusFound, h.frontendURL+"/providers?oauth_error=exchange_failed&error_description="+err.Error())
		return
	}

	// Redirect to frontend with success
	c.Redirect(http.StatusFound, h.frontendURL+"/providers?oauth_success=true&provider_id="+strconv.FormatUint(uint64(provider.ID), 10))
}

// Disconnect handles POST /oauth/anthropic/disconnect/:id - disconnects OAuth from a provider
func (h *OAuthHandler) Disconnect(c *gin.Context) {
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

	if err := h.oauthService.DisconnectOAuth(userID, uint(providerID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OAuth disconnected successfully"})
}

// TestOAuthConnection handles POST /oauth/anthropic/test/:id - tests OAuth connection
func (h *OAuthHandler) TestOAuthConnection(c *gin.Context) {
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

	// Get the provider
	provider, err := h.oauthService.GetProviderByID(userID, uint(providerID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := h.oauthService.TestOAuthConnection(provider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OAuth connection successful"})
}
