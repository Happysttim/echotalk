package repositories

import (
	"context"
	"echotalk/internal/model"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AnswerRepository interface {
	Create(ctx context.Context, answer *model.Answer) (*model.Answer, error)
	Update(ctx context.Context, answer *model.Answer) error
	FindByID(ctx context.Context, id string) (*model.Answer, error)
	FindBySurveyID(ctx context.Context, id string) ([]*model.Answer, error)
	Delete(ctx context.Context, id string) error
	DeleteBySurveyID(ctx context.Context, id string) error
	FindAll(ctx context.Context) ([]*model.Answer, error)
	FindByFilter(ctx context.Context, filter bson.M, opts ...options.Lister[options.FindOptions]) ([]*model.Answer, error)
}

type MongoAnswerRepository struct {
	collection *mongo.Collection
}

func NewMongoAnswerRepository(collection *mongo.Collection) AnswerRepository {
	return &MongoAnswerRepository{collection: collection}
}

func (r *MongoAnswerRepository) Create(ctx context.Context, answer *model.Answer) (*model.Answer, error) {
	result, err := r.collection.InsertOne(ctx, answer)
	if err != nil {
		return nil, err
	}

	doc := *answer
	doc.ID = result.InsertedID.(bson.ObjectID)
	return &doc, nil
}

func (r *MongoAnswerRepository) Update(ctx context.Context, answer *model.Answer) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": answer.ID}, bson.M{"$set": answer})
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoAnswerRepository) FindByID(ctx context.Context, id string) (*model.Answer, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var answer model.Answer

	if err := r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&answer); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &answer, nil
}

func (r *MongoAnswerRepository) FindBySurveyID(ctx context.Context, id string) ([]*model.Answer, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	cursor, err := r.collection.Find(ctx, bson.M{"survey_id": objectID})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	answers := make([]*model.Answer, 0)
	if err := cursor.All(ctx, &answers); err != nil {
		return nil, err
	}
	return answers, nil

}

func (r *MongoAnswerRepository) FindByFilter(ctx context.Context, filter bson.M, opts ...options.Lister[options.FindOptions]) ([]*model.Answer, error) {
	cursor, err := r.collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	answers := make([]*model.Answer, 0)

	if err := cursor.All(ctx, answers); err != nil {
		return nil, err
	}

	return answers, nil
}

func (r *MongoAnswerRepository) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *MongoAnswerRepository) DeleteBySurveyID(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"survey_id": objectID})
	return err
}

func (r *MongoAnswerRepository) FindAll(ctx context.Context) ([]*model.Answer, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	answers := make([]*model.Answer, 0)
	if err := cursor.All(ctx, &answers); err != nil {
		return nil, err
	}
	return answers, nil

}
