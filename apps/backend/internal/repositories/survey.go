package repositories

import (
	"context"
	"echotalk/internal/config"
	"echotalk/internal/database"
	"echotalk/internal/model"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SurveyRepository interface {
	Create(ctx context.Context, survey *model.Survey) (*model.Survey, error)
	FindByID(ctx context.Context, id string) (*model.Survey, error)
	Update(ctx context.Context, survey *model.Survey) error
	CheckExpires(ctx context.Context) error
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context) ([]*model.Survey, error)
	FindDetailByID(ctx context.Context, id string) (*model.SurveyDetail, error)
	FindByFilter(ctx context.Context, filter bson.D, sort bson.D, limit int64) ([]*model.SurveyDetail, error)
}

type MongoSurveyRepository struct {
	collection *mongo.Collection
}

func NewMongoSurveyRepository() SurveyRepository {
	return &MongoSurveyRepository{collection: database.Database.Collection(config.CollectionSurvey)}
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

func (r *MongoSurveyRepository) CheckExpires(ctx context.Context) error {
	_, err := r.collection.UpdateMany(ctx, bson.M{
		"closed":     false,
		"expires_at": bson.M{"$lte": time.Now()},
	}, bson.M{
		"$set": bson.M{
			"closed": true,
		},
	})

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

func (r *MongoSurveyRepository) FindByFilter(ctx context.Context, filter bson.D, sort bson.D, limit int64) ([]*model.SurveyDetail, error) {
	cursor, err := r.collection.Aggregate(ctx, r.pipeline(filter, sort, limit))
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	surveys := make([]*model.SurveyDetail, 0)

	if err := cursor.All(ctx, &surveys); err != nil {
		return nil, err
	}

	return surveys, nil
}

func (r *MongoSurveyRepository) FindDetailByID(ctx context.Context, id string) (*model.SurveyDetail, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Aggregate(ctx, r.pipeline(
		bson.D{{Key: "_id", Value: objectID}},
		bson.D{},
		1,
	))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var survey model.SurveyDetail
	if !cursor.Next(ctx) {
		if err := cursor.Err(); err != nil {
			return nil, err
		}
		return nil, nil
	}

	if err := cursor.Decode(&survey); err != nil {
		return nil, err
	}

	return &survey, nil
}

func (r *MongoSurveyRepository) FindAll(ctx context.Context) ([]*model.Survey, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	surveys := make([]*model.Survey, 0)
	if err := cursor.All(ctx, &surveys); err != nil {
		return nil, err
	}
	return surveys, nil
}

func (r *MongoSurveyRepository) pipeline(match bson.D, sort bson.D, limit int64) mongo.Pipeline {
	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: match},
		},
	}

	if len(sort) > 0 {
		pipeline = append(pipeline, bson.D{
			{Key: "$sort", Value: sort},
		})
	}

	if limit > 0 {
		pipeline = append(pipeline, bson.D{
			{Key: "$limit", Value: limit},
		})
	}

	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: config.CollectionUser},
		{Key: "localField", Value: "author_id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "author"},
	}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$author"},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}}
	pipeline = append(pipeline, lookupStage, unwindStage)

	return pipeline
}
