package handler

import (
	"echotalk/internal/auth"
	"echotalk/internal/errors"
	"echotalk/internal/model"
	"echotalk/internal/response"
	"echotalk/internal/service"
	"net/http"
	"strings"

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

// @Router /auth/google [get]
func (handler *AuthHandler) GoogleLogin(c *gin.Context) {
	googleAuth := new(model.GoogleAuthRequest)
	if err := c.ShouldBindQuery(googleAuth); err != nil {
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()

	loginResponse, err := handler.authService.LoginWithGoogle(ctx, googleAuth)

	if err != nil {
		response.Failed(c, http.StatusInternalServerError, errors.ErrInternalServer.Error())
		return
	}

	c.SetCookie(RefreshToken, loginResponse.RefreshToken, RefreshTokenMaxAge, "/", "", false, true)
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
		response.Failed(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error())
		return
	}

	c.SetCookie(RefreshToken, loginResponse.RefreshToken, RefreshTokenMaxAge, "/", "", false, true)
	response.OKWithData(c, http.StatusOK, map[string]string{
		"accessToken": loginResponse.AccessToken,
	})
}

// @Router /auth/register [post]
func (handler *AuthHandler) LocalRegister(c *gin.Context) {
	localAuth := new(model.LocalAuthRegisterRequest)
	email := c.GetString("email")

	if err := c.ShouldBindJSON(localAuth); err != nil || email == "" {
		response.Failed(c, http.StatusBadRequest, errors.ErrInvalidInput.Error())
		return
	}

	if strings.Compare(localAuth.Email, email) != 0 {
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()

	loginResponse, err := handler.authService.RegisterWithLocal(ctx, localAuth)

	if err != nil {
		response.Failed(c, http.StatusInternalServerError, errors.ErrInternalServer.Error())
		return
	}

	c.SetCookie(RefreshToken, loginResponse.RefreshToken, RefreshTokenMaxAge, "/", "", false, true)
	c.SetCookie(VerifyToken, "", -1, "/", "", false, true)
	response.OKWithData(c, http.StatusOK, map[string]string{
		"accessToken": loginResponse.AccessToken,
	})
}

// @Router /auth/change [post]
func (handler *AuthHandler) PasswordChange(c *gin.Context) {
	localAuth := new(model.LocalAuthRequest)
	email := c.GetString("email")

	if err := c.ShouldBindJSON(localAuth); err != nil || email == "" {
		response.Failed(c, http.StatusBadRequest, errors.ErrInvalidInput.Error())
		return
	}

	if strings.Compare(localAuth.Email, email) != 0 {
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()
	if err := handler.authService.ChangeWithLocal(ctx, localAuth); err != nil {
		response.Failed(c, http.StatusBadRequest, errors.ErrInvalidInput.Error())
		return
	}
	c.SetCookie(VerifyToken, "", -1, "/", "", false, true)
	response.OK(c, http.StatusAccepted)
}

// @Router /auth/refresh [get]
func (handler *AuthHandler) Refresh(c *gin.Context) {
	cookie, err := c.Request.Cookie(RefreshToken)
	if err != nil {
		response.Failed(c, http.StatusUnauthorized, err.Error())
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
		response.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	if exists == nil {
		response.FailedWithReason(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error(), "invalid user data")
		return
	}

	jwtToken, err := auth.CreateToken(exists.ID.Hex())

	if err != nil {
		response.FailedWithReason(c, http.StatusUnauthorized, err.Error(), "internal server error")
		return
	}

	if err := handler.authService.RevokeRefreshToken(ctx, refreshToken); err != nil {
		response.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie(RefreshToken, jwtToken.RefreshToken, RefreshTokenMaxAge, "/", "", false, true)
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
		response.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie(RefreshToken, "", -1, "/", "", false, true)
	response.OK(c, http.StatusOK)
}

// @Router /auth [delete]
func (handler *AuthHandler) DeleteAccount(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		response.Failed(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error())
		return
	}

	ctx := c.Request.Context()
	user, err := handler.userService.GetUserByID(ctx, userID)

	if err != nil || user == nil {
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	if err := handler.userService.DeleteUser(ctx, user.ID.Hex()); err != nil {
		response.Failed(c, http.StatusInternalServerError, errors.ErrInternalServer.Error())
		return
	}

	c.SetCookie(RefreshToken, "", -1, "/", "", false, true)
	response.OK(c, http.StatusOK)
}
