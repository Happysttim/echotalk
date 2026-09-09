package main

import (
	"context"
	"echotalk/internal/config"
	"echotalk/internal/database"
	"echotalk/internal/handler"
	"echotalk/internal/redis"
	"echotalk/internal/repositories"
	"echotalk/internal/router"
	"echotalk/internal/server"
	"echotalk/internal/service"
	"echotalk/internal/timer"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	AuthHandler   *handler.AuthHandler
	SurveyHandler *handler.SurveyHandler
	AnswerHandler *handler.AnswerHandler
	VerifyHandler *handler.VerifyHandler
}

func main() {
	engine := gin.New()
	ticker := timer.NewTimeTicker(time.Minute)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	defer ticker.Stop()

	if err := database.InitDatabase(); err != nil {
		panic("database initialize is fail")
	}

	if config.Config == nil {
		panic("config is not loaded")
	}

	cfg := config.Config

	handlers := GetHandler()
	router.RegisterRoutes(engine, handlers.AuthHandler, handlers.SurveyHandler, handlers.AnswerHandler, handlers.VerifyHandler)

	ticker.Start(handlers.SurveyHandler.Ticker)

	if err := server.Run(ctx, cfg.Port, engine); err != nil {
		log.Fatalln(err)
	}
}

func GetHandler() *Handlers {
	userRepository := repositories.NewMongoUserRepository()
	sessionRepository := repositories.NewMongoSessionRepository()
	surveyRepository := repositories.NewMongoSurveyRepository()
	rateUpRepository := repositories.NewRateUpRepository()
	answerRepository := repositories.NewMongoAnswerRepository()

	userService := service.NewUserService(userRepository)
	authService := service.NewAuthService(sessionRepository, userService)
	answerService := service.NewAnswerService(answerRepository, rateUpRepository, userService)
	surveyService := service.NewSurveyService(surveyRepository, userService)

	redisUser := redis.RedisUser{}

	return &Handlers{
		AuthHandler:   handler.NewAuthHandler(authService, userService),
		SurveyHandler: handler.NewSurveyHandler(surveyService, answerService, userService),
		AnswerHandler: handler.NewAnswerHandler(answerService, surveyService),
		VerifyHandler: handler.NewVerifyHandler(userService, redisUser),
	}
}
