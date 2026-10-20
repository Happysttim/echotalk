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

type AnswerRepository interface {
	Create(ctx context.Context, answer *model.Answer) (*model.Answer, error)
	Update(ctx context.Context, answer *model.Answer) error
	RateUp(ctx context.Context, answerId string) error
	FindByID(ctx context.Context, id string) (*model.Answer, error)
	FindDetailByID(ctx context.Context, id string) (*model.AnswerDetail, error)
	FindBySurveyID(ctx context.Context, id string) ([]*model.AnswerDetail, error)
	Delete(ctx context.Context, id string) error
	DeleteBySurveyID(ctx context.Context, id string) error
	FindAll(ctx context.Context) ([]*model.AnswerDetail, error)
	FindByFilter(ctx context.Context, filter bson.D, sort bson.D, limit int64) ([]*model.AnswerDetail, error)
	FindMyAnswerDetails(ctx context.Context, filter bson.D, sort bson.D, limit int64) ([]*model.MyAnswerDetail, error)
	pipeline(match bson.D, sort bson.D, limit int64) mongo.Pipeline
}

type MongoAnswerRepository struct {
	collection *mongo.Collection
}

func NewMongoAnswerRepository() AnswerRepository {
	return &MongoAnswerRepository{collection: database.Database.Collection(config.CollectionAnswer)}
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

func (r *MongoAnswerRepository) RateUp(ctx context.Context, answerId string) error {
	objectID, err := bson.ObjectIDFromHex(answerId)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{
		"$inc": bson.M{
			"rate_up": 1,
		},
	})

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

func (r *MongoAnswerRepository) FindBySurveyID(ctx context.Context, id string) ([]*model.AnswerDetail, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	cursor, err := r.collection.Aggregate(ctx, r.pipeline(
		bson.D{{Key: "survey_id", Value: objectID}},
		bson.D{},
		0,
	))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	answers := make([]*model.AnswerDetail, 0)
	if err := cursor.All(ctx, &answers); err != nil {
		return nil, err
	}
	return answers, nil

}

func (r *MongoAnswerRepository) FindDetailByID(ctx context.Context, id string) (*model.AnswerDetail, error) {
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

	var answer model.AnswerDetail
	if !cursor.Next(ctx) {
		if err := cursor.Err(); err != nil {
			return nil, err
		}
		return nil, nil
	}

	if err := cursor.Decode(&answer); err != nil {
		return nil, err
	}

	return &answer, nil
}

func (r *MongoAnswerRepository) FindByFilter(ctx context.Context, filter bson.D, sort bson.D, limit int64) ([]*model.AnswerDetail, error) {
	cursor, err := r.collection.Aggregate(ctx, r.pipeline(filter, sort, limit))
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	answers := make([]*model.AnswerDetail, 0)

	if err := cursor.All(ctx, &answers); err != nil {
		return nil, err
	}

	return answers, nil
}

func (r *MongoAnswerRepository) FindMyAnswerDetails(ctx context.Context, filter bson.D, sort bson.D, limit int64) ([]*model.MyAnswerDetail, error) {
	cursor, err := r.collection.Aggregate(ctx, r.myAnswerPipeline(filter, sort, limit))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	answers := make([]*model.MyAnswerDetail, 0)
	if err := cursor.All(ctx, &answers); err != nil {
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

	_, err = r.collection.DeleteMany(ctx, bson.M{"survey_id": objectID})
	return err
}

func (r *MongoAnswerRepository) FindAll(ctx context.Context) ([]*model.AnswerDetail, error) {
	cursor, err := r.collection.Aggregate(ctx, r.pipeline(bson.D{}, bson.D{}, -1))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	answers := make([]*model.AnswerDetail, 0)
	if err := cursor.All(ctx, &answers); err != nil {
		return nil, err
	}
	return answers, nil

}

func (r *MongoAnswerRepository) pipeline(match bson.D, sort bson.D, limit int64) mongo.Pipeline {
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

func (r *MongoAnswerRepository) myAnswerPipeline(match bson.D, sort bson.D, limit int64) mongo.Pipeline {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: match}},
	}

	if len(sort) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: sort}})
	}
	if limit > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$limit", Value: limit}})
	}

	pipeline = append(pipeline,
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: config.CollectionUser},
			{Key: "localField", Value: "author_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "author"},
		}}},
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$author"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: config.CollectionSurvey},
			{Key: "localField", Value: "survey_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "survey"},
		}}},
		bson.D{{Key: "$unwind", Value: "$survey"}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$survey._id"},
			{Key: "title", Value: bson.D{{Key: "$first", Value: "$survey.title"}}},
			{Key: "is_public", Value: bson.D{{Key: "$first", Value: "$survey.is_public"}}},
			{Key: "closed", Value: bson.D{{Key: "$first", Value: "$survey.closed"}}},
			{Key: "author_id", Value: bson.D{{Key: "$first", Value: "$survey.author_id"}}},
			{Key: "expires_at", Value: bson.D{{Key: "$first", Value: "$survey.expires_at"}}},
			{Key: "created_at", Value: bson.D{{Key: "$first", Value: "$survey.created_at"}}},
			{Key: "updated_at", Value: bson.D{{Key: "$first", Value: "$survey.updated_at"}}},
			{Key: "answers", Value: bson.D{{Key: "$push", Value: bson.D{
				{Key: "id", Value: "$_id"},
				{Key: "survey_id", Value: "$survey_id"},
				{Key: "author_id", Value: "$author_id"},
				{Key: "is_anonymous", Value: "$is_anonymous"},
				{Key: "anonymous", Value: "$anonymous"},
				{Key: "content", Value: "$content"},
				{Key: "rate_up", Value: "$rate_up"},
				{Key: "created_at", Value: "$created_at"},
				{Key: "updated_at", Value: "$updated_at"},
				{Key: "author", Value: "$author"},
			}}}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "survey_answers", Value: bson.D{
				{Key: "id", Value: "$_id"},
				{Key: "title", Value: "$title"},
				{Key: "is_public", Value: "$is_public"},
				{Key: "closed", Value: "$closed"},
				{Key: "author_id", Value: "$author_id"},
				{Key: "expires_at", Value: "$expires_at"},
				{Key: "created_at", Value: "$created_at"},
				{Key: "updated_at", Value: "$updated_at"},
				{Key: "answers", Value: "$answers"},
			}},
		}}},
	)

	return pipeline
}
