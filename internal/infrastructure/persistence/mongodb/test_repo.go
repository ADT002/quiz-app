package persistence

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/repository"
)

type TestTemplateMongoRepository struct {
	CollRepo repository.CRUDMongoDB[entity.TestTemplete]
}

func NewTestTemplateMongoRepository() repository.TestTemplateRepository {
	collRepo := NewCollRepository[entity.TestTemplete]("dbapp", "test_templates")
	return &TestTemplateMongoRepository{
		CollRepo: collRepo,
	}
}

func (r *TestTemplateMongoRepository) GetTestTemplates(ctx context.Context, userID string) ([]entity.TestTemplete, error) {
	filter := bson.M{"user_id": userID}
	results, err := r.CollRepo.GetAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get test templates: %w", err)
	}
	return results, nil
}

func (r *TestTemplateMongoRepository) GetTestTemplatesByIDs(ctx context.Context, userID string, testIDs []primitive.ObjectID) ([]entity.TestTemplete, error) {
	filter := bson.M{"_id": bson.M{"$in": testIDs}}
	results, err := r.CollRepo.GetAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get test templates by IDs: %w", err)
	}
	return results, nil
}

func (r *TestTemplateMongoRepository) CreateTestTemplate(ctx context.Context, test entity.TestTemplete) (primitive.ObjectID, error) {
	test.ID = primitive.NewObjectID()

	result, err := r.CollRepo.Create(ctx, test)
	fmt.Println("Inserted ID:", result)
	if err != nil || result.InsertedID == primitive.NilObjectID {
		return primitive.NilObjectID, fmt.Errorf("failed to create test template: %w", err)
	}

	return test.ID, nil
}

func (r *TestTemplateMongoRepository) UpdateTestTemplate(ctx context.Context, test entity.TestTemplete) (entity.TestTemplete, error) {
	filter := bson.M{"user_id": test.UserID, "_id": test.ID}
	results, err := r.CollRepo.Update(ctx, filter, bson.M{"$set": test})
	if err != nil || results.MatchedCount == 0 {
		return entity.TestTemplete{}, fmt.Errorf("failed to update test template: %w", err)
	}
	return test, nil
}

func (r *TestTemplateMongoRepository) DeleteTestTemplate(ctx context.Context, userID string, id primitive.ObjectID) error {
	filter := bson.M{"user_id": userID, "_id": id}
	_, err := r.CollRepo.Delete(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete test template: %w", err)
	}
	return nil
}
