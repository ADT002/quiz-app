package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CRUDMongoDB[T any] interface {
	Create(ctx context.Context, document T) (*mongo.InsertOneResult, error)

	GetAll(ctx context.Context, filter any) ([]T, error)
	GetAllWithOption(ctx context.Context, filter any, opts *options.FindOptions) ([]T, error)

	GetWithProjection(ctx context.Context, filter, projection any) ([]T, error)
	GetOneWithProjection(ctx context.Context, filter, projection any) (bson.M, error)

	GetFilter(ctx context.Context, filter any) (bson.M, error)

	Update(ctx context.Context, filter, update any) (*mongo.UpdateResult, error)
	Delete(ctx context.Context, filter any) (*mongo.DeleteResult, error)

	AggregateLookup(
		ctx context.Context,
		pipeline mongo.Pipeline,
	) ([]bson.M, error)
}
