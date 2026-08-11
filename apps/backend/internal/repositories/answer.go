package repositories

import (
	"context"
	"echotalk/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AnswerRepository interface {
	Create(answer *model.Answer) (*model.Answer, error)
	Update(answer *model.Answer) error
	FindByID(id string) (*model.Answer, error)
	Delete(id string) error
	FindAll() ([]*model.Answer, error)
}

type MongoAnswerRepository struct {
	collection *mongo.Collection
}

func NewMongoAnswerRepository(collection *mongo.Collection) AnswerRepository {
	return &MongoAnswerRepository{collection: collection}
}

func (r *MongoAnswerRepository) Create(answer *model.Answer) (*model.Answer, error) {
	ctx := context.Background()
	result, err := r.collection.InsertOne(ctx, answer)
	if err != nil {
		return nil, err
	}

	doc := *answer
	doc.ID = result.InsertedID.(bson.ObjectID)
	return &doc, nil
}

func (r *MongoAnswerRepository) Update(answer *model.Answer) error {
	ctx := context.Background()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": answer.ID}, bson.M{"$set": answer})
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoAnswerRepository) FindByID(id string) (*model.Answer, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	var answer model.Answer

	if err := r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&answer); err != nil {
		return nil, err
	}

	return &answer, nil
}

func (r *MongoAnswerRepository) Delete(id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *MongoAnswerRepository) FindAll() ([]*model.Answer, error) {
	ctx := context.Background()
	cursor, err := r.collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var answers []*model.Answer
	if err := cursor.All(ctx, &answers); err != nil {
		return nil, err
	}
	return answers, nil

}
