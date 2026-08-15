package repositories

import (
	"context"
	"echotalk/internal/model"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RateUpRepository interface {
	RateUp(ctx context.Context, rateUp *model.RateUp) error
}

type MongoRateUpRepository struct {
	collection *mongo.Collection
}

func NewRateUpRepository(collection *mongo.Collection) RateUpRepository {
	return &MongoRateUpRepository{collection: collection}
}

func (r *MongoRateUpRepository) RateUp(ctx context.Context, rateUp *model.RateUp) error {
	_, err := r.collection.InsertOne(ctx, rateUp)

	return err
}
