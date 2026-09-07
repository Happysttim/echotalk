package handler

import (
	"echotalk/internal/auth"
	"echotalk/internal/errors"
	"echotalk/internal/mail"
	"echotalk/internal/redis"
	"echotalk/internal/response"
	"echotalk/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	VerifyToken       = "VERIFY_TOKEN"
	VerifyTokenMaxAge = 24 * 60 * 60
)

type VerifyHandler struct {
	userService *service.UserService
	redisUser   redis.RedisUser
}

func NewVerifyHandler(userService *service.UserService, redisUser redis.RedisUser) *VerifyHandler {
	return &VerifyHandler{
		userService: userService,
		redisUser:   redisUser,
	}
}

// @Router /verify/:email [get]
func (handler *VerifyHandler) EmailVerify(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()
	exists, err := handler.userService.GetUserByEmail(ctx, email)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	if exists != nil {
		response.Failed(c, http.StatusConflict, "email already exists")
		return
	}

	if err := handler.redisUser.IsVerify(ctx, email); err == nil {
		response.Failed(c, http.StatusConflict, "email already verified")
		return
	}

	uuid := uuid.New().String()
	if err := handler.redisUser.SetVerify(ctx, uuid, email); err != nil {
		response.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := mail.SendVerifyEmail(email, uuid); err != nil {
		response.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.OK(c, http.StatusOK)
}

// @Router /code/:code [get]
func (handler *VerifyHandler) CodeVerify(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()
	email, err := handler.redisUser.GetDelVerify(ctx, code)

	if err != nil || email == "" {
		response.Failed(c, http.StatusNotFound, "invalid verify code")
		return
	}

	token, err := auth.CreateVerifiedToken(email)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie(VerifyToken, token, VerifyTokenMaxAge, "/", "", false, true)
	response.OK(c, http.StatusOK)
}
