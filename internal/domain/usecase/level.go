package usecase

import (
	"context"
	"errors"
	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ErrLevelInUse blocks deleting a level still referenced by any question.
var ErrLevelInUse = errors.New("level is referenced by a question")

type LevelUseCase struct {
	LevelRepo    repository.LevelRepository
	QuestionRepo repository.QuestionRepository
}

func NewLevelUseCase(lr repository.LevelRepository, qr repository.QuestionRepository) *LevelUseCase {
	return &LevelUseCase{
		LevelRepo:    lr,
		QuestionRepo: qr,
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

// DeleteLevel blocks delete if any question still references the level.
func (uc *LevelUseCase) DeleteLevel(ctx context.Context, level entity.Level) error {
	n, err := uc.QuestionRepo.CountByLevel(ctx, level.UserID, level.ID)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrLevelInUse
	}
	return uc.LevelRepo.DeleteLevel(ctx, level)
}

// UpdateLevel updates an existing Level
func (uc *LevelUseCase) GetAllLevels(ctx context.Context, filter string) ([]entity.Level, error) {
	return uc.LevelRepo.GetAllLevels(ctx, filter)
}
