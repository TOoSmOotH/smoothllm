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

	// Helper to render a page that closes the popup
	renderClosePage := func(success bool, message string) {
		status := "error"
		if success {
			status = "success"
		}
		html := `<!DOCTYPE html>
<html>
<head>
    <title>OAuth ` + status + `</title>
    <style>
        body { font-family: system-ui, sans-serif; display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0; background: #1a1a2e; color: #fff; }
        .container { text-align: center; padding: 2rem; }
        .icon { font-size: 48px; margin-bottom: 1rem; }
        .message { font-size: 18px; margin-bottom: 1rem; }
        .close-note { font-size: 14px; color: #888; }
    </style>
</head>
<body>
    <div class="container">
        <div class="icon">` + func() string { if success { return "✓" } else { return "✗" } }() + `</div>
        <div class="message">` + message + `</div>
        <div class="close-note">This window will close automatically...</div>
    </div>
    <script>
        setTimeout(function() { window.close(); }, 2000);
    </script>
</body>
</html>`
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, html)
	}

	// Handle OAuth errors
	if errorParam != "" {
		message := "OAuth error: " + errorParam
		if errorDescription != "" {
			message = errorDescription
		}
		renderClosePage(false, message)
		return
	}

	if code == "" || state == "" {
		renderClosePage(false, "Missing required parameters")
		return
	}

	// Build the redirect URI (must match what was sent in authorize)
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	redirectURI := scheme + "://" + c.Request.Host + "/api/v1/oauth/anthropic/callback"

	// Exchange the code for tokens
	_, err := h.oauthService.ExchangeCode(code, state, redirectURI)
	if err != nil {
		renderClosePage(false, "Failed to connect: "+err.Error())
		return
	}

	renderClosePage(true, "Successfully connected to Claude Max!")
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
