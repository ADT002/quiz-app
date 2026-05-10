package usecase

import (
	"context"
	"time"

	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassUseCase struct {
	repoClass       repository.ClassRepository
	repoTestOfClass repository.TestOfClassRepository
}

func NewClassUseCase(
	repoClass repository.ClassRepository,
	repoTestOfClass repository.TestOfClassRepository,
) *ClassUseCase {
	return &ClassUseCase{
		repoClass:       repoClass,
		repoTestOfClass: repoTestOfClass,
	}
}

func (uc *ClassUseCase) CreateClass(ctx context.Context, class entity.Class) (primitive.ObjectID, error) {
	class.CreatedAt = time.Now()
	class.UpdatedAt = time.Now()
	class.ID = primitive.NewObjectID()
	return uc.repoClass.CreateClass(ctx, class)
}

func (uc *ClassUseCase) UpdateClass(ctx context.Context, class entity.Class) (entity.Class, error) {
	return uc.repoClass.UpdateClass(ctx, class)
}

func (uc *ClassUseCase) DeleteClass(ctx context.Context, userID string, id primitive.ObjectID) error {
	return uc.repoClass.DeleteClass(ctx, userID, id)
}

func (uc *ClassUseCase) GetAllClass(ctx context.Context, userID string) ([]entity.Class, error) {
	return uc.repoClass.GetClasses(ctx, userID)
}

func (uc *ClassUseCase) JoinClass(ctx context.Context, classID primitive.ObjectID, emailAuthor, emailIDStudent, emailStudent string) error {
	return uc.repoClass.JoinClass(ctx, classID, emailAuthor, emailIDStudent, emailStudent)
}

func (uc *ClassUseCase) ApproveStudent(ctx context.Context, classID primitive.ObjectID, ownerID, studentID string) error {
	return uc.repoClass.ApproveStudent(ctx, classID, ownerID, studentID)
}

func (uc *ClassUseCase) RejectStudent(ctx context.Context, classID primitive.ObjectID, ownerID, studentID string) error {
	return uc.repoClass.RejectStudent(ctx, classID, ownerID, studentID)
}

func (uc *ClassUseCase) GetAllTestOfClass(ctx context.Context, email string, id primitive.ObjectID) ([]entity.TestOfClass, error) {
	return uc.repoTestOfClass.GetTestsOfClass(ctx, id)
}

func (uc *ClassUseCase) GetQuestionOfTest(ctx context.Context, classID, testID primitive.ObjectID, email string) ([]primitive.ObjectID, entity.TestOfClass, error) {
	test, err := uc.repoTestOfClass.GetTestOfClassByID(ctx, testID)
	if err != nil {
		return nil, entity.TestOfClass{}, err
	}

	var questionIDs []primitive.ObjectID
	for _, qID := range test.QuestionIDs {
		objID, err := primitive.ObjectIDFromHex(qID)
		if err != nil {
			return nil, entity.TestOfClass{}, err
		}
		questionIDs = append(questionIDs, objID)
	}

	test.QuestionIDs = nil
	return questionIDs, test, nil
}
