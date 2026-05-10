package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type FileAWSRepository struct {
	S3Client *s3.S3
	Bucket   string
}

func NewFileAWSRepository(bucket string) *FileAWSRepository {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-2"),
		Credentials: credentials.NewEnvCredentials(),
	}))

	return &FileAWSRepository{
		S3Client: s3.New(sess),
		Bucket:   bucket,
	}
}

// PresignUploadURL signs a PUT URL for the given opaque S3 key.
// The key MUST NOT contain user-identifying info — see CLAUDE.md E#15.
func (r *FileAWSRepository) PresignUploadURL(
	ctx context.Context,
	s3Key string,
	contentType string,
	exp time.Duration,
) (string, error) {
	req, _ := r.S3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(r.Bucket),
		Key:         aws.String(s3Key),
		ContentType: aws.String(contentType),
	})
	return req.Presign(exp)
}

// PresignDownloadURL signs a GET URL for the given opaque S3 key.
func (r *FileAWSRepository) PresignDownloadURL(
	ctx context.Context,
	s3Key string,
	exp time.Duration,
) (string, error) {
	req, _ := r.S3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(r.Bucket),
		Key:    aws.String(s3Key),
	})
	return req.Presign(exp)
}

func (r *FileAWSRepository) DeleteObject(
	ctx context.Context,
	s3Key string,
) error {
	_, err := r.S3Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.Bucket),
		Key:    aws.String(s3Key),
	})
	return err
}
