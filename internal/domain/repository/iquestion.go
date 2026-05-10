package repository

import (
	"context"
	entity "quiz-app/internal/domain/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// QuestionFilters captures all list query parameters for questions.
// Empty fields = no filter on that axis.
type QuestionFilters struct {
	OwnerID string
	TopicID string
	LevelID string
	Type    string
	Search  string
	Page    int64
	Limit   int64
}

type QuestionRepository interface {
	CreateQuestion(
		ctx context.Context,
		question entity.Question) (primitive.ObjectID, error)

	// ListQuestions returns paginated, filterable questions owned by f.OwnerID.
	// Returns (items, total, error). E#10 list shape requirement.
	ListQuestions(
		ctx context.Context,
		f QuestionFilters) ([]entity.Question, int64, error)

	UpdateQuestion(
		ctx context.Context,
		question entity.Question) (entity.Question, error)

	DeleteQuestion(
		ctx context.Context,
		question entity.Question) error

	GetAllQuestions(
		ctx context.Context,
		question_ids []primitive.ObjectID) ([]entity.Question, error)
}
