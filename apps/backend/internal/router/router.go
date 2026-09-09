package router

import (
	"echotalk/internal/auth"
	"echotalk/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	engine *gin.Engine,
	authHandler *handler.AuthHandler,
	surveyHandler *handler.SurveyHandler,
	answerHandler *handler.AnswerHandler,
	verifyHandler *handler.VerifyHandler,
) {
	users := engine.Group("/auth")
	{
		users.GET("/refresh", authHandler.Refresh)
		uselessAuth := users.Group("", auth.AuthUnrequired())

		uselessAuth.GET("/google", authHandler.GoogleLogin)
		uselessAuth.POST("/local", authHandler.LocalLogin)
		uselessAuth.POST("/register", authHandler.LocalRegister).Use(auth.VerifyRequired())
		uselessAuth.POST("/change", authHandler.PasswordChange).Use(auth.VerifyRequired())

		needAuth := users.Group("", auth.AuthRequired())
		needAuth.GET("/logout", authHandler.Logout)
		needAuth.DELETE("", authHandler.DeleteAccount)
	}

	answers := engine.Group("/answers")
	{
		needAuth := answers.Group("", auth.AuthRequired())
		needAuth.POST("", answerHandler.CreateAnswer)
		needAuth.PATCH("", answerHandler.UpdateAnswer)
		needAuth.DELETE("", answerHandler.DeleteAnswer)
		needAuth.POST("/rateup", answerHandler.RateUp)

		answers.GET("/:answerId", answerHandler.GetAnswer)
		answers.GET("", answerHandler.GetAnswerFeed)
	}

	surveys := engine.Group("/surveys")
	{
		needAuth := surveys.Group("", auth.AuthRequired())
		needAuth.POST("", surveyHandler.CreateSurvey)
		needAuth.PATCH("", surveyHandler.UpdateSurvey)
		needAuth.DELETE("", surveyHandler.DeleteSurvey)

		surveys.GET("/:surveyId", surveyHandler.GetSurvey)
		surveys.GET("", surveyHandler.GetSurveyPage)
	}

	engine.GET("/verify", verifyHandler.EmailVerify).Use(auth.AuthUnrequired())
	engine.GET("/code", verifyHandler.CodeVerify).Use(auth.AuthUnrequired())
}
