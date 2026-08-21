package handler

import (
	"echotalk/internal/auth"
	"echotalk/internal/errors"
	"echotalk/internal/model"
	"echotalk/internal/response"
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
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()

	loginResponse, err := handler.authService.LoginWithGoogle(ctx, googleAuth)

	if err != nil {
		response.Failed(c, http.StatusInternalServerError, errors.ErrInternalServer.Error())
		return
	}

	c.SetCookie(RefreshToken, loginResponse.RefreshToken, RefreshTokenMaxAge, "/", "", true, true)
	response.OKWithData(c, http.StatusOK, map[string]string{
		"accessToken": loginResponse.AccessToken,
	})
}

// @Router /auth/local [post]
func (handler *AuthHandler) LocalLogin(c *gin.Context) {
	localAuth := new(model.LocalAuthRequest)
	if err := c.ShouldBindJSON(localAuth); err != nil {
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()

	loginResponse, err := handler.authService.LoginWithLocal(ctx, localAuth)

	if err != nil {
		response.Failed(c, http.StatusInternalServerError, errors.ErrInternalServer.Error())
		return
	}

	c.SetCookie(RefreshToken, loginResponse.RefreshToken, RefreshTokenMaxAge, "/", "", true, true)
	response.OKWithData(c, http.StatusOK, map[string]string{
		"accessToken": loginResponse.AccessToken,
	})
}

// @Router /auth/register [post]
func (handler *AuthHandler) LocalRegister(c *gin.Context) {
	localAuth := new(model.LocalAuthRequest)
	if err := c.ShouldBindJSON(localAuth); err != nil {
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()

	loginResponse, err := handler.authService.RegisterWithLocal(ctx, localAuth)

	if err != nil {
		response.Failed(c, http.StatusInternalServerError, errors.ErrInternalServer.Error())
		return
	}

	c.SetCookie(RefreshToken, loginResponse.RefreshToken, RefreshTokenMaxAge, "/", "", true, true)
	response.OKWithData(c, http.StatusOK, map[string]string{
		"accessToken": loginResponse.AccessToken,
	})
}

// @Router /auth/refresh [get]
func (handler *AuthHandler) Refresh(c *gin.Context) {
	cookie, err := c.Request.Cookie(RefreshToken)
	if err != nil {
		response.Failed(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error())
		return
	}

	ctx := c.Request.Context()
	refreshToken := cookie.Value

	if handler.authService.IsUselessRefreshToken(ctx, refreshToken) {
		response.FailedWithReason(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error(), "expired token error")
		return
	}

	claims, err := auth.ParseRefreshToken(refreshToken)

	if err != nil || claims == nil {
		response.FailedWithReason(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error(), "invalid refresh token")
		return
	}

	userID := claims.UserID

	exists, err := handler.userService.GetUserByID(ctx, userID)

	if err != nil {
		response.Failed(c, http.StatusInternalServerError, errors.ErrInternalServer.Error())
		return
	}

	if exists == nil {
		response.FailedWithReason(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error(), "invalid user data")
		return
	}

	jwtToken, err := auth.CreateToken(exists.ID.Hex())

	if err != nil {
		response.FailedWithReason(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error(), "internal server error")
		return
	}

	c.SetCookie(RefreshToken, jwtToken.RefreshToken, RefreshTokenMaxAge, "/", "", true, true)
	response.OKWithData(c, http.StatusOK, map[string]string{
		"accessToken": jwtToken.AccessToken,
	})
}

// @Router /auth/logout [get]
func (handler *AuthHandler) Logout(c *gin.Context) {
	cookie, err := c.Request.Cookie(RefreshToken)
	if err != nil {
		response.Failed(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error())
		return
	}

	ctx := c.Request.Context()
	refreshToken := cookie.Value

	if handler.authService.IsUselessRefreshToken(ctx, refreshToken) {
		response.FailedWithReason(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error(), "expired token error")
		return
	}

	if err := handler.authService.Logout(ctx, refreshToken); err != nil {
		response.Failed(c, http.StatusInternalServerError, errors.ErrInternalServer.Error())
		return
	}

	c.SetCookie(RefreshToken, "", -1, "/", "", true, true)
	response.OK(c, http.StatusOK)
}
