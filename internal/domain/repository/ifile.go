package repository

import (
	"context"

	entity "quiz-app/internal/domain/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FileRepository — owner-scoped CRUD over `files` collection.
// All read paths take owner_id; cross-owner lookup goes through FindByID
// (auth check lives in handler, not repo).
type FileRepository interface {
	Create(ctx context.Context, file entity.File) (any, error)

	FindByID(ctx context.Context, id primitive.ObjectID) (*entity.File, error)

	FindByOwnerAndFilename(
		ctx context.Context,
		ownerID, filename string,
	) (*entity.File, error)

	ListByOwner(
		ctx context.Context,
		ownerID, mimePrefix string,
		page, limit int64,
	) ([]entity.File, int64, error)

	DeleteOwned(
		ctx context.Context,
		id primitive.ObjectID,
		ownerID string,
	) error
}
