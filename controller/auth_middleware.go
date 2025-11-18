package controller

import (
	"net/http"

	"github.com/SerbanEduard/ProiectColectivBackEnd/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

// RequireOwner ensures the authenticated subject (sub claim) matches the provided
// path parameter (for example :id). Use as a route-level middleware after JWTAuthMiddleware.
func RequireOwner(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsI, ok := c.Get("userClaims")
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		claims, ok := claimsI.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		subI, ok := claims["sub"]
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		sub, ok := subI.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		if sub != c.Param(paramName) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}
