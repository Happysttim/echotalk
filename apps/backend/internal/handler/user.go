package handler

import (
	"echotalk/internal/auth"
	"echotalk/internal/model"
	"echotalk/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// @Router /users [post]
func (handler *UserHandler) Create(c *gin.Context) {
	userRequest := new(model.CreateUserRequest)
	if err := c.ShouldBindJSON(userRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	user, err := handler.userService.CreateUser(userRequest)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	tokenString, err := auth.CreateAccessToken(string(user.ID.Hex()))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.SetCookie("JWT_TOKEN", tokenString, int((30 * time.Minute).Seconds()), "/", "", true, true)
	c.JSON(http.StatusCreated, gin.H{"status": "OK"})
}
