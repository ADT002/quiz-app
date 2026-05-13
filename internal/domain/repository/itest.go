package repository

import (
	"context"
	entity "quiz-app/internal/domain/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestTemplateRepository interface {
	CreateTestTemplate(
		ctx context.Context,
		test entity.TestTemplete) (primitive.ObjectID, error)

	GetTestTemplates(
		ctx context.Context,
		userID string) ([]entity.TestTemplete, error)

	GetTestTemplatesByIDs(
		ctx context.Context,
		userID string,
		testIDs []primitive.ObjectID) ([]entity.TestTemplete, error)

	UpdateTestTemplate(
		ctx context.Context,
		test entity.TestTemplete) (entity.TestTemplete, error)

	DeleteTestTemplate(
		ctx context.Context,
		userID string,
		id primitive.ObjectID) error

	// CountUsingQuestion returns how many of ownerID's test templates have
	// questionID in their question_ids array.
	CountUsingQuestion(
		ctx context.Context,
		ownerID string,
		questionID primitive.ObjectID) (int64, error)

	// PullQuestionFromAll removes questionID from question_ids in every
	// test template owned by ownerID. Returns number of docs modified.
	PullQuestionFromAll(
		ctx context.Context,
		ownerID string,
		questionID primitive.ObjectID) (int64, error)
}
