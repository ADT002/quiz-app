package repository

import (
	"context"
	entity "quiz-app/internal/domain/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestOfClassRepository interface {
	CreateTestOfClass(
		ctx context.Context,
		test entity.TestOfClass) (primitive.ObjectID, error)

	GetTestsOfClass(
		ctx context.Context,
		classID primitive.ObjectID) ([]entity.TestOfClass, error)

	GetTestOfClassByID(
		ctx context.Context,
		testID primitive.ObjectID) (entity.TestOfClass, error)

	UpdateTestOfClass(
		ctx context.Context,
		test entity.TestOfClass) (entity.TestOfClass, error)

	DeleteTestOfClass(
		ctx context.Context,
		userID string,
		id primitive.ObjectID) error
}
