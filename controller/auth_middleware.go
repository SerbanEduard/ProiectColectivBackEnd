package controller

import (
	"net/http"

	"github.com/SerbanEduard/ProiectColectivBackEnd/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthMiddleware verifies the Authorization header and stores claims in context
//
//	@Summary		JWT Authentication Middleware
//	@Description	Middleware to verify JWT token from Authorization header or query parameter
//	@Security		Bearer
//	@Param			Authorization	header		string				false	"Bearer token"
//	@Param			token			query		string				false	"Token for WebSocket connections"
//	@Success		200				{string}	string				"Token is valid"
//	@Failure		401				{object}	map[string]string	"Unauthorized"
//	@Router			/auth/middleware [post]
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

// RequireOwner ensures the authenticated subject matches the provided path parameter
//
//	@Summary		Owner Authorization Middleware
//	@Description	Middleware to ensure the authenticated user matches the resource owner
//	@Security		Bearer
//	@Param			id	path		string				true	"Resource ID that must match authenticated user ID"
//	@Success		200	{string}	string				"User is authorized"
//	@Failure		403	{object}	map[string]string	"Forbidden"
//	@Router			/auth/owner/{id} [post]
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
