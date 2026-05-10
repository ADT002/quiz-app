package persistence

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/repository"
)

type TopicMongoRepository struct {
	CollRepo repository.CRUDMongoDB[entity.Topic]
	coll     *mongo.Collection
}

func NewTopicMongoRepository() repository.TopicRepository {
	collRepo := NewCollRepository[entity.Topic]("dbapp", "topic")
	coll := GetMongoClient().Database("dbapp").Collection("topic")
	return &TopicMongoRepository{CollRepo: collRepo, coll: coll}
}

func (r *TopicMongoRepository) CreateTopic(ctx context.Context, t entity.Topic) (entity.Topic, error) {
	_, err := r.CollRepo.Create(ctx, t)
	if err != nil {
		return t, fmt.Errorf("failed to create Topic: %w", err)
	}
	return t, nil
}

func (r *TopicMongoRepository) UpdateTopic(ctx context.Context, t entity.Topic) (entity.Topic, error) {
	res, err := r.CollRepo.Update(ctx, bson.M{"user_id": t.UserID, "_id": t.ID}, bson.M{"$set": t})
	if err != nil {
		return t, fmt.Errorf("failed to update Topic: %w", err)
	}
	if res.MatchedCount == 0 {
		return t, fmt.Errorf("topic not found")
	}
	return t, nil
}

func (r *TopicMongoRepository) DeleteTopic(ctx context.Context, t entity.Topic) error {
	filter := bson.M{"_id": t.ID, "user_id": t.UserID}
	_, err := r.CollRepo.Delete(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete Topic: %w", err)
	}
	return nil
}

// GetAllTopics returns topics for the given user_id sorted by topic_no asc.
// Stable ordering is required by the Question Bank UI per F4.
func (r *TopicMongoRepository) GetAllTopics(
	ctx context.Context,
	userID string,
) ([]entity.Topic, error) {
	opts := options.Find().SetSort(bson.D{
		{Key: "topic_no", Value: 1},
		{Key: "_id", Value: 1},
	})
	return r.CollRepo.GetAllWithOption(ctx, bson.M{"user_id": userID}, opts)
}

// ReorderTopics rewrites topic_no in bulk via a single BulkWrite op.
// Topics not in orderedIDs are left untouched.
func (r *TopicMongoRepository) ReorderTopics(
	ctx context.Context,
	userID string,
	orderedIDs []string,
) error {
	if len(orderedIDs) == 0 {
		return nil
	}
	models := make([]mongo.WriteModel, 0, len(orderedIDs))
	for i, hex := range orderedIDs {
		id, err := primitive.ObjectIDFromHex(hex)
		if err != nil {
			return fmt.Errorf("invalid topic id %q: %w", hex, err)
		}
		models = append(models,
			mongo.NewUpdateOneModel().
				SetFilter(bson.M{"_id": id, "user_id": userID}).
				SetUpdate(bson.M{"$set": bson.M{"topic_no": i + 1}}),
		)
	}
	_, err := r.coll.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
	if err != nil {
		return fmt.Errorf("reorder topics failed: %w", err)
	}
	return nil
}
