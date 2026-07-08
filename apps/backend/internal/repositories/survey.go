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

func NewMongoSurveyRepository(collection *mongo.Collection) SurveyRepository {
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

func (r *MongoSurveyRepository) FindByID(id string) (*model.Survey, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	var survey model.Survey
	if err := r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&survey); err != nil {
		return nil, err
	}

	return &survey, nil
}

func (r *MongoSurveyRepository) Delete(id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *MongoSurveyRepository) FindAll() ([]*model.Survey, error) {
	ctx := context.Background()
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var surveys []*model.Survey
	if err := cursor.All(ctx, &surveys); err != nil {
		return nil, err
	}
	return surveys, nil
}
