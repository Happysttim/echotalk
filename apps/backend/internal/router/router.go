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
		uselessAuth.POST("/register", auth.VerifyRequired(), authHandler.LocalRegister)
		uselessAuth.POST("/change", auth.VerifyRequired(), authHandler.PasswordChange)

		needAuth := users.Group("", auth.AuthRequired())
		needAuth.GET("/logout", authHandler.Logout)
		needAuth.DELETE("", authHandler.DeleteAccount)
	}

	answers := engine.Group("/answers")
	{
		answers.POST("", auth.AuthNoRequired(), answerHandler.CreateAnswer)
		answers.DELETE("", auth.AuthNoRequired(), answerHandler.DeleteAnswer)

		needAuth := answers.Group("", auth.AuthRequired())
		needAuth.GET("/me", answerHandler.GetMyAnswers)
		needAuth.PATCH("", answerHandler.UpdateAnswer)
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
		needAuth.GET("/me", surveyHandler.GetMySurveyPage)

		surveys.GET("/:surveyId", surveyHandler.GetSurvey)
		surveys.GET("", surveyHandler.GetSurveyPage)
	}

	engine.GET("/verify", auth.AuthUnrequired(), verifyHandler.EmailVerify)
	engine.GET("/verify/me", auth.VerifyRequired(), verifyHandler.VerifyMe)
	engine.GET("/code", auth.AuthUnrequired(), verifyHandler.CodeVerify)
}
