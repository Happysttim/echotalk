package repositories

import (
	"context"
	"echotalk/internal/model"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SessionRepository interface {
	Create(ctx context.Context, session *model.Session) (*model.Session, error)
	FindByID(ctx context.Context, id string) (*model.Session, error)
	FindByToken(ctx context.Context, hashedToken string) (*model.Session, error)
	Update(ctx context.Context, session *model.Session) error
}

type MongoSessionRepository struct {
	collection *mongo.Collection
}

func NewMongoSessionRepository(collection *mongo.Collection) SessionRepository {
	return &MongoSessionRepository{collection: collection}
}

func (r *MongoSessionRepository) Create(ctx context.Context, session *model.Session) (*model.Session, error) {
	result, err := r.collection.InsertOne(ctx, session)
	if err != nil {
		return nil, err
	}

	doc := *session
	doc.ID = result.InsertedID.(bson.ObjectID)
	return &doc, nil
}

func (r *MongoSessionRepository) FindByID(ctx context.Context, id string) (*model.Session, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var session model.Session
	if err := r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&session); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &session, nil
}

func (r *MongoSessionRepository) FindByToken(ctx context.Context, hashedToken string) (*model.Session, error) {
	var session model.Session
	if err := r.collection.FindOne(ctx, bson.M{"refresh_token_hash": hashedToken}).Decode(&session); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &session, nil
}

func (r *MongoSessionRepository) Update(ctx context.Context, session *model.Session) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": session.ID}, bson.M{"$set": session})
	if err != nil {
		return err
	}

	return nil
}
