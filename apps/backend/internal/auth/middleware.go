package auth

import (
	"echotalk/internal/errors"
	"echotalk/internal/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorized := c.GetHeader("Authorization")
		tokenString := extractAuthorized(authorized)

		if tokenString == "" {
			abortUnauthorized(c)
			return
		}

		tokenClaims, err := ParseAccessToken(tokenString)

		if err != nil || tokenClaims == nil || tokenClaims.UserID == "" || tokenClaims.TokenType != "access" {
			abortUnauthorized(c)
			return
		}

		c.Set("userID", tokenClaims.UserID)
		c.Next()
	}
}

func AuthNoRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorized := c.GetHeader("Authorization")
		if authorized == "" {
			c.Next()
			return
		}

		tokenString := extractAuthorized(authorized)

		if tokenString == "" {
			abortUnauthorized(c)
			return
		}

		tokenClaims, err := ParseAccessToken(tokenString)

		if err != nil || tokenClaims == nil || tokenClaims.UserID == "" || tokenClaims.TokenType != "access" {
			abortUnauthorized(c)
			return
		}

		c.Set("userID", tokenClaims.UserID)
		c.Next()
	}
}

func AuthUnrequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorized := c.GetHeader("Authorization")
		tokenString := extractAuthorized(authorized)

		if tokenString != "" {
			abortBadRequest(c)
			return
		}

		c.Next()
	}
}

func VerifyRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("VERIFY_TOKEN")
		if err != nil {
			abortUnauthorized(c)
			return
		}

		claims, err := ParseVerifiedToken(cookie)
		if err != nil || claims == nil {
			abortUnauthorized(c)
			return
		}

		c.Set("email", claims.Email)
		c.Next()
	}
}

func extractAuthorized(authorized string) string {
	if !strings.HasPrefix(authorized, "Bearer") {
		return ""
	}

	return strings.TrimPrefix(authorized, "Bearer ")
}

func abortUnauthorized(c *gin.Context) {
	response.Failed(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error())
}

func abortBadRequest(c *gin.Context) {
	response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
}
