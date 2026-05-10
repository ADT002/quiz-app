package persistence

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/repository"
)

type TestOfClassMongoRepository struct {
	CollRepo repository.CRUDMongoDB[entity.TestOfClass]
}

func NewTestOfClassMongoRepository() repository.TestOfClassRepository {
	collRepo := NewCollRepository[entity.TestOfClass]("dbapp", "test_of_class")
	return &TestOfClassMongoRepository{
		CollRepo: collRepo,
	}
}

func (r *TestOfClassMongoRepository) CreateTestOfClass(ctx context.Context, test entity.TestOfClass) (primitive.ObjectID, error) {
	test.ID = primitive.NewObjectID()

	result, err := r.CollRepo.Create(ctx, test)
	fmt.Println("Inserted TestOfClass ID:", result)
	if err != nil || result.InsertedID == primitive.NilObjectID {
		return primitive.NilObjectID, fmt.Errorf("failed to create test of class: %w", err)
	}

	return test.ID, nil
}

func (r *TestOfClassMongoRepository) GetTestsOfClass(ctx context.Context, classID primitive.ObjectID) ([]entity.TestOfClass, error) {
	filter := bson.M{"class_id": classID}
	results, err := r.CollRepo.GetAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get tests of class: %w", err)
	}
	return results, nil
}

func (r *TestOfClassMongoRepository) GetTestOfClassByID(ctx context.Context, testID primitive.ObjectID) (entity.TestOfClass, error) {
	filter := bson.M{"_id": testID}
	raw, err := r.CollRepo.GetFilter(ctx, filter)
	if err != nil {
		return entity.TestOfClass{}, fmt.Errorf("failed to get test of class: %w", err)
	}
	if raw == nil {
		return entity.TestOfClass{}, fmt.Errorf("test of class not found")
	}

	var test entity.TestOfClass
	bsonBytes, err := bson.Marshal(raw)
	if err != nil {
		return entity.TestOfClass{}, fmt.Errorf("failed to marshal test of class: %w", err)
	}
	if err := bson.Unmarshal(bsonBytes, &test); err != nil {
		return entity.TestOfClass{}, fmt.Errorf("failed to unmarshal test of class: %w", err)
	}
	return test, nil
}

func (r *TestOfClassMongoRepository) UpdateTestOfClass(ctx context.Context, test entity.TestOfClass) (entity.TestOfClass, error) {
	filter := bson.M{"user_id": test.UserID, "_id": test.ID}
	result, err := r.CollRepo.Update(ctx, filter, bson.M{"$set": test})
	if err != nil || result.MatchedCount == 0 {
		return entity.TestOfClass{}, fmt.Errorf("failed to update test of class: %w", err)
	}
	return test, nil
}

func (r *TestOfClassMongoRepository) DeleteTestOfClass(ctx context.Context, userID string, id primitive.ObjectID) error {
	filter := bson.M{"user_id": userID, "_id": id}
	_, err := r.CollRepo.Delete(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete test of class: %w", err)
	}
	return nil
}
