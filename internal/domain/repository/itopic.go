package repository

import (
	"context"
	entity "quiz-app/internal/domain/entity"
)

type TopicRepository interface {
	CreateTopic(
		ctx context.Context,
		topic entity.Topic) (entity.Topic, error)

	UpdateTopic(
		ctx context.Context,
		topic entity.Topic) (entity.Topic, error)

	DeleteTopic(
		ctx context.Context,
		topic entity.Topic) error

	GetAllTopics(
		ctx context.Context,
		filter string) ([]entity.Topic, error)

	// ReorderTopics atomically rewrites topic_no for the given user, in the
	// order of the supplied IDs (1-based). Topics not in the list keep their
	// existing topic_no. Per CLAUDE.md F4 — order_question.topic_no must be
	// unique within an owner so list returns deterministic ordering.
	ReorderTopics(
		ctx context.Context,
		userID string,
		orderedIDs []string) error
}
