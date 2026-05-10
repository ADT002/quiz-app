package usecase

import (
	"context"
	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LevelUseCase struct {
	LevelRepo repository.LevelRepository
}

func NewLevelUseCase(tr repository.LevelRepository) *LevelUseCase {
	return &LevelUseCase{
		LevelRepo: tr,
	}
}

// CreateLevel creates a new Level
func (uc *LevelUseCase) CreateLevel(ctx context.Context, level entity.Level) (primitive.ObjectID, error) {
	level.ID = primitive.NewObjectID()

	newLevel, err := uc.LevelRepo.CreateLevel(ctx, level)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return newLevel.ID, nil
}

// GetLevelByID retrieves a Level by its ID
func (uc *LevelUseCase) UpdateLevel(ctx context.Context, level entity.Level) (entity.Level, error) {
	return uc.LevelRepo.UpdateLevel(ctx, level)
}

// GetLevelByID retrieves a Level by its ID
func (uc *LevelUseCase) DeleteLevel(ctx context.Context, level entity.Level) error {
	return uc.LevelRepo.DeleteLevel(ctx, level)
}

// UpdateLevel updates an existing Level
func (uc *LevelUseCase) GetAllLevels(ctx context.Context, filter string) ([]entity.Level, error) {
	return uc.LevelRepo.GetAllLevels(ctx, filter)
}
