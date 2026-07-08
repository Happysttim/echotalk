package repositories

import (
	"context"
	"echotalk/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	Create(user *model.CreateUserRequest) (*model.User, error)
	FindByID(id string) (*model.User, error)
	Delete(id string) error
	FindAll() ([]*model.User, error)
}

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(collection *mongo.Collection) UserRepository {
	return &MongoUserRepository{collection: collection}
}

func (r *MongoUserRepository) Create(user *model.CreateUserRequest) (*model.User, error) {
	ID := bson.NewObjectID()

	doc := model.User{
		ID:        ID,
		GoogleID:  user.GoogleID,
		Email:     user.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ctx := context.Background()
	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}
	doc.ID = result.InsertedID.(bson.ObjectID)
	return &doc, nil
}

func (r *MongoUserRepository) FindByID(id string) (*model.User, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	var user model.User
	if err := r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *MongoUserRepository) Delete(id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *MongoUserRepository) FindAll() ([]*model.User, error) {
	ctx := context.Background()
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
