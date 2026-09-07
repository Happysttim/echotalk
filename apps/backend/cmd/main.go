package main

import (
	"echotalk/internal/auth"
	"echotalk/internal/config"
	"echotalk/internal/database"
	"echotalk/internal/handler"
	"echotalk/internal/redis"
	"echotalk/internal/repositories"
	"echotalk/internal/service"
	"log"
	"strconv"

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

	if err := database.InitDatabase(); err != nil {
		panic("database initialize is fail")
	}

	if config.Config == nil {
		panic("config is not loaded")
	}

	cfg := config.Config

	handlers := GetHandler()

	users := engine.Group("/auth")
	{
		users.GET("/refresh", handlers.AuthHandler.Refresh)
		uselessAuth := users.Group("", auth.AuthUnrequired())

		uselessAuth.GET("/google", handlers.AuthHandler.GoogleLogin)
		uselessAuth.POST("/local", handlers.AuthHandler.LocalLogin)
		uselessAuth.POST("/register", handlers.AuthHandler.LocalRegister).Use(auth.VerifyRequired())

		needAuth := users.Group("", auth.AuthRequired())
		needAuth.GET("/logout", handlers.AuthHandler.Logout)
		needAuth.DELETE("", handlers.AuthHandler.DeleteAccount)
	}

	answers := engine.Group("/answers")
	{
		needAuth := answers.Group("", auth.AuthRequired())
		needAuth.POST("", handlers.AnswerHandler.CreateAnswer)
		needAuth.PATCH("", handlers.AnswerHandler.UpdateAnswer)
		needAuth.DELETE("", handlers.AnswerHandler.DeleteAnswer)
		needAuth.POST("/rateup", handlers.AnswerHandler.RateUp)

		answers.GET("/:answerId", handlers.AnswerHandler.GetAnswer)
		answers.GET("", handlers.AnswerHandler.GetAnswerFeed)
	}

	surveys := engine.Group("/surveys")
	{
		needAuth := surveys.Group("", auth.AuthRequired())
		needAuth.POST("", handlers.SurveyHandler.CreateSurvey)
		needAuth.PATCH("", handlers.SurveyHandler.UpdateSurvey)
		needAuth.DELETE("", handlers.SurveyHandler.DeleteSurvey)

		surveys.GET("/:surveyId", handlers.SurveyHandler.GetSurvey)
		surveys.GET("", handlers.SurveyHandler.GetSurveyPage)
	}

	engine.GET("/verify/:email", handlers.VerifyHandler.EmailVerify).Use(auth.AuthUnrequired())
	engine.GET("/code/:code", handlers.VerifyHandler.CodeVerify).Use(auth.AuthUnrequired())

	if err := engine.Run(":" + strconv.Itoa(int(cfg.Port))); err != nil {
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
