package repositories

import (
	"context"
	"echotalk/internal/model"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type SurveyRepository interface {
	Create(ctx context.Context, survey *model.Survey) (*model.Survey, error)
	FindByID(ctx context.Context, id string) (*model.Survey, error)
	Update(ctx context.Context, survey *model.Survey) error
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context) ([]*model.Survey, error)
	FindByFilter(ctx context.Context, filter bson.M, opts ...options.Lister[options.FindOptions]) ([]*model.Survey, error)
}

type MongoSurveyRepository struct {
	collection *mongo.Collection
}

func NewMongoSurveyRepository(collection *mongo.Collection) SurveyRepository {
	return &MongoSurveyRepository{collection: collection}
}

func (r *MongoSurveyRepository) Create(ctx context.Context, survey *model.Survey) (*model.Survey, error) {
	result, err := r.collection.InsertOne(ctx, survey)
	if err != nil {
		return nil, err
	}

	doc := *survey
	doc.ID = result.InsertedID.(bson.ObjectID)
	return &doc, nil
}

func (r *MongoSurveyRepository) Update(ctx context.Context, survey *model.Survey) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": survey.ID}, bson.M{"$set": survey})
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoSurveyRepository) FindByID(ctx context.Context, id string) (*model.Survey, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var survey model.Survey
	if err := r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&survey); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &survey, nil
}

func (r *MongoSurveyRepository) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *MongoSurveyRepository) FindByFilter(ctx context.Context, filter bson.M, opts ...options.Lister[options.FindOptions]) ([]*model.Survey, error) {
	cursor, err := r.collection.Find(ctx, filter, opts...)
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

func (r *MongoSurveyRepository) FindAll(ctx context.Context) ([]*model.Survey, error) {
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
