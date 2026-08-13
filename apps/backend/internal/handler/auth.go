package handler

import (
	"echotalk/internal/auth"
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
	userService *service.UserService
}

func NewAuthHandler(authService *service.AuthService, userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
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

// @Router /auth/refresh [get]
func (handler *AuthHandler) Refresh(c *gin.Context) {
	cookie, err := c.Request.Cookie(RefreshToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ctx := c.Request.Context()
	refreshToken := cookie.Value

	if handler.authService.IsUselessRefreshToken(ctx, refreshToken) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "expired token error"})
		return
	}

	claims, err := auth.ParseRefreshToken(refreshToken)

	if err != nil || claims == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	userID := claims.UserID

	exists, err := handler.userService.GetUserByID(ctx, userID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if exists == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user data"})
		return
	}

	jwtToken, err := auth.CreateToken(exists.ID.Hex())

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.SetCookie(RefreshToken, jwtToken.RefreshToken, RefreshTokenMaxAge, "/", "", true, true)
	c.JSON(http.StatusCreated, gin.H{"status": "OK", "accessToken": jwtToken.AccessToken})
}

// @Router /auth/logout [get]
func (handler *AuthHandler) Logout(c *gin.Context) {
	cookie, err := c.Request.Cookie(RefreshToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ctx := c.Request.Context()
	refreshToken := cookie.Value

	if handler.authService.IsUselessRefreshToken(ctx, refreshToken) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "expired token error"})
		return
	}

	if err := handler.authService.Logout(ctx, refreshToken); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.SetCookie(RefreshToken, "", -1, "/", "", true, true)
	c.JSON(http.StatusCreated, gin.H{"status": "OK"})
}
