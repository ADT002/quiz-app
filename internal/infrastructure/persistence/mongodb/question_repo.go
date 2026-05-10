package persistence

import (
	"context"
	"fmt"
	"regexp"
	"time"

	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/repository"
	"quiz-app/internal/domain/usecase"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const questionsDB = "dbapp"
const questionsColl = "questions"

type QuestionMongoRepository struct {
	CollRepo repository.CRUDMongoDB[entity.Question]
	coll     *mongo.Collection
}

func NewQuestionMongoRepository() repository.QuestionRepository {
	collRepo := NewCollRepository[entity.Question](questionsDB, questionsColl)
	coll := GetMongoClient().Database(questionsDB).Collection(questionsColl)
	return &QuestionMongoRepository{CollRepo: collRepo, coll: coll}
}

func (r *QuestionMongoRepository) GetAllQuestions(ctx context.Context, ids []primitive.ObjectID) ([]entity.Question, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("no question IDs provided")
	}
	filter := bson.M{"_id": bson.M{"$in": ids}}
	projection := bson.M{
		"createdAt": 0,
		"updatedAt": 0,
		"metadata":  0,
	}
	return r.CollRepo.GetWithProjection(ctx, filter, projection)
}

func (r *QuestionMongoRepository) CreateQuestion(ctx context.Context, q entity.Question) (primitive.ObjectID, error) {
	result, err := r.CollRepo.Create(ctx, q)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to create question: %w", err)
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// ListQuestions implements paginated + filterable owner-scoped query per E#10.
func (r *QuestionMongoRepository) ListQuestions(
	ctx context.Context,
	f repository.QuestionFilters,
) ([]entity.Question, int64, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 20
	}

	filter := bson.M{"metadata.user_id": f.OwnerID}
	if f.TopicID != "" {
		if oid, err := primitive.ObjectIDFromHex(f.TopicID); err == nil {
			filter["topic"] = oid
		}
	}
	if f.LevelID != "" {
		if oid, err := primitive.ObjectIDFromHex(f.LevelID); err == nil {
			filter["level"] = oid
		}
	}
	if f.Type != "" {
		filter["type"] = f.Type
	}
	if f.Search != "" {
		// Case-insensitive substring match against question content text.
		// Quoted to avoid regex injection from arbitrary user input.
		filter["question_content.content.text"] = bson.M{
			"$regex":   regexp.QuoteMeta(f.Search),
			"$options": "i",
		}
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetLimit(f.Limit).
		SetSkip((f.Page - 1) * f.Limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	items, err := r.CollRepo.GetAllWithOption(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// sanitizeByType keeps only the fields relevant to q.Type and emits $unset for
// the rest, so changing question type wipes stale fields. Names must match
// CLAUDE.md D2 — see the constants in usecase/question.go.
func sanitizeByType(q *entity.Question) (unset bson.M) {
	unset = bson.M{}
	switch q.Type {
	case usecase.SingleQuestion, usecase.MultiQuestion:
		q.FillInTheBlanks = nil
		q.OrderItems = nil
		q.MatchItems = nil
		q.MatchOptions = nil
		unset["fill_in_the_blanks"] = ""
		unset["order_items"] = ""
		unset["match_items"] = ""
		unset["match_options"] = ""

	case usecase.FillInTheBlankQuestion:
		q.Options = nil
		q.OrderItems = nil
		q.MatchItems = nil
		q.MatchOptions = nil
		unset["options"] = ""
		unset["order_items"] = ""
		unset["match_items"] = ""
		unset["match_options"] = ""

	case usecase.OrderQuestion:
		q.Options = nil
		q.FillInTheBlanks = nil
		q.MatchItems = nil
		q.MatchOptions = nil
		unset["options"] = ""
		unset["fill_in_the_blanks"] = ""
		unset["match_items"] = ""
		unset["match_options"] = ""

	case usecase.MatchQuestion:
		q.Options = nil
		q.FillInTheBlanks = nil
		q.OrderItems = nil
		unset["options"] = ""
		unset["fill_in_the_blanks"] = ""
		unset["order_items"] = ""
	}
	return
}

func (r *QuestionMongoRepository) UpdateQuestion(ctx context.Context, q entity.Question) (entity.Question, error) {
	filter := bson.M{
		"metadata.user_id": q.Metadata.User_ID,
		"_id":              q.ID,
	}
	unset := sanitizeByType(&q)
	q.Updated_At = time.Now()

	update := bson.M{"$set": q, "$unset": unset}
	result, err := r.CollRepo.Update(ctx, filter, update)
	if err != nil {
		return entity.Question{}, fmt.Errorf("update failed: %w", err)
	}
	if result.MatchedCount == 0 {
		return entity.Question{}, fmt.Errorf("question not found")
	}
	return q, nil
}

func (r *QuestionMongoRepository) DeleteQuestion(ctx context.Context, q entity.Question) error {
	filter := bson.M{"metadata.user_id": q.Metadata.User_ID, "_id": q.ID}
	_, err := r.CollRepo.Delete(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete question: %w", err)
	}
	return nil
}
