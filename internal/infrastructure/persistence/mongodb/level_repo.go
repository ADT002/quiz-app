package persistence

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/repository"
)

// REPOSITORY
type LevelMongoRepository struct {
	CollRepo repository.CRUDMongoDB[entity.Level]
}

// NewLevelMongoRepository tạo một instance mới của LevelMongoRepository
func NewLevelMongoRepository() repository.LevelRepository {
	collRepo := NewCollRepository[entity.Level]("dbapp", "level")
	return &LevelMongoRepository{
		CollRepo: collRepo,
	}
}

func (r *LevelMongoRepository) CreateLevel(ctx context.Context, t entity.Level) (entity.Level, error) {
	_, err := r.CollRepo.Create(ctx, t)
	if err != nil {
		return t, fmt.Errorf("failed to create Level: %w", err)
	}
	return t, nil
}

func (r *LevelMongoRepository) UpdateLevel(ctx context.Context, t entity.Level) (entity.Level, error) {
	inserted, err := r.CollRepo.Update(ctx, bson.M{"user_id": t.UserID, "_id": t.ID}, bson.M{"$set": t})
	fmt.Println(inserted.MatchedCount, err)
	if err != nil || inserted.MatchedCount == 0 {
		return t, fmt.Errorf("failed to update Level: %w", err)
	}

	return t, nil
}

func (r *LevelMongoRepository) DeleteLevel(ctx context.Context, t entity.Level) error {
	filter := bson.M{"_id": t.ID}
	_, err := r.CollRepo.Delete(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to update Level: %w", err)
	}
	return nil
}

func (r *LevelMongoRepository) GetAllLevels(
	ctx context.Context,
	filter string,
) ([]entity.Level, error) {
	return r.CollRepo.GetAll(ctx, bson.M{"user_id": filter})
}
