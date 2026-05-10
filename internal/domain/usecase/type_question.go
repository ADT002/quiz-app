package usecase

import (
	"context"
	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TypeQuestionUseCase struct {
	TypeQuestionRepo repository.TypeQuestionRepository
}

func NewTypeQuestionUseCase(tr repository.TypeQuestionRepository) *TypeQuestionUseCase {
	return &TypeQuestionUseCase{
		TypeQuestionRepo: tr,
	}
}

// CreateTypeQuestion creates a new TypeQuestion
func (uc *TypeQuestionUseCase) UpdateType(ctx context.Context, t entity.TypeQuestion) (entity.TypeQuestion, error) {
	newTypeQuestion, err := uc.TypeQuestionRepo.UpdateType(ctx, t)
	if err != nil {
		return t, err
	}
	return newTypeQuestion, nil
}

// GetTypeQuestionByID retrieves a TypeQuestion by its ID
func (uc *TypeQuestionUseCase) CreateType(ctx context.Context, t entity.TypeQuestion) (entity.TypeQuestion, error) {
	return uc.TypeQuestionRepo.CreateType(ctx, t)
}

// GetTypeQuestionByID retrieves a TypeQuestion by its ID
func (uc *TypeQuestionUseCase) DeleteType(ctx context.Context, t entity.TypeQuestion) error {
	return uc.TypeQuestionRepo.DeleteType(ctx, t)
}

// UpdateTypeQuestion updates an existing TypeQuestion
func (uc *TypeQuestionUseCase) GetAllTypes(
	ctx context.Context,
	filter primitive.ObjectID,
) ([]entity.TypeQuestion, error) {
	return uc.TypeQuestionRepo.GetAllTypes(ctx, filter)
}
