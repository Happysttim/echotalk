package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorized := c.GetHeader("Authorization")
		tokenString := extractAuthorized(authorized)

		if tokenString == "" {
			abortJSON(c)
			return
		}

		tokenClaims, err := ParseAccessToken(tokenString)

		if err != nil || tokenClaims == nil || tokenClaims.UserID == "" || tokenClaims.TokenType != "access" {
			abortJSON(c)
			return
		}

		c.Set("userID", tokenClaims.ID)
		c.Next()
	}
}

func extractAuthorized(authorized string) string {
	if !strings.HasPrefix(authorized, "Bearer") {
		return ""
	}

	return strings.TrimPrefix(authorized, "Bearer ")
}

func abortJSON(c *gin.Context) {
	c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
}
