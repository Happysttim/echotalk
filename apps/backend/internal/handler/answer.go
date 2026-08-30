package handler

import (
	"echotalk/internal/errors"
	"echotalk/internal/model"
	"echotalk/internal/response"
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
		response.Failed(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error())
		return
	}
	if err := c.ShouldBindJSON(payload); err != nil {
		response.Failed(c, http.StatusUnauthorized, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()

	answer, err := handler.answerService.CreateAnswer(ctx, payload, userID)
	if err != nil {
		response.Failed(c, http.StatusUnauthorized, errors.ErrInternalServer.Error())
		return
	}

	response.OKWithData(c, http.StatusCreated, answer)
}

// @Router /answers [delete]
func (handler *AnswerHandler) DeleteAnswer(c *gin.Context) {
	payload := new(model.DeleteAnswerRequest)

	if err := c.ShouldBindJSON(payload); err != nil {
		response.Failed(c, http.StatusUnauthorized, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()
	answer, err := handler.answerService.GetAnswerByID(ctx, payload.AnswerID)

	if err != nil {
		response.Failed(c, http.StatusUnauthorized, errors.ErrInternalServer.Error())
		return
	}

	if err := handler.answerService.DeleteAnswer(ctx, answer.ID.Hex()); err != nil {
		response.Failed(c, http.StatusUnauthorized, errors.ErrInternalServer.Error())
		return
	}

	response.OK(c, http.StatusOK)
}

// @Router /answers/:answerId [get]
func (handler *AnswerHandler) GetAnswer(c *gin.Context) {
	answerId := c.Param("answerId")
	if answerId == "" {
		response.Failed(c, http.StatusUnauthorized, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()
	answer, err := handler.answerService.GetAnswerByID(ctx, answerId)

	if err != nil {
		response.Failed(c, http.StatusUnauthorized, errors.ErrInternalServer.Error())
		return
	}

	response.OKWithData(c, http.StatusOK, answer)
}

// @Router /answers [patch]
func (handler *AnswerHandler) UpdateAnswer(c *gin.Context) {
	payload := new(model.UpdateAnswerRequest)
	userID := c.GetString("userID")

	if userID == "" {
		response.Failed(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error())
		return
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Failed(c, http.StatusUnauthorized, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()
	answer, err := handler.answerService.GetAnswerByID(ctx, payload.AnswerID)

	if err != nil {
		response.Failed(c, http.StatusUnauthorized, errors.ErrInternalServer.Error())
		return
	}

	if strings.Compare(answer.AuthorID.Hex(), userID) != 0 {
		response.Failed(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error())
		return
	}

	if err := handler.answerService.UpdateAnswer(ctx, payload); err != nil {
		response.Failed(c, http.StatusUnauthorized, errors.ErrInternalServer.Error())
		return
	}

	response.OK(c, http.StatusOK)
}

// @Router /answers?cursor={string} [get]
func (handler *AnswerHandler) GetAnswerFeed(c *gin.Context) {
	cursorBase64 := c.Query("cursor")

	if cursorBase64 == "" {
		response.Failed(c, http.StatusUnauthorized, errors.ErrBadRequest.Error())
		return
	}

	decodeBytes, err := base64.RawURLEncoding.DecodeString(cursorBase64)
	if err != nil {
		response.Failed(c, http.StatusUnauthorized, errors.ErrBadRequest.Error())
		return
	}

	cursor := new(AnswerCursor)
	if err := json.Unmarshal(decodeBytes, cursor); err != nil {
		response.Failed(c, http.StatusUnauthorized, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()
	survey, err := handler.surveyService.GetSurveyByID(ctx, cursor.SurveyID)

	if err != nil {
		response.Failed(c, http.StatusUnauthorized, errors.ErrInternalServer.Error())
		return
	}

	if survey == nil {
		response.Failed(c, http.StatusNotFound, errors.ErrNotFound.Error())
		return
	}

	answer, err := handler.answerService.GetAnswerByID(ctx, cursor.AnswerID)

	if err != nil {
		response.Failed(c, http.StatusUnauthorized, errors.ErrInternalServer.Error())
		return
	}

	if answer == nil || answer.SurveyID != survey.ID {
		response.Failed(c, http.StatusNotFound, errors.ErrNotFound.Error())
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
		response.Failed(c, http.StatusUnauthorized, errors.ErrInternalServer.Error())
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
		response.Failed(c, http.StatusUnauthorized, errors.ErrInternalServer.Error())
		return
	}

	nextSkip = base64.RawURLEncoding.EncodeToString(jsonBytes)

	response.OKWithData(c, http.StatusOK, map[string]any{"skip": nextSkip, "hasNext": hasNext, "answers": answers[:LimitSmall]})
}

// @Router /answers/rateup [post]
func (handler *AnswerHandler) RateUp(c *gin.Context) {
	payload := new(model.RateUpRequest)
	userIdString := c.GetString("userId")

	if err := c.ShouldBindJSON(payload); err != nil {
		response.Failed(c, http.StatusNotFound, errors.ErrNotFound.Error())
		return
	}

	if userIdString == "" {
		response.Failed(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error())
		return
	}

	ctx := c.Request.Context()
	userId, err := bson.ObjectIDFromHex(userIdString)

	if err != nil {
		response.Failed(c, http.StatusUnauthorized, errors.ErrInternalServer.Error())
		return
	}

	if err := handler.answerService.RateUp(ctx, payload.AnswerID, userId); err != nil {
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	response.OK(c, http.StatusOK)
}
