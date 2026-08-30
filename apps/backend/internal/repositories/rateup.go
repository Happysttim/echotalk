package repositories

import (
	"context"
	"echotalk/internal/config"
	"echotalk/internal/database"
	"echotalk/internal/model"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RateUpRepository interface {
	RateUp(ctx context.Context, rateUp *model.RateUp) error
}

type MongoRateUpRepository struct {
	collection *mongo.Collection
}

func NewRateUpRepository() RateUpRepository {
	return &MongoRateUpRepository{collection: database.Database.Collection(config.CollectionRateUp)}
}

func (r *MongoRateUpRepository) RateUp(ctx context.Context, rateUp *model.RateUp) error {
	_, err := r.collection.InsertOne(ctx, rateUp)

	return err
}
