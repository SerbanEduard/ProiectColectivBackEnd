package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserIDFromContext(c *gin.Context) (string, error) {
	if value, exists := c.Get("userClaims"); exists {
		claims, ok := value.(jwt.MapClaims)
		if !ok {
			return "", fmt.Errorf("userClaims is not of type jwt.MapClaims")
		}

		userID, err := claims.GetSubject()
		if err != nil {
			return "", err
		}

		return userID, nil
	}
	return "", fmt.Errorf("userClaims not found in context")
}
