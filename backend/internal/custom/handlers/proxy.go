package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smoothweb/backend/internal/custom/services"
)

// ProxyHandler handles LLM proxy API endpoints
type ProxyHandler struct {
	proxyService *services.ProxyService
}

// NewProxyHandler creates a new ProxyHandler instance
func NewProxyHandler(proxyService *services.ProxyService) *ProxyHandler {
	return &ProxyHandler{
		proxyService: proxyService,
	}
}

// ChatCompletions handles POST /v1/chat/completions
// This is the main OpenAI-compatible endpoint that proxies requests to the configured provider
func (h *ProxyHandler) ChatCompletions(c *gin.Context) {
	// Extract the proxy API key from the Authorization header
	apiKey, err := h.proxyService.GetProxyKeyFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"message": err.Error(),
				"type":    "authentication_error",
				"code":    "invalid_api_key",
			},
		})
		return
	}

	// Validate the key
	proxyKey, err := h.proxyService.ValidateKey(apiKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"message": err.Error(),
				"type":    "authentication_error",
				"code":    "invalid_api_key",
			},
		})
		return
	}

	// Proxy the request to the provider
	result, err := h.proxyService.ProxyRequest(c, proxyKey)
	if err != nil {
		// If ProxyRequest already wrote to the response (via c.Data), don't write again
		if c.Writer.Written() {
			return
		}

		// Handle provider errors appropriately
		statusCode := http.StatusBadGateway
		if result != nil && result.StatusCode > 0 {
			statusCode, _ = h.proxyService.HandleProviderError(result.StatusCode, result.ErrorMessage)
		}

		c.JSON(statusCode, gin.H{
			"error": gin.H{
				"message": err.Error(),
				"type":    "api_error",
				"code":    "proxy_error",
			},
		})
		return
	}

	// Response was already written by ProxyRequest via c.Data()
	// No need to write anything else here
}

// ListModels handles GET /v1/models
// Returns a list of available models based on the proxy key's provider
func (h *ProxyHandler) ListModels(c *gin.Context) {
	// Extract the proxy API key from the Authorization header
	apiKey, err := h.proxyService.GetProxyKeyFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"message": err.Error(),
				"type":    "authentication_error",
				"code":    "invalid_api_key",
			},
		})
		return
	}

	// Validate the key
	proxyKey, err := h.proxyService.ValidateKey(apiKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"message": err.Error(),
				"type":    "authentication_error",
				"code":    "invalid_api_key",
			},
		})
		return
	}

	// Get the list of available models for this key
	models, err := h.proxyService.ListModelsForKey(proxyKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": err.Error(),
				"type":    "api_error",
				"code":    "server_error",
			},
		})
		return
	}

	c.JSON(http.StatusOK, models)
}

// Messages handles POST /v1/messages
// This is the Anthropic-compatible endpoint that proxies requests to Claude Max providers
func (h *ProxyHandler) Messages(c *gin.Context) {
	// Extract the proxy API key from the Authorization header
	apiKey, err := h.proxyService.GetProxyKeyFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"type":    "authentication_error",
				"message": err.Error(),
			},
		})
		return
	}

	// Validate the key
	proxyKey, err := h.proxyService.ValidateKey(apiKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"type":    "authentication_error",
				"message": err.Error(),
			},
		})
		return
	}

	// Proxy the request to Anthropic
	result, err := h.proxyService.ProxyAnthropicPassthrough(c, proxyKey)
	if err != nil {
		// If ProxyAnthropicPassthrough already wrote to the response, don't write again
		if c.Writer.Written() {
			return
		}

		// Handle provider errors appropriately
		statusCode := http.StatusBadGateway
		if result != nil && result.StatusCode > 0 {
			statusCode = result.StatusCode
		}

		c.JSON(statusCode, gin.H{
			"error": gin.H{
				"type":    "api_error",
				"message": err.Error(),
			},
		})
		return
	}

	// Response was already written by ProxyAnthropicPassthrough via c.Data()
}
