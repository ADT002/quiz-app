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

type TestOfClassMongoRepository struct {
	CollRepo repository.CRUDMongoDB[entity.TestOfClass]
	coll     *mongo.Collection
}

func NewTestOfClassMongoRepository() repository.TestOfClassRepository {
	collRepo := NewCollRepository[entity.TestOfClass]("dbapp", "test_of_class")
	coll := GetMongoClient().Database("dbapp").Collection("test_of_class")
	return &TestOfClassMongoRepository{
		CollRepo: collRepo,
		coll:     coll,
	}
}

func (r *TestOfClassMongoRepository) CreateTestOfClass(ctx context.Context, test entity.TestOfClass) (primitive.ObjectID, error) {
	test.ID = primitive.NewObjectID()

	result, err := r.CollRepo.Create(ctx, test)
	fmt.Println("Inserted TestOfClass ID:", result)
	if err != nil || result.InsertedID == primitive.NilObjectID {
		return primitive.NilObjectID, fmt.Errorf("failed to create test of class: %w", err)
	}

	return test.ID, nil
}

func (r *TestOfClassMongoRepository) GetTestsOfClass(ctx context.Context, classID primitive.ObjectID) ([]entity.TestOfClass, error) {
	filter := bson.M{"class_id": classID}
	results, err := r.CollRepo.GetAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get tests of class: %w", err)
	}
	return results, nil
}

func (r *TestOfClassMongoRepository) GetTestOfClassByID(ctx context.Context, testID primitive.ObjectID) (entity.TestOfClass, error) {
	filter := bson.M{"_id": testID}
	raw, err := r.CollRepo.GetFilter(ctx, filter)
	if err != nil {
		return entity.TestOfClass{}, fmt.Errorf("failed to get test of class: %w", err)
	}
	if raw == nil {
		return entity.TestOfClass{}, fmt.Errorf("test of class not found")
	}

	var test entity.TestOfClass
	bsonBytes, err := bson.Marshal(raw)
	if err != nil {
		return entity.TestOfClass{}, fmt.Errorf("failed to marshal test of class: %w", err)
	}
	if err := bson.Unmarshal(bsonBytes, &test); err != nil {
		return entity.TestOfClass{}, fmt.Errorf("failed to unmarshal test of class: %w", err)
	}
	return test, nil
}

func (r *TestOfClassMongoRepository) UpdateTestOfClass(ctx context.Context, test entity.TestOfClass) (entity.TestOfClass, error) {
	filter := bson.M{"user_id": test.UserID, "_id": test.ID}
	result, err := r.CollRepo.Update(ctx, filter, bson.M{"$set": test})
	if err != nil || result.MatchedCount == 0 {
		return entity.TestOfClass{}, fmt.Errorf("failed to update test of class: %w", err)
	}
	return test, nil
}

// CountUsingQuestion counts owner's test_of_class docs whose question_ids
// contain questionID (stored as hex string per BaseTest.QuestionIDs).
func (r *TestOfClassMongoRepository) CountUsingQuestion(
	ctx context.Context,
	ownerID string,
	questionID primitive.ObjectID,
) (int64, error) {
	return r.coll.CountDocuments(ctx, bson.M{
		"user_id":      ownerID,
		"question_ids": questionID.Hex(),
	})
}

func (r *TestOfClassMongoRepository) PullQuestionFromAll(
	ctx context.Context,
	ownerID string,
	questionID primitive.ObjectID,
) (int64, error) {
	res, err := r.coll.UpdateMany(
		ctx,
		bson.M{"user_id": ownerID, "question_ids": questionID.Hex()},
		bson.M{"$pull": bson.M{"question_ids": questionID.Hex()}},
	)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

func (r *TestOfClassMongoRepository) DeleteTestOfClass(ctx context.Context, userID string, id primitive.ObjectID) error {
	filter := bson.M{"user_id": userID, "_id": id}
	_, err := r.CollRepo.Delete(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete test of class: %w", err)
	}
	return nil
}
