package usecase

import (
	"context"
	"errors"
	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ErrTopicInUse blocks deleting a topic still referenced by any question.
var ErrTopicInUse = errors.New("topic is referenced by a question")

type TopicUseCase struct {
	TopicRepo    repository.TopicRepository
	QuestionRepo repository.QuestionRepository
}

func NewTopicUseCase(tr repository.TopicRepository, qr repository.QuestionRepository) *TopicUseCase {
	return &TopicUseCase{
		TopicRepo:    tr,
		QuestionRepo: qr,
	}
}

// CreateTopic creates a new Topic
func (uc *TopicUseCase) CreateTopic(ctx context.Context, topic entity.Topic) (any, error) {
	topic.ID = primitive.NewObjectID()
	newTopic, err := uc.TopicRepo.CreateTopic(ctx, topic)
	if err != nil {
		return nil, err
	}
	return newTopic, nil
}

// GetTopicByID retrieves a Topic by its ID
func (uc *TopicUseCase) UpdateTopic(ctx context.Context, t entity.Topic) (entity.Topic, error) {
	return uc.TopicRepo.UpdateTopic(ctx, t)
}

// DeleteTopic blocks delete if any question still references the topic.
func (uc *TopicUseCase) DeleteTopic(ctx context.Context, t entity.Topic) error {
	n, err := uc.QuestionRepo.CountByTopic(ctx, t.UserID, t.ID)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrTopicInUse
	}
	return uc.TopicRepo.DeleteTopic(ctx, t)
}

// UpdateTopic updates an existing Topic
func (uc *TopicUseCase) GetAllTopics(ctx context.Context, filter string) ([]entity.Topic, error) {
	return uc.TopicRepo.GetAllTopics(ctx, filter)
}

// ReorderTopics atomically rewrites topic_no for the given user.
func (uc *TopicUseCase) ReorderTopics(ctx context.Context, userID string, orderedIDs []string) error {
	return uc.TopicRepo.ReorderTopics(ctx, userID, orderedIDs)
}
