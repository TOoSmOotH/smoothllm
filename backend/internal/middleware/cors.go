package middleware

import (
	"github.com/gin-gonic/gin"
)

func CORS(allowedOrigins []string) gin.HandlerFunc {
	// Check if wildcard is enabled
	allowAll := false
	for _, o := range allowedOrigins {
		if o == "*" {
			allowAll = true
			break
		}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if allowAll && origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
