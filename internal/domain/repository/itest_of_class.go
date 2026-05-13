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

	// CountUsingQuestion returns how many test_of_class docs (owned via class
	// owner or any class) have questionID in question_ids. We do NOT filter
	// by user_id because a question may also be embedded in tests created by
	// the owner inside a class — owner_id is on the TestOfClass row.
	CountUsingQuestion(
		ctx context.Context,
		ownerID string,
		questionID primitive.ObjectID) (int64, error)

	// PullQuestionFromAll removes questionID from question_ids in every
	// test_of_class doc owned by ownerID.
	PullQuestionFromAll(
		ctx context.Context,
		ownerID string,
		questionID primitive.ObjectID) (int64, error)
}
