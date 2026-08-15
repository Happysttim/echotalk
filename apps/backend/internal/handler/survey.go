package handler

import (
	"echotalk/internal/service"

	"github.com/gin-gonic/gin"
)

type SurveyHandler struct {
	surveyService *service.SurveyService
}

func NewSurveyHandler(surveyService *service.SurveyService) *SurveyHandler {
	return &SurveyHandler{
		surveyService: surveyService,
	}
}

// @Router /surveys/:id [get]
func (handler *SurveyHandler) GetSurvey(c *gin.Context) {

}
