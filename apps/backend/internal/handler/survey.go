package handler

import (
	"echotalk/internal/model"
	"echotalk/internal/service"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
	"time"

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
var Types = []string{TypeID, TypeUpdatedAt}

type SurveyCursor struct {
	SurveyID  string `json:"survey_id"`
	UpdatedAt string `json:"updated_at"`
}

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

// @Router /surveys?type={string}&cursor={string}&limit={number} [get]
func (handler *SurveyHandler) GetSurveyPage(c *gin.Context) {
	findType := c.Query("type")
	cursorBase64 := c.Query("cursor")
	limit := c.Query("limit")

	filter := bson.M{}
	sort := bson.D{}

	if !slices.Contains(Types, findType) {
		findType = TypeID
	}

	if findType == TypeUpdatedAt {
		sort = bson.D{
			{Key: "_id", Value: -1},
			{Key: "updated_at", Value: -1},
		}
	} else {
		sort = bson.D{
			{Key: "_id", Value: -1},
		}
	}

	if cursorBase64 != "" {
		decodeBytes, err := base64.RawURLEncoding.DecodeString(cursorBase64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}

		var cursor SurveyCursor
		if err := json.Unmarshal(decodeBytes, &cursor); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}

		if findType == TypeUpdatedAt {
			surveyID, err := bson.ObjectIDFromHex(cursor.SurveyID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
				return
			}

			updatedAt, err := time.Parse(time.RFC3339Nano, cursor.UpdatedAt)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
				return
			}

			filter["$or"] = bson.A{
				bson.M{
					"updated_at": updatedAt,
					"_id": bson.M{
						"$lt": surveyID,
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
			if err == nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
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
		options.Find().SetSort(sort).SetLimit(int64(limitNumber)),
	)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var nextSkip string
	var nextID string
	var nextAt string

	if len(surveys) > 0 {
		nextID = surveys[len(surveys)-1].ID.Hex()
		nextAt = surveys[len(surveys)-1].UpdatedAt.Format(time.RFC3339Nano)
	}

	if nextID != "" && nextAt != "" {
		nextCursor := SurveyCursor{
			SurveyID:  nextID,
			UpdatedAt: nextAt,
		}

		jsonBytes, err := json.Marshal(nextCursor)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		nextSkip = base64.RawURLEncoding.EncodeToString(jsonBytes)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "skip": nextSkip, "surveys": surveys})
}
