package repository

import (
	"context"
	"quiz-app/internal/domain/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubmissionRepository interface {
	ExportData(
		ctx context.Context,
		classID primitive.ObjectID,
		testID primitive.ObjectID,
	) (*entity.ClassExport, error)

	// ListFinishedSubmissionsForClass trả mọi bài đã nộp (submitted / auto_submitted) trong lớp.
	ListFinishedSubmissionsForClass(
		ctx context.Context,
		classID primitive.ObjectID,
	) ([]entity.FinishedSubmissionRow, error)
}
