package handler

import (
	"echotalk/internal/auth"
	"echotalk/internal/errors"
	"echotalk/internal/mail"
	"echotalk/internal/model"
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

// @Router /verify?email={string}&verifyType={string} [get]
func (handler *VerifyHandler) EmailVerify(c *gin.Context) {
	email := c.Query("email")
	verifyTypeQuery := c.Query("verifyType")

	if email == "" || !redis.IsVerifyType(verifyTypeQuery) {
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()
	exists, err := handler.userService.GetUserByEmail(ctx, email)

	if err != nil {
		response.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	var verifyType redis.VerifyType
	var location string

	if exists != nil {
		if !redis.IsPasswordVerify(verifyTypeQuery) || exists.AuthProvider != model.AuthProviderLocal {
			response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
			return
		}
		verifyType = redis.PasswordVerify
		location = "change"
	} else {
		if !redis.IsRegisterVerify(verifyTypeQuery) {
			response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
			return
		}
		verifyType = redis.RegisterVerify
		location = "register"
	}

	if err := handler.redisUser.IsVerifyEmail(ctx, email, verifyType); err == nil {
		response.Failed(c, http.StatusConflict, "email already verified")
		return
	}

	uuid := uuid.New().String()
	if err := mail.SendVerifyEmail(email, uuid, location); err != nil {
		response.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := handler.redisUser.SetVerify(ctx, uuid, email, verifyType); err != nil {
		response.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.OK(c, http.StatusOK)
}

// @Router /code?code={string}&verifyType={string} [get]
func (handler *VerifyHandler) CodeVerify(c *gin.Context) {
	code := c.Query("code")
	verifyTypeQuery := c.Query("verifyType")

	if code == "" || !redis.IsVerifyType(verifyTypeQuery) {
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	var verifyType redis.VerifyType = redis.UnknownVerify

	if redis.IsPasswordVerify(verifyTypeQuery) {
		verifyType = redis.PasswordVerify
	} else if redis.IsRegisterVerify(verifyTypeQuery) {
		verifyType = redis.RegisterVerify
	} else {
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()
	email, err := handler.redisUser.GetDelVerify(ctx, code, verifyType)

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
