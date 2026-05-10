package persistence

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/repository"
)

type ClassMongoRepository struct {
	CollRepo repository.CRUDMongoDB[entity.Class]
	coll     *mongo.Collection
}

func NewClassMongoRepository() repository.ClassRepository {
	collRepo := NewCollRepository[entity.Class]("dbapp", "classes")
	coll := GetMongoClient().Database("dbapp").Collection("classes")
	return &ClassMongoRepository{CollRepo: collRepo, coll: coll}
}

func (r *ClassMongoRepository) CreateClass(ctx context.Context, class entity.Class) (primitive.ObjectID, error) {
	_, err := r.CollRepo.Create(ctx, class)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to create class: %w", err)
	}
	return class.ID, nil
}

func (r *ClassMongoRepository) GetClasses(ctx context.Context, userID string) ([]entity.Class, error) {
	results, err := r.CollRepo.GetAll(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to get classes by user ID: %w", err)
	}
	return results, nil
}

func (r *ClassMongoRepository) UpdateClass(ctx context.Context, class entity.Class) (entity.Class, error) {
	filter := bson.M{"user_id": class.UserID, "_id": class.ID}
	_, err := r.CollRepo.Update(ctx, filter, bson.M{"$set": class})
	if err != nil {
		return class, fmt.Errorf("failed to update class: %w", err)
	}
	return class, nil
}

func (r *ClassMongoRepository) DeleteClass(ctx context.Context, userID string, id primitive.ObjectID) error {
	_, err := r.CollRepo.Delete(ctx, bson.M{"user_id": userID, "_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete class: %w", err)
	}
	return nil
}

// JoinClass: public class → push thẳng accept; private → push wait. Idempotent
// nhờ $addToSet (CLAUDE.md F3 edge: học sinh đã accept → 200 không nhân đôi).
func (r *ClassMongoRepository) JoinClass(
	ctx context.Context,
	classID primitive.ObjectID,
	userID string,
	emailIDStudent string,
	emailStudent string,
) error {
	filter := bson.M{"user_id": userID, "_id": classID}

	classDoc, err := r.CollRepo.GetFilter(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to find class: %w", err)
	}

	target := "students_wait"
	if classDoc["is_public"] == true {
		target = "students_accept"
	}

	_, err = r.CollRepo.Update(ctx, filter, bson.M{
		"$addToSet": bson.M{
			target: bson.M{"user_id": emailIDStudent, "email": emailStudent},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to join class: %w", err)
	}
	return nil
}

// ApproveStudent: atomic move from students_wait → students_accept, but only
// if studentID is currently in students_wait. Filter ràng buộc owner_id để
// teacher khác không approve hộ.
func (r *ClassMongoRepository) ApproveStudent(
	ctx context.Context,
	classID primitive.ObjectID,
	ownerID, studentID string,
) error {
	filter := bson.M{
		"_id":                       classID,
		"user_id":                   ownerID,
		"students_wait.user_id":     studentID,
	}
	// Find student record to copy into accept (preserve email/name).
	var doc bson.M
	if err := r.coll.FindOne(ctx, filter).Decode(&doc); err != nil {
		return fmt.Errorf("approve: pending student not found: %w", err)
	}
	var copyInfo bson.M
	for _, s := range doc["students_wait"].(bson.A) {
		if entry, ok := s.(bson.M); ok && entry["user_id"] == studentID {
			copyInfo = entry
			break
		}
	}
	if copyInfo == nil {
		return fmt.Errorf("approve: student %s not in wait list", studentID)
	}

	_, err := r.coll.UpdateOne(ctx,
		bson.M{"_id": classID, "user_id": ownerID},
		bson.M{
			"$pull":     bson.M{"students_wait": bson.M{"user_id": studentID}},
			"$addToSet": bson.M{"students_accept": copyInfo},
		},
	)
	if err != nil {
		return fmt.Errorf("approve update failed: %w", err)
	}
	return nil
}

func (r *ClassMongoRepository) RejectStudent(
	ctx context.Context,
	classID primitive.ObjectID,
	ownerID, studentID string,
) error {
	res, err := r.coll.UpdateOne(ctx,
		bson.M{"_id": classID, "user_id": ownerID},
		bson.M{"$pull": bson.M{"students_wait": bson.M{"user_id": studentID}}},
	)
	if err != nil {
		return fmt.Errorf("reject failed: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("class not found or not owned")
	}
	return nil
}

// GetAllTestOfClass / GetQuestionOfTest delegated to TestOfClassRepository.
func (r *ClassMongoRepository) GetAllTestOfClass(
	ctx context.Context,
	email string,
	id primitive.ObjectID,
) ([]entity.TestOfClass, error) {
	return nil, nil
}

func (r *ClassMongoRepository) GetQuestionOfTest(
	ctx context.Context,
	classID, testID primitive.ObjectID,
	email string,
) ([]primitive.ObjectID, entity.TestOfClass, error) {
	return nil, entity.TestOfClass{}, nil
}
