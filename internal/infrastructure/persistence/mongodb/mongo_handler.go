package persistence

import (
	"context"
	"fmt"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CollRepository[T any] struct {
	Collection *mongo.Collection
}

func NewCollRepository[T any](dbName, collName string) repository.CRUDMongoDB[T] {
	client := GetMongoClient()
	db := client.Database(dbName)
	return &CollRepository[T]{
		Collection: db.Collection(collName),
	}

}

func (r *CollRepository[T]) Create(
	ctx context.Context,
	document T,
) (*mongo.InsertOneResult, error) {

	result, err := r.Collection.InsertOne(ctx, document)
	if err != nil {
		return result, fmt.Errorf("failed to create document: %w", err)
	}

	return result, nil
}
func (r *CollRepository[T]) GetAll(
	ctx context.Context,
	filter any,
) ([]T, error) {

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// GetAllWithOption retrieves documents based on the provided filter and options (e.g., limit, skip)
func (r *CollRepository[T]) GetAllWithOption(ctx context.Context, filter any, opts *options.FindOptions) ([]T, error) {
	// Create an empty slice to store the results
	var results []T

	// Execute the query with the filter and options
	cursor, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor to decode each document and add to results slice
	for cursor.Next(ctx) {
		var doc T
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode document: %w", err)
		}
		results = append(results, doc)
	}

	// Check if there was an error during iteration
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor iteration error: %w", err)
	}

	return results, nil
}
func (r *CollRepository[T]) GetWithProjection(ctx context.Context, filter, projection any) ([]T, error) {
	var results []T

	cursor, err := r.Collection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, fmt.Errorf("failed to execute find query: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item T
		if err := cursor.Decode(&item); err != nil {
			return nil, fmt.Errorf("failed to decode document: %w", err)
		}
		results = append(results, item)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor encountered an error: %w", err)
	}

	return results, nil
}

func (r *CollRepository[T]) GetOneWithProjection(ctx context.Context, filter, projection any) (bson.M, error) {
	var result bson.M
	err := r.Collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // không có tài liệu phù hợp
		}
		return nil, fmt.Errorf("failed to execute findOne query: %w", err)
	}
	return result, nil
}

func (r *CollRepository[T]) GetFilter(ctx context.Context, filter any) (bson.M, error) {
	var result map[string]interface{}
	err := r.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No documents found
		}
		return nil, fmt.Errorf("failed to find document: %w", err)
	}
	return result, nil
}

func (r *CollRepository[T]) Update(ctx context.Context, filter, update any) (*mongo.UpdateResult, error) {
	result, err := r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}
	return result, nil
}

func (r *CollRepository[T]) UpdateMany(ctx context.Context, filter, update any) (*mongo.UpdateResult, error) {
	result, err := r.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update documents: %w", err)
	}
	return result, nil
}

func (r *CollRepository[T]) Delete(ctx context.Context, filter any) (*mongo.DeleteResult, error) {
	result, err := r.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete document: %w", err)
	}
	return result, nil
}

func (r *CollRepository[T]) AggregateLookup(
	ctx context.Context,
	pipeline mongo.Pipeline,
) ([]bson.M, error) {

	cursor, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to execute aggregate: %w", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode aggregate results: %w", err)
	}

	return results, nil
}
