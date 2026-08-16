package handler

import (
	"echotalk/internal/model"
	"echotalk/internal/service"
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	LimitSmall  = 10
	LimitMedium = 30
	LimitLarge  = 30
)

const (
	TypeID        = "_id"
	TypeUpdatedAt = "updated_at"
)

var Limits = []int{LimitSmall, LimitMedium, LimitLarge}

type SurveyHandler struct {
	surveyService *service.SurveyService
	userService   *service.UserService
}

func NewSurveyHandler(surveyService *service.SurveyService, userService *service.UserService) *SurveyHandler {
	return &SurveyHandler{
		surveyService: surveyService,
		userService:   userService,
	}
}

// @Router /surveys [post]
func (handler *SurveyHandler) CreateSurvey(c *gin.Context) {
	payload := new(model.CreateSurveyRequest)
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

	survey, err := handler.surveyService.CreateSurvey(ctx, payload, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created", "survey": survey})
}

// @Router /surveys/:surveyId [get]
func (handler *SurveyHandler) GetSurvey(c *gin.Context) {
	surveyId := c.Param("surveyId")
	if surveyId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	ctx := c.Request.Context()
	survey, err := handler.surveyService.GetSurveyByID(ctx, surveyId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "survey": survey})
}

// @Router /surveys?skip={string}&limit={number} [get]
func (handler *SurveyHandler) GetSurveyPage(c *gin.Context) {
	skip := c.Query("skip")
	limit := c.Query("number")

	limitNumber, err := strconv.Atoi(limit)

	if err != nil {
		limitNumber = LimitMedium
	}

	if !slices.Contains(Limits, limitNumber) {
		limitNumber = LimitMedium
	}

	ctx := c.Request.Context()
	filter := bson.M{
		"_id": bson.M{
			"$lt": skip,
		},
	}
	surveys, err := handler.surveyService.GetFilteredSurveys(
		ctx,
		filter,
		options.Find().SetSort(bson.D{
			{Key: "_id", Value: -1},
		}).SetLimit(int64(limitNumber)),
	)

	c.JSON(http.StatusOK, gin.H{"status": "ok", "skip": surveys[len(surveys)-1].ID.Hex(), "surveys": surveys})
}
