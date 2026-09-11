package handler

import (
	"context"
	"echotalk/internal/errors"
	"echotalk/internal/model"
	"echotalk/internal/response"
	"echotalk/internal/service"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type SurveyCursor struct {
	SurveyID  string `json:"survey_id"`
	UpdatedAt string `json:"updated_at"`
}

type SurveyHandler struct {
	surveyService *service.SurveyService
	answerService *service.AnswerService
	userService   *service.UserService
}

func NewSurveyHandler(surveyService *service.SurveyService, answerService *service.AnswerService, userService *service.UserService) *SurveyHandler {
	return &SurveyHandler{
		surveyService: surveyService,
		answerService: answerService,
		userService:   userService,
	}
}

// @Router /surveys [post]
func (handler *SurveyHandler) CreateSurvey(c *gin.Context) {
	payload := new(model.CreateSurveyRequest)
	userID := c.GetString("userID")

	if userID == "" {
		response.Failed(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error())
		return
	}
	if err := c.ShouldBindJSON(payload); err != nil {
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()

	survey, err := handler.surveyService.CreateSurvey(ctx, payload, userID)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, errors.ErrInternalServer.Error())
		return
	}

	response.OKWithData(c, http.StatusCreated, survey)
}

// @Router /surveys [delete]
func (handler *SurveyHandler) DeleteSurvey(c *gin.Context) {
	payload := new(model.DeleteSurveyRequest)
	userID := c.GetString("userID")

	if userID == "" {
		response.Failed(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error())
		return
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()
	survey, err := handler.surveyService.GetSurveyByID(ctx, payload.SurveyID)

	if err != nil {
		response.Failed(c, http.StatusInternalServerError, errors.ErrInternalServer.Error())
		return
	}

	if survey == nil {
		response.Failed(c, http.StatusNotFound, errors.ErrNotFound.Error())
		return
	}

	if strings.Compare(survey.AuthorID.Hex(), userID) != 0 {
		response.Failed(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error())
		return
	}

	if err := handler.answerService.DeleteAnswerBySurvey(ctx, survey.ID.Hex()); err != nil {
		response.Failed(c, http.StatusInternalServerError, errors.ErrInternalServer.Error())
		return
	}

	if err := handler.surveyService.DeleteSurvey(ctx, payload.SurveyID); err != nil {
		response.Failed(c, http.StatusInternalServerError, errors.ErrInternalServer.Error())
		return
	}

	response.OK(c, http.StatusOK)
}

// @Router /surveys/:surveyId [get]
func (handler *SurveyHandler) GetSurvey(c *gin.Context) {
	surveyId := c.Param("surveyId")
	if surveyId == "" {
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()
	survey, err := handler.surveyService.GetSurveyByID(ctx, surveyId)

	if err != nil {
		response.Failed(c, http.StatusInternalServerError, errors.ErrInternalServer.Error())
		return
	}

	response.OKWithData(c, http.StatusOK, survey)
}

// @Router /surveys [patch]
func (handler *SurveyHandler) UpdateSurvey(c *gin.Context) {
	payload := new(model.UpdateSurveyRequest)
	userID := c.GetString("userID")

	if userID == "" {
		response.Failed(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error())
		return
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
		return
	}

	ctx := c.Request.Context()
	survey, err := handler.surveyService.GetSurveyByID(ctx, payload.SurveyID)

	if err != nil {
		response.Failed(c, http.StatusInternalServerError, errors.ErrInternalServer.Error())
		return
	}

	if strings.Compare(survey.AuthorID.Hex(), userID) != 0 {
		response.Failed(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error())
		return
	}

	if err := handler.surveyService.UpdateSurvey(ctx, payload); err != nil {
		response.Failed(c, http.StatusInternalServerError, errors.ErrInternalServer.Error())
		return
	}

	response.OK(c, http.StatusOK)
}

// @Router /surveys/me?type={string}&cursor={string}&limit={number} [get]
func (handler *SurveyHandler) GetMySurveyPage(c *gin.Context) {
	handler.surveyPage(c, true)
}

// @Router /surveys?type={string}&cursor={string}&limit={number} [get]
func (handler *SurveyHandler) GetSurveyPage(c *gin.Context) {
	handler.surveyPage(c, false)
}

func (handler *SurveyHandler) surveyPage(c *gin.Context, me bool) {
	findType := c.Query("type")
	cursorBase64 := c.Query("cursor")
	limit := c.Query("limit")

	filter := bson.M{"is_public": true}
	sort := bson.D{}

	if me {
		userID := c.GetString("userID")
		objectID, err := bson.ObjectIDFromHex(userID)

		if err != nil {
			response.Failed(c, http.StatusUnauthorized, errors.ErrUnauthorized.Error())
			return
		}

		filter["author_id"] = objectID
	}

	if !slices.Contains(SurveyTypes, findType) {
		findType = TypeID
	}

	if findType == TypeUpdatedAt {
		sort = bson.D{
			{Key: "updated_at", Value: -1},
			{Key: "_id", Value: -1},
		}
	} else {
		sort = bson.D{
			{Key: "_id", Value: -1},
		}
	}

	if cursorBase64 != "" {
		decodeBytes, err := base64.RawURLEncoding.DecodeString(cursorBase64)
		if err != nil {
			response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
			return
		}

		cursor := new(SurveyCursor)
		if err := json.Unmarshal(decodeBytes, cursor); err != nil {
			response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
			return
		}

		if findType == TypeUpdatedAt {
			surveyId, err := bson.ObjectIDFromHex(cursor.SurveyID)
			if err != nil {
				response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
				return
			}

			updatedAt, err := time.Parse(time.RFC3339Nano, cursor.UpdatedAt)
			if err != nil {
				response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
				return
			}

			filter["$or"] = bson.A{
				bson.M{
					"updated_at": updatedAt,
					"_id": bson.M{
						"$lt": surveyId,
					},
				},
				bson.M{
					"updated_at": bson.M{
						"$lt": updatedAt,
					},
				},
			}
		} else {
			surveyID, err := bson.ObjectIDFromHex(cursor.SurveyID)
			if err != nil {
				response.Failed(c, http.StatusBadRequest, errors.ErrBadRequest.Error())
				return
			}
			filter["_id"] = bson.M{
				"$lt": surveyID,
			}
		}
	}

	limitNumber, err := strconv.Atoi(limit)

	if err != nil {
		limitNumber = LimitMedium
	}

	if !slices.Contains(Limits, limitNumber) {
		limitNumber = LimitMedium
	}

	ctx := c.Request.Context()
	surveys, err := handler.surveyService.GetFilteredSurveys(
		ctx,
		filter,
		options.Find().SetSort(sort).SetLimit(int64(limitNumber)+1),
	)

	if err != nil {
		response.Failed(c, http.StatusInternalServerError, errors.ErrInternalServer.Error())
		return
	}

	var nextSkip string
	var nextID string
	var nextAt string
	hasNext := len(surveys) > limitNumber

	if hasNext {
		last := surveys[limitNumber-1]
		nextID = last.ID.Hex()
		nextAt = last.UpdatedAt.Format(time.RFC3339Nano)
		surveys = surveys[:limitNumber]
	}

	if nextID != "" && nextAt != "" {
		nextCursor := SurveyCursor{
			SurveyID:  nextID,
			UpdatedAt: nextAt,
		}

		jsonBytes, err := json.Marshal(nextCursor)
		if err != nil {
			response.Failed(c, http.StatusInternalServerError, errors.ErrInternalServer.Error())
			return
		}

		nextSkip = base64.RawURLEncoding.EncodeToString(jsonBytes)
	}

	response.OKWithData(c, http.StatusOK, map[string]any{"skip": nextSkip, "hasNext": hasNext, "surveys": surveys})
}

func (handler *SurveyHandler) Ticker(ctx context.Context) {
	handler.surveyService.CheckExpire(ctx)
}
