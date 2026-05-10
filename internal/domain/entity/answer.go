package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExportUserSubmission struct {
	Email     string `bson:"email" json:"email"`
	UserID    string `bson:"user_id" json:"user_id"`
	FirstName string `bson:"first_name" json:"first_name"`
	LastName  string `bson:"last_name" json:"last_name"`

	Score     float64 `bson:"score" json:"score"`
	Submitted bool    `bson:"submitted" json:"submitted"`

	StartTime *time.Time `bson:"start_time,omitempty" json:"start_time,omitempty"`
	EndTime   *time.Time `bson:"end_time,omitempty" json:"end_time,omitempty"`
}

type ExportTest struct {
	TestID      primitive.ObjectID `bson:"test_id" json:"test_id"`
	TestName    string             `bson:"test_name" json:"test_name"`
	Description string             `bson:"descript" json:"descript"`

	StartTime       time.Time `bson:"start_time" json:"start_time"`
	EndTime         time.Time `bson:"end_time" json:"end_time"`
	DurationMinutes int       `bson:"duration_minutes" json:"duration_minutes"`

	TestScore  float64 `bson:"test_score" json:"test_score"`
	AuthorMail string  `bson:"author_mail" json:"author_mail"`

	Users []ExportUserSubmission `bson:"users" json:"users"`
}

type ClassExport struct {
	ClassID   primitive.ObjectID `bson:"class_id" json:"class_id"`
	ClassName string             `bson:"class_name" json:"class_name"`
	Test      ExportTest         `bson:"test" json:"test"`
}
