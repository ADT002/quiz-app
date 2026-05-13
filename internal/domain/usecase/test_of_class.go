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
	Submissions     repository.SubmissionRepository
}

func NewTestOfClassUseCase(
	repo repository.TestOfClassRepository,
	submissions repository.SubmissionRepository,
) *TestOfClassUseCase {
	return &TestOfClassUseCase{
		TestOfClassRepo: repo,
		Submissions:     submissions,
	}
}

func (uc *TestOfClassUseCase) CreateTestOfClass(ctx context.Context, test entity.TestOfClass) (primitive.ObjectID, error) {
	return uc.TestOfClassRepo.CreateTestOfClass(ctx, test)
}

func (uc *TestOfClassUseCase) GetTestsOfClass(ctx context.Context, classID primitive.ObjectID) ([]entity.TestOfClass, error) {
	tests, err := uc.TestOfClassRepo.GetTestsOfClass(ctx, classID)
	if err != nil {
		return nil, err
	}
	rows, err := uc.Submissions.ListFinishedSubmissionsForClass(ctx, classID)
	if err != nil {
		return nil, err
	}
	byTest := make(map[primitive.ObjectID][]entity.TestOfClassUserSubmit)
	for _, r := range rows {
		em := r.Email
		if em == "" {
			em = r.StudentID
		}
		byTest[r.TestOfClassID] = append(byTest[r.TestOfClassID], entity.TestOfClassUserSubmit{
			UserEmail:    em,
			Score:        r.Score,
			EmailID:      r.StudentID,
			SubmissionID: r.SubmissionID.Hex(),
		})
	}
	for i := range tests {
		if list := byTest[tests[i].ID]; len(list) > 0 {
			tests[i].UserSubmit = list
		}
	}
	return tests, nil
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
