package repositories

import (
	"context"
	"echotalk/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SurveyRepository interface {
	Create(survey *model.CreateSurveyRequest) (*model.Survey, error)
	FindByID(id string) (*model.Survey, error)
	Delete(id string) error
	FindAll() ([]*model.Survey, error)
}

type MongoSurveyRepository struct {
	collection *mongo.Collection
}

func NewMongoSurveyRepository(collection *mongo.Collection) *MongoSurveyRepository {
	return &MongoSurveyRepository{collection: collection}
}

func (r *MongoSurveyRepository) Create(survey *model.CreateSurveyRequest) (*model.Survey, error) {
	ID := bson.NewObjectID()
	ctx := context.Background()

	var user model.User
	if err := r.collection.Database().Collection("users").FindOne(ctx, bson.M{"ID": survey.Author}).Decode(&user); err != nil {
		return nil, err
	}

	doc := model.Survey{
		ID:        ID,
		Title:     survey.Title,
		Content:   survey.Content,
		Author:    user,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	doc.ID = result.InsertedID.(bson.ObjectID)
	return &doc, nil
}
