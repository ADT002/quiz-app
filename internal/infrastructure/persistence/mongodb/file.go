package persistence

import (
	"context"
	"errors"

	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const filesDB = "dbapp"
const filesColl = "files"

type FileMongoRepository struct {
	coll *mongo.Collection
}

func NewFileMongoRepository() repository.FileRepository {
	coll := GetMongoClient().Database(filesDB).Collection(filesColl)
	// Indexes: per CLAUDE.md I3.
	// - (owner_id, filename) unique → cấm trùng tên trong cùng owner.
	// - owner_id alone → list-by-owner pagination.
	_, _ = coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "owner_id", Value: 1}, {Key: "filename", Value: 1}},
			Options: options.Index().
				SetName("uniq_owner_filename").
				SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "owner_id", Value: 1}, {Key: "created_at", Value: -1}},
			Options: options.Index().SetName("owner_created"),
		},
	})
	return &FileMongoRepository{coll: coll}
}

func (r *FileMongoRepository) Create(
	ctx context.Context,
	file entity.File,
) (any, error) {
	return r.coll.InsertOne(ctx, file)
}

func (r *FileMongoRepository) FindByID(
	ctx context.Context,
	id primitive.ObjectID,
) (*entity.File, error) {
	var f entity.File
	if err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&f); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &f, nil
}

func (r *FileMongoRepository) FindByOwnerAndFilename(
	ctx context.Context,
	ownerID, filename string,
) (*entity.File, error) {
	var f entity.File
	err := r.coll.
		FindOne(ctx, bson.M{"owner_id": ownerID, "filename": filename}).
		Decode(&f)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &f, nil
}

// ListByOwner returns paginated files belonging to ownerID, optionally filtered
// by content-type prefix (e.g. "image/", "audio/", or "" for all).
func (r *FileMongoRepository) ListByOwner(
	ctx context.Context,
	ownerID, mimePrefix string,
	page, limit int64,
) ([]entity.File, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	filter := bson.M{"owner_id": ownerID}
	if mimePrefix != "" {
		filter["content_type"] = bson.M{"$regex": "^" + mimePrefix}
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip((page - 1) * limit).
		SetLimit(limit)

	cur, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var items []entity.File
	if err := cur.All(ctx, &items); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *FileMongoRepository) DeleteOwned(
	ctx context.Context,
	id primitive.ObjectID,
	ownerID string,
) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": id, "owner_id": ownerID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
