package usecase

import (
	"context"

	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileUseCase struct {
	repo repository.FileRepository
}

func NewFileUseCase(repo repository.FileRepository) *FileUseCase {
	return &FileUseCase{repo: repo}
}

func (uc *FileUseCase) Create(ctx context.Context, file entity.File) (any, error) {
	return uc.repo.Create(ctx, file)
}

func (uc *FileUseCase) FindByID(ctx context.Context, id primitive.ObjectID) (*entity.File, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *FileUseCase) FindByOwnerAndFilename(
	ctx context.Context,
	ownerID, filename string,
) (*entity.File, error) {
	return uc.repo.FindByOwnerAndFilename(ctx, ownerID, filename)
}

func (uc *FileUseCase) ListByOwner(
	ctx context.Context,
	ownerID, mimePrefix string,
	page, limit int64,
) ([]entity.File, int64, error) {
	return uc.repo.ListByOwner(ctx, ownerID, mimePrefix, page, limit)
}

func (uc *FileUseCase) DeleteOwned(
	ctx context.Context,
	id primitive.ObjectID,
	ownerID string,
) error {
	return uc.repo.DeleteOwned(ctx, id, ownerID)
}
