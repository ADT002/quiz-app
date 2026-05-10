package repository

import (
	"context"
	entity "quiz-app/internal/domain/entity"
)

type LevelRepository interface {
	CreateLevel(
		ctx context.Context,
		Level entity.Level) (entity.Level, error)

	UpdateLevel(
		ctx context.Context,
		level entity.Level) (entity.Level, error)

	DeleteLevel(
		ctx context.Context,
		level entity.Level) error

	GetAllLevels(
		ctx context.Context,
		filter string) ([]entity.Level, error)
}
