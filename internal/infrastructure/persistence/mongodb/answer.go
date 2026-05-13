package persistence

import (
	"context"
	"strings"
	"time"

	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubmissionMongoRepository struct {
	Collection *mongo.Collection
}

func NewSubmissionMongoRepository() repository.SubmissionRepository {
	client := GetMongoClient()
	return &SubmissionMongoRepository{
		Collection: client.Database("dbapp").Collection("submissions"),
	}
}

func (ar *SubmissionMongoRepository) ListFinishedSubmissionsForClass(
	ctx context.Context,
	classID primitive.ObjectID,
) ([]entity.FinishedSubmissionRow, error) {
	cur, err := ar.Collection.Find(ctx, bson.M{
		"class_id": classID,
		"status": bson.M{
			"$in": bson.A{"submitted", "auto_submitted"},
		},
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []entity.FinishedSubmissionRow
	for cur.Next(ctx) {
		var row entity.FinishedSubmissionRow
		if err := cur.Decode(&row); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, cur.Err()
}

func (ar *SubmissionMongoRepository) ExportData(
	ctx context.Context,
	classID primitive.ObjectID,
	testID primitive.ObjectID,
) (*entity.ClassExport, error) {
	client := GetMongoClient()
	db := client.Database("dbapp")

	var toc entity.TestOfClass
	err := db.Collection("test_of_class").FindOne(ctx, bson.M{
		"_id":      testID,
		"class_id": classID,
	}).Decode(&toc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	var cls entity.Class
	if err := db.Collection("classes").FindOne(ctx, bson.M{"_id": classID}).Decode(&cls); err != nil {
		return nil, err
	}

	cur, err := ar.Collection.Find(ctx, bson.M{
		"test_of_class_id": testID,
		"class_id":         classID,
		"status": bson.M{
			"$in": bson.A{"submitted", "auto_submitted"},
		},
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	submittedStudent := make(map[string]bool)
	var submittedUsers []entity.ExportUserSubmission

	for cur.Next(ctx) {
		var s struct {
			StudentID   string     `bson:"student_id"`
			Email       string     `bson:"email"`
			Score       float64    `bson:"score"`
			StartedAt   *time.Time `bson:"started_at,omitempty"`
			SubmittedAt *time.Time `bson:"submitted_at,omitempty"`
		}
		if err := cur.Decode(&s); err != nil {
			return nil, err
		}
		submittedStudent[s.StudentID] = true

		fn, ln := nameFromClassRoster(cls, s.StudentID)
		em := strings.TrimSpace(s.Email)
		if em == "" {
			em = s.StudentID
		}

		submittedUsers = append(submittedUsers, entity.ExportUserSubmission{
			Email:     em,
			UserID:    s.StudentID,
			FirstName: fn,
			LastName:  ln,
			Score:     s.Score,
			Submitted: true,
			StartTime: s.StartedAt,
			EndTime:   s.SubmittedAt,
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	var unsubmitted []entity.ExportUserSubmission
	for _, st := range cls.StudentsAccept {
		if submittedStudent[st.StudentID] {
			continue
		}
		fn, ln := splitDisplayName(st.StudentName)
		em := strings.TrimSpace(st.StudentEmail)
		if em == "" {
			em = st.StudentID
		}
		unsubmitted = append(unsubmitted, entity.ExportUserSubmission{
			Email:     em,
			UserID:    st.StudentID,
			FirstName: fn,
			LastName:  ln,
			Score:     0,
			Submitted: false,
		})
	}

	all := append(append([]entity.ExportUserSubmission{}, submittedUsers...), unsubmitted...)

	return &entity.ClassExport{
		ClassID:   classID,
		ClassName: cls.ClassName,
		Test: entity.ExportTest{
			TestID:          testID,
			TestName:        toc.TestName,
			Description:     toc.Descript,
			StartTime:       parseExportTime(toc.StartTime),
			EndTime:         parseExportTime(toc.EndTime),
			DurationMinutes: toc.DurationMinutes,
			TestScore:       float64(toc.TestScore),
			AuthorMail:      toc.AuthorMail,
			Users:           all,
		},
	}, nil
}

func nameFromClassRoster(cls entity.Class, userID string) (first, last string) {
	for _, st := range cls.StudentsAccept {
		if st.StudentID == userID {
			return splitDisplayName(st.StudentName)
		}
	}
	return "", ""
}

func splitDisplayName(full string) (first, last string) {
	full = strings.TrimSpace(full)
	if full == "" {
		return "", ""
	}
	parts := strings.Fields(full)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.Join(parts[1:], " ")
}

func parseExportTime(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}
	}
	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
	}
	for _, ly := range layouts {
		if t, err := time.Parse(ly, s); err == nil {
			return t
		}
	}
	return time.Time{}
}
