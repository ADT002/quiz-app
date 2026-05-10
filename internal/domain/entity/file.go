package entity

import (
	"path/filepath"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// File metadata stored in Mongo. The S3 object itself is stored under an
// OPAQUE key (file_id-based), never under a path containing user email or
// other PII — see CLAUDE.md invariant E#15.
type File struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"file_id"`
	OwnerID     string             `bson:"owner_id" json:"-"`           // teacher user_id; never returned to non-owner
	Filename    string             `bson:"filename" json:"filename"`    // display name
	S3Key       string             `bson:"s3_key" json:"-"`             // server-only, opaque
	ContentType string             `bson:"content_type" json:"content_type"`
	Size        int64              `bson:"size" json:"size"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

// BuildS3Key returns an opaque S3 object key derived from the file_id.
// Includes the original extension (lowercased) so S3 + clients pick the right
// content-type when streaming. Does NOT include any user-identifying info.
func BuildS3Key(id primitive.ObjectID, filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return "files/" + id.Hex()
	}
	return "files/" + id.Hex() + ext
}

// NewFile constructs a File ready for insert. Caller fills S3Key after
// computing it from ID via BuildS3Key.
func NewFile(ownerID, filename, contentType string, size int64) File {
	return File{
		OwnerID:     ownerID,
		Filename:    filename,
		ContentType: contentType,
		Size:        size,
		CreatedAt:   time.Now().UTC(),
	}
}
