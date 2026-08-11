package repositories

import (
	"context"
	"echotalk/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SurveyRepository interface {
	Create(survey *model.Survey) (*model.Survey, error)
	FindByID(id string) (*model.Survey, error)
	Update(survey *model.Survey) error
	Delete(id string) error
	FindAll() ([]*model.Survey, error)
	FindByFilter(filter bson.M) ([]*model.Survey, error)
}

type MongoSurveyRepository struct {
	collection *mongo.Collection
}

func NewMongoSurveyRepository(collection *mongo.Collection) SurveyRepository {
	return &MongoSurveyRepository{collection: collection}
}

func (r *MongoSurveyRepository) Create(survey *model.Survey) (*model.Survey, error) {
	ctx := context.Background()

	result, err := r.collection.InsertOne(ctx, survey)
	if err != nil {
		return nil, err
	}

	doc := *survey
	doc.ID = result.InsertedID.(bson.ObjectID)
	return &doc, nil
}

func (r *MongoSurveyRepository) Update(survey *model.Survey) error {
	ctx := context.Background()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": survey.ID}, bson.M{"$set": survey})
	if err != nil {
		return err
	}

	return nil
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

func (r *MongoSurveyRepository) FindByFilter(filter bson.M) ([]*model.Survey, error) {
	ctx := context.Background()
	cursor, err := r.collection.Find(ctx, filter)
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
