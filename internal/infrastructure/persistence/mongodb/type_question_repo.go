package persistence

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/repository"
)

// REPOSITORY
type TypeQuestionMongoRepository struct {
	CollRepo repository.CRUDMongoDB[entity.TypeQuestion]
}

// NewTypeQuestionMongoRepository tạo một instance mới của TypeQuestionMongoRepository
func NewTypeQuestionMongoRepository() repository.TypeQuestionRepository {
	collRepo := NewCollRepository[entity.TypeQuestion]("dbapp", "type_question")
	return &TypeQuestionMongoRepository{
		CollRepo: collRepo,
	}
}

func (r *TypeQuestionMongoRepository) CreateType(ctx context.Context, t entity.TypeQuestion) (entity.TypeQuestion, error) {
	_, err := r.CollRepo.Create(ctx, t)
	if err != nil {
		return t, fmt.Errorf("failed to create TypeQuestion: %w", err)
	}
	return t, nil
}

// GetTypeQuestionByAuthorEmail fetches all TypeQuestiones for a given author's email ID.
func (r *TypeQuestionMongoRepository) UpdateType(ctx context.Context, t entity.TypeQuestion) (entity.TypeQuestion, error) {
	_, err := r.CollRepo.Create(ctx, t)
	if err != nil {
		return t, fmt.Errorf("failed to create TypeQuestion: %w", err)
	}
	return t, nil
}

// UpdateTypeQuestion implements repository.TypeQuestionRepository.UpdateTypeQuestion
func (r *TypeQuestionMongoRepository) DeleteType(ctx context.Context, t entity.TypeQuestion) error {
	filter := bson.M{"_id": t.ID}
	_, err := r.CollRepo.Update(ctx, filter, bson.M{"$set": t})
	if err != nil {
		return fmt.Errorf("failed to update TypeQuestion: %w", err)
	}
	return nil
}

// GetAllTypeQuestion implements repository.TypeQuestionRepository.GetAllTypes
func (r *TypeQuestionMongoRepository) GetAllTypes(
	ctx context.Context,
	filter primitive.ObjectID,
) ([]entity.TypeQuestion, error) {
	return r.CollRepo.GetAll(ctx, filter)
}
