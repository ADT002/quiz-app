package usecase

import (
	"context"
	"fmt"
	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestOfClassUseCase struct {
	TestOfClassRepo repository.TestOfClassRepository
}

func NewTestOfClassUseCase(repo repository.TestOfClassRepository) *TestOfClassUseCase {
	return &TestOfClassUseCase{
		TestOfClassRepo: repo,
	}
}

func (uc *TestOfClassUseCase) CreateTestOfClass(ctx context.Context, test entity.TestOfClass) (primitive.ObjectID, error) {
	return uc.TestOfClassRepo.CreateTestOfClass(ctx, test)
}

func (uc *TestOfClassUseCase) GetTestsOfClass(ctx context.Context, classID primitive.ObjectID) ([]entity.TestOfClass, error) {
	return uc.TestOfClassRepo.GetTestsOfClass(ctx, classID)
}

func (uc *TestOfClassUseCase) GetTestOfClassByID(ctx context.Context, testID string) (entity.TestOfClass, error) {
	objID, err := primitive.ObjectIDFromHex(testID)
	if err != nil {
		fmt.Println("test_of_class.go/usecase: invalid testID:", testID, err)
		return entity.TestOfClass{}, err
	}
	return uc.TestOfClassRepo.GetTestOfClassByID(ctx, objID)
}

func (uc *TestOfClassUseCase) UpdateTestOfClass(ctx context.Context, test entity.TestOfClass) (entity.TestOfClass, error) {
	test.UpdatedAt = time.Now()
	return uc.TestOfClassRepo.UpdateTestOfClass(ctx, test)
}

func (uc *TestOfClassUseCase) DeleteTestOfClass(ctx context.Context, userID string, testID primitive.ObjectID) error {
	return uc.TestOfClassRepo.DeleteTestOfClass(ctx, userID, testID)
}
