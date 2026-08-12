package handler

import (
	"echotalk/internal/model"
	"echotalk/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	RefreshToken = "REFRESH_TOKEN"
)

const (
	RefreshTokenMaxAge = 2 * 60 * 60
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// @Router /auth/google [post]
func (handler *AuthHandler) GoogleLogin(c *gin.Context) {
	googleAuth := new(model.GoogleAuthRequest)
	if err := c.ShouldBindJSON(googleAuth); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	ctx := c.Request.Context()

	response, err := handler.authService.LoginWithGoogle(ctx, googleAuth)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.SetCookie(RefreshToken, response.RefreshToken, RefreshTokenMaxAge, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"status": "OK", "accessToken": response.AccessToken})
}

// @Router /auth/local [post]
func (handler *AuthHandler) LocalLogin(c *gin.Context) {
	localAuth := new(model.LocalAuthRequest)
	if err := c.ShouldBindJSON(localAuth); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	ctx := c.Request.Context()

	response, err := handler.authService.LoginWithLocal(ctx, localAuth)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.SetCookie(RefreshToken, response.RefreshToken, RefreshTokenMaxAge, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"status": "OK", "accessToken": response.AccessToken})
}

// @Router /auth/register [post]
func (handler *AuthHandler) LocalRegister(c *gin.Context) {
	localAuth := new(model.LocalAuthRequest)
	if err := c.ShouldBindJSON(localAuth); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	ctx := c.Request.Context()

	response, err := handler.authService.RegisterWithLocal(ctx, localAuth)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.SetCookie(RefreshToken, response.RefreshToken, RefreshTokenMaxAge, "/", "", true, true)
	c.JSON(http.StatusCreated, gin.H{"status": "OK", "accessToken": response.AccessToken})
}
