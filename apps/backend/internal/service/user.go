import (
	"apps/backend/internal/model"
	"apps/backend/internal/repositories"
	"context"
	"time"
)

type UserService struct {
	userRepo *repositories.MongoUserRepository
}

func NewUserSerevice(userRepo *repositories.MongoUserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}