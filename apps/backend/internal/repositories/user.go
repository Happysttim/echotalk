package repositories

import (
	"context"
	"echotalk/internal/config"
	"echotalk/internal/database"
	"echotalk/internal/model"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context) ([]*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByLocal(ctx context.Context, email string, passwordHash string) (*model.User, error)
	FindByProvider(ctx context.Context, authProvider model.AuthProvider, providerID string) (*model.User, error)
}

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository() UserRepository {
	return &MongoUserRepository{collection: database.Database.Collection(config.CollectionUser)}
}

func (r *MongoUserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	doc := user
	doc.ID = result.InsertedID.(bson.ObjectID)
	return doc, nil
}

func (r *MongoUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *MongoUserRepository) FindByProvider(ctx context.Context, authProvider model.AuthProvider, providerID string) (*model.User, error) {
	var user model.User
	if err := r.collection.FindOne(ctx, bson.M{"auth_provider": authProvider, "provider_id": providerID}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *MongoUserRepository) FindByLocal(ctx context.Context, email string, passwordHash string) (*model.User, error) {
	var user model.User
	if err := r.collection.FindOne(ctx, bson.M{"email": email, "password_hash": passwordHash}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *MongoUserRepository) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *MongoUserRepository) FindAll(ctx context.Context) ([]*model.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	users := make([]*model.User, 0)
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
