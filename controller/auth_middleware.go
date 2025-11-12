package controller

import (
	"net/http"

	"github.com/SerbanEduard/ProiectColectivBackEnd/config"
	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware verifies the Authorization header and stores claims in context.
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		// expected: "Bearer <token>"
		var tokenString string
		if len(auth) > 7 && auth[:7] == "Bearer " {
			tokenString = auth[7:]
		} else {
			tokenString = auth
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
