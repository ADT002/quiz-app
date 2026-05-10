package usecase

import (
	"context"
	"fmt"
	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestTemplateUseCase struct {
	TestRepo repository.TestTemplateRepository
}

func NewTestTemplateUseCase(tr repository.TestTemplateRepository) *TestTemplateUseCase {
	return &TestTemplateUseCase{
		TestRepo: tr,
	}
}

func (uc *TestTemplateUseCase) CreateTestTemplate(ctx context.Context, test entity.TestTemplete) (primitive.ObjectID, error) {
	return uc.TestRepo.CreateTestTemplate(ctx, test)
}

func (uc *TestTemplateUseCase) GetTestTemplatesByIDs(ctx context.Context, userID string, testIDs []string) ([]entity.TestTemplete, error) {
	objectIDs := make([]primitive.ObjectID, 0, len(testIDs))

	for _, id := range testIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			fmt.Println("test.go/usecase: invalid testID:", id, err)
			return nil, err
		}
		objectIDs = append(objectIDs, objID)
	}

	return uc.TestRepo.GetTestTemplatesByIDs(ctx, userID, objectIDs)
}

func (uc *TestTemplateUseCase) GetTestTemplates(ctx context.Context, userID string) ([]entity.TestTemplete, error) {
	return uc.TestRepo.GetTestTemplates(ctx, userID)
}

func (uc *TestTemplateUseCase) UpdateTestTemplate(ctx context.Context, test entity.TestTemplete) (entity.TestTemplete, error) {
	test.UpdatedAt = time.Now()
	return uc.TestRepo.UpdateTestTemplate(ctx, test)
}

func (uc *TestTemplateUseCase) DeleteTestTemplate(ctx context.Context, userID string, testID primitive.ObjectID) error {
	return uc.TestRepo.DeleteTestTemplate(ctx, userID, testID)
}
