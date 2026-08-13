package database

import (
	"context"
	"echotalk/internal/config"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

const connectTimeout = 5 * time.Second

var Client *mongo.Client
var Database *mongo.Database

func init() {
	if config.Config.MongoURI == "" {
		panic("invalid mongouri")
	}

	ctx := context.Background()

	client, err := mongo.Connect(options.Client().ApplyURI(config.Config.MongoURI))
	if err != nil {
		panic("connect mongodb error:" + err.Error())
	}

	pingCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		_ = client.Disconnect(context.Background())
		panic("ping mongodb error: " + err.Error())
	}

	Client = client
}

func InitDatabase() error {
	if Client == nil {
		return errors.New("mongo client is nil")
	}

	if config.Config.DBName == "" {
		return errors.New("invalid mongo database name")
	}

	Database = Client.Database(config.Config.DBName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	return initSessionIndexes(ctx)
}

func initSessionIndexes(ctx context.Context) error {
	collection := Database.Collection(config.CollectionSession)
	_, err := collection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.D{
				{Key: "expires_at", Value: 1},
			},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
	)

	return err
}

func Close() error {
	ctx := context.Background()
	if Client == nil {
		panic("mongo client is nil")
	}
	return Client.Disconnect(ctx)
}
