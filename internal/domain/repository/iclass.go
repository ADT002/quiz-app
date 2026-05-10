package repository

import (
	"context"
	entity "quiz-app/internal/domain/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassRepository interface {
	CreateClass(
		ctx context.Context,
		class entity.Class) (primitive.ObjectID, error)

	GetClasses(
		ctx context.Context,
		userID string) ([]entity.Class, error)

	UpdateClass(
		ctx context.Context,
		class entity.Class) (entity.Class, error)

	DeleteClass(
		ctx context.Context,
		userID string,
		id primitive.ObjectID) error

	JoinClass(
		ctx context.Context,
		classID primitive.ObjectID,
		emailAuthor, emailIDStudent,
		emailID string) error

	// ApproveStudent moves student_id từ students_wait → students_accept (atomic).
	// Owner check: ràng buộc qua filter user_id = ownerID.
	ApproveStudent(
		ctx context.Context,
		classID primitive.ObjectID,
		ownerID, studentID string,
	) error

	// RejectStudent xoá student_id khỏi students_wait. Không thêm vào accept.
	RejectStudent(
		ctx context.Context,
		classID primitive.ObjectID,
		ownerID, studentID string,
	) error

	GetAllTestOfClass(
		ctx context.Context,
		email string,
		id primitive.ObjectID,
	) ([]entity.TestOfClass, error)

	GetQuestionOfTest(
		ctx context.Context,
		classID, testID primitive.ObjectID,
		email string,
	) ([]primitive.ObjectID, entity.TestOfClass, error)
}
