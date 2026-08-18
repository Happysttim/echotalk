package handler

import (
	"echotalk/internal/model"
	"echotalk/internal/service"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AnswerCursor struct {
	SurveyID string `json:"survey_id"`
	AnswerID string `json:"answer_id"`
}

type AnswerHandler struct {
	answerService *service.AnswerService
	surveyService *service.SurveyService
}

func NewAnswerHandler(answerService *service.AnswerService, surveyService *service.SurveyService) *AnswerHandler {
	return &AnswerHandler{
		answerService: answerService,
		surveyService: surveyService,
	}
}

// @Router /answers [post]
func (handler *AnswerHandler) CreateAnswer(c *gin.Context) {
	payload := new(model.CreateAnswerRequest)
	userID := c.GetString("userID")

	if userID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := c.ShouldBindJSON(payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	ctx := c.Request.Context()

	answer, err := handler.answerService.CreateAnswer(ctx, payload, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created", "answer": answer})
}

// @Router /answers/:answerId [get]
func (handler *AnswerHandler) GetAnswer(c *gin.Context) {
	answerId := c.Param("answerId")
	if answerId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	ctx := c.Request.Context()
	answer, err := handler.answerService.GetAnswerByID(ctx, answerId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "answer": answer})
}

// @Router /answers [patch]
func (handler *AnswerHandler) UpdateAnswer(c *gin.Context) {
	payload := new(model.UpdateAnswerRequest)
	userID := c.GetString("userID")

	if userID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	ctx := c.Request.Context()
	answer, err := handler.answerService.GetAnswerByID(ctx, payload.AnswerID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if strings.Compare(answer.AuthorID.Hex(), userID) != 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := handler.answerService.UpdateAnswer(ctx, payload); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// @Router /answers?cursor={string} [get]
func (handler *AnswerHandler) GetAnswerFeed(c *gin.Context) {
	cursorBase64 := c.Query("cursor")

	if cursorBase64 == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	decodeBytes, err := base64.RawURLEncoding.DecodeString(cursorBase64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	cursor := new(AnswerCursor)
	if err := json.Unmarshal(decodeBytes, cursor); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	ctx := c.Request.Context()
	survey, err := handler.surveyService.GetSurveyByID(ctx, cursor.SurveyID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if survey == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	answer, err := handler.answerService.GetAnswerByID(ctx, cursor.AnswerID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if answer == nil || answer.SurveyID != survey.ID {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	answers, err := handler.answerService.GetFilteredAnswers(
		ctx,
		bson.M{
			"_id": bson.M{
				"$lt": answer.ID,
			},
			"survey_id": survey.ID,
		},
		options.Find().SetSort(
			bson.D{
				{Key: "_id", Value: -1},
				{Key: "rate_up", Value: -1},
			},
		).SetLimit(LimitSmall+1),
	)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var nextSkip string
	var nextID string
	hasNext := len(answers) > LimitSmall
	nextCursor := AnswerCursor{
		SurveyID: survey.ID.Hex(),
	}
	if hasNext && len(answers) > 0 {
		last := answers[len(answers)-1]
		nextID = last.ID.Hex()
	}

	if nextID != "" {
		nextCursor.AnswerID = nextID
	}

	jsonBytes, err := json.Marshal(nextCursor)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	nextSkip = base64.RawURLEncoding.EncodeToString(jsonBytes)

	c.JSON(http.StatusOK, gin.H{"status": "ok", "skip": nextSkip, "hasNext": hasNext, "answers": answers[:LimitSmall]})
}
