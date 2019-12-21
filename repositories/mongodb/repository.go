package mongodb

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/vrnvu/tiny-to-go/tiny"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewMongoRepository(mongoURL, mongoDB string, mongoTimeout int) (tiny.Repository, error) {
	repo := &mongoRepository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}
	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewMongoRepo")
	}
	repo.client = client
	return repo, nil
}

func (r *mongoRepository) Find(code string) (*tiny.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	redirect := &tiny.Redirect{}
	collection := r.client.Database(r.database).Collection("redirects")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&redirect)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(tiny.ErrRedirectNotFound, "repository.Redirect.Find")
		}
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	return redirect, nil
}

func (r *mongoRepository) Save(redirect *tiny.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	collection := r.client.Database(r.database).Collection("redirects")
	_, err := collection.InsertOne(
		ctx,
		bson.M{
			"code":      redirect.Code,
			"url":       redirect.URL,
			"timestamp": redirect.Timestamp,
		},
	)
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}
