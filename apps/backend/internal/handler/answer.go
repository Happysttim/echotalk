package handler

import (
	"echotalk/internal/errors"
	"echotalk/internal/model"
	"echotalk/internal/response"
	"echotalk/internal/service"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

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
		log.Println("User ID is empty")
		response.Failed(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error())
		return
	}
	if err := c.ShouldBindJSON(payload); err != nil {
		log.Println("Error binding JSON: ", err)
		response.Failed(c, http.StatusUnauthorized, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()
	survey, err := handler.surveyService.GetSurveyByID(ctx, payload.SurveyID)

	if err != nil || survey == nil {
		log.Println("Error fetching survey")
		response.Failed(c, http.StatusNotFound, errors.ErrSurveyNotFound.Error())
		return
	}

	if survey.Closed || survey.ExpiresAt.Before(time.Now()) {
		log.Println("Survey is closed")
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	answer, err := handler.answerService.CreateAnswer(ctx, payload, userID)
	if err != nil {
		log.Println("Error creating answer: ", err)
		response.Failed(c, http.StatusUnauthorized, errors.ErrInternalServer.Error())
		return
	}

	response.OKWithData(c, http.StatusCreated, answer)
}

// @Router /answers [delete]
func (handler *AnswerHandler) DeleteAnswer(c *gin.Context) {
	payload := new(model.DeleteAnswerRequest)
	userId := c.GetString("userID")

	if err := c.ShouldBindJSON(payload); err != nil {
		log.Println("Error binding JSON: ", err)
		response.Failed(c, http.StatusUnauthorized, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()
	answer, err := handler.answerService.GetAnswerByID(ctx, payload.AnswerID)

	if err != nil || answer == nil {
		log.Println("Error fetching answer")
		response.Failed(c, http.StatusNotFound, errors.ErrAnswerNotFound.Error())
		return
	}

	if strings.Compare(answer.AuthorID.Hex(), userId) != 0 {
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	if err := handler.answerService.DeleteAnswer(ctx, answer.ID.Hex()); err != nil {
		log.Println("Error deleting answer: ", err)
		response.Failed(c, http.StatusUnauthorized, errors.ErrInternalServer.Error())
		return
	}

	response.OK(c, http.StatusOK)
}

// @Router /answers/:answerId [get]
func (handler *AnswerHandler) GetAnswer(c *gin.Context) {
	answerId := c.Param("answerId")
	if answerId == "" {
		log.Println("Answer ID is empty")
		response.Failed(c, http.StatusUnauthorized, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()
	answer, err := handler.answerService.GetAnswerByID(ctx, answerId)

	if err != nil || answer == nil {
		log.Println("Error fetching answer")
		response.Failed(c, http.StatusNotFound, errors.ErrAnswerNotFound.Error())
		return
	}

	response.OKWithData(c, http.StatusOK, answer)
}

// @Router /answers [patch]
func (handler *AnswerHandler) UpdateAnswer(c *gin.Context) {
	payload := new(model.UpdateAnswerRequest)
	userID := c.GetString("userID")

	if userID == "" {
		log.Println("User ID is empty")
		response.Failed(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error())
		return
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println("Error binding JSON: ", err)
		response.Failed(c, http.StatusUnauthorized, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()
	answer, err := handler.answerService.GetAnswerByID(ctx, payload.AnswerID)

	if err != nil || answer == nil {
		log.Println("Error fetching answer")
		response.Failed(c, http.StatusNotFound, errors.ErrAnswerNotFound.Error())
		return
	}

	if strings.Compare(answer.AuthorID.Hex(), userID) != 0 {
		log.Println("User ID does not match answer author ID")
		response.Failed(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error())
		return
	}

	if err := handler.answerService.UpdateAnswer(ctx, payload); err != nil {
		log.Println("Error updating answer: ", err)
		response.Failed(c, http.StatusUnauthorized, errors.ErrInternalServer.Error())
		return
	}

	response.OK(c, http.StatusOK)
}

// @Router /answers?cursor={string} [get]
func (handler *AnswerHandler) GetAnswerFeed(c *gin.Context) {
	cursorBase64 := c.Query("cursor")
	var survey *model.Survey
	var answer *model.Answer
	var answers []*model.Answer
	var err error
	ctx := c.Request.Context()

	if cursorBase64 == "" {
		log.Println("Cursor is empty")
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	decodeBytes, err := base64.RawURLEncoding.DecodeString(cursorBase64)
	if err != nil {
		log.Println("Error decoding cursor: ", err)
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	cursor := new(AnswerCursor)
	if err := json.Unmarshal(decodeBytes, cursor); err != nil {
		log.Println("Error unmarshaling cursor: ", err)
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	survey, err = handler.surveyService.GetSurveyByID(ctx, cursor.SurveyID)
	if err != nil {
		log.Println("Error fetching survey: ", err)
		response.Failed(c, http.StatusInternalServerError, errors.ErrInternalServer.Error())
		return
	}
	if survey == nil {
		log.Println("Survey not found")
		response.Failed(c, http.StatusNotFound, errors.ErrNotFound.Error())
		return
	}

	if cursor.AnswerID != "" {
		answer, err = handler.answerService.GetAnswerByID(ctx, cursor.AnswerID)
		if err != nil {
			log.Println("Error fetching answer: ", err)
			response.Failed(c, http.StatusInternalServerError, errors.ErrInternalServer.Error())
			return
		}
		if answer == nil || answer.SurveyID != survey.ID {
			log.Println("Answer not found or does not belong to the specified survey")
			response.Failed(c, http.StatusNotFound, errors.ErrNotFound.Error())
			return
		}
		answers, err = handler.answerService.GetFilteredAnswers(
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
	} else {
		answers, err = handler.answerService.GetFilteredAnswers(
			ctx,
			bson.M{
				"survey_id": survey.ID,
			},
			options.Find().SetSort(
				bson.D{
					{Key: "_id", Value: -1},
					{Key: "rate_up", Value: -1},
				},
			).SetLimit(LimitSmall+1),
		)
	}

	var nextSkip string
	var nextID string
	hasNext := len(answers) > LimitSmall
	nextCursor := AnswerCursor{
		SurveyID: cursor.SurveyID,
	}
	if hasNext && len(answers) > 0 {
		last := answers[len(answers)-1]
		nextID = last.ID.Hex()
	}

	if nextID != "" {
		nextCursor.AnswerID = nextID
	}

	if nextCursor.AnswerID != "" {
		jsonBytes, err := json.Marshal(nextCursor)
		if err != nil {
			log.Println("Error marshaling next cursor: ", err)
			response.Failed(c, http.StatusInternalServerError, errors.ErrInternalServer.Error())
			return
		}
		nextSkip = base64.RawURLEncoding.EncodeToString(jsonBytes)
	}

	response.OKWithData(c, http.StatusOK, map[string]any{"skip": nextSkip, "hasNext": hasNext, "answers": answers})
}

// @Router /answers/rateup [post]
func (handler *AnswerHandler) RateUp(c *gin.Context) {
	payload := new(model.RateUpRequest)
	userIdString := c.GetString("userID")

	if err := c.ShouldBindJSON(payload); err != nil {
		log.Println("Error binding JSON: ", err)
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	if userIdString == "" {
		log.Println("User ID is empty")
		response.Failed(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error())
		return
	}

	ctx := c.Request.Context()
	userId, err := bson.ObjectIDFromHex(userIdString)

	if err != nil {
		log.Println("Error converting user ID to ObjectID: ", err)
		response.Failed(c, http.StatusUnauthorized, errors.ErrInternalServer.Error())
		return
	}

	if err := handler.answerService.RateUp(ctx, payload.AnswerID, userId); err != nil {
		log.Println("Error rating up answer: ", err)
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	response.OK(c, http.StatusOK)
}
