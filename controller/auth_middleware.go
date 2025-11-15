package controller

import (
	"net/http"

	"github.com/SerbanEduard/ProiectColectivBackEnd/config"
	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware verifies the Authorization header and stores claims in context.
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// expected: "Bearer <token>" (HTTP) or "?token=<token>" (WebSocket)
		var tokenString string

		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString = authHeader[7:]
			} else {
				tokenString = authHeader
			}
		} else {
			// It might be a WebSocket request
			if c.GetHeader("Upgrade") == "websocket" {
				tokenString = c.Query("token")
			}
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		claims, err := config.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("userClaims", claims)

		c.Next()
	}
}
