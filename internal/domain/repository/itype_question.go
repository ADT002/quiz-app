package repository

import (
	"context"
	entity "quiz-app/internal/domain/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TypeQuestionRepository interface {
	CreateType(
		ctx context.Context,
		t entity.TypeQuestion) (entity.TypeQuestion, error)

	UpdateType(
		ctx context.Context,
		t entity.TypeQuestion) (entity.TypeQuestion, error)

	DeleteType(
		ctx context.Context,
		t entity.TypeQuestion) error

	GetAllTypes(
		ctx context.Context,
		filter primitive.ObjectID) ([]entity.TypeQuestion, error)
}
